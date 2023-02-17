package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"photo/src/function"
	"photo/src/middleware"
	"photo/src/model/admin"
	"photo/src/model/agent"
	"photo/src/model/common"
	"photo/src/model/h5"
	"photo/src/model/notify"
	"time"
)

const qiniuUrl = "https://7niu.trumall.cn/" //"http://qiniu.xf6699.cn/"

func main() {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.NoRoute(NotPage)
	manager := r.Group("/", HandleAdminToken)
	{
		manager.POST("/checkadminlogin", checkadminlogin)
		manager.POST("/adminlogin", adminlogin)
		manager.POST("/adminlogout", adminlogout)
		manager.GET("/getAgentList", getAgentList)
		manager.GET("/getLevelAgent", getLevelAgent)
		manager.POST("/addAgent", addAgent)
		manager.POST("/editAgent", editAgent)
		manager.GET("/delAgent", delAgent)
		manager.GET("/delLevelAgent", delLevelAgent)
		manager.GET("/getModelList", getModelList)
		manager.GET("/getModel", getModel)
		manager.GET("/getWorks", getWorks)
		manager.POST("/addWorks", addWorks)
		manager.POST("/editWorks", editWorks)
		manager.GET("/delWorks", delWorks)
		manager.GET("/delVideo", delVideo)
		manager.GET("/stateWorks", stateWorks)
		manager.GET("/stateVideo", stateVideo)
		manager.GET("/getWorksData", getWorksData)
		manager.POST("/updateWorksImage", updateWorksImage)
		manager.POST("/editModel", editModel)
		manager.POST("/addModel", addModel)
		manager.GET("/delModel", delModel)
		manager.GET("/getUptoken", getUptoken)
		manager.POST("/uploadPhoto", uploadPhoto)
		manager.POST("/uploadFile", uploadFile)
		manager.GET("/delMovie", delMovie)
		manager.GET("/getVideoList", getVideoList)
		manager.POST("/addVideo", addVideo)
		manager.POST("/editVideo", editVideo)
		manager.GET("/getOrderList", getOrderList)
		manager.GET("/getDrawingLisst", getDrawingLisst)
		manager.GET("/getLuodiDomainList", getLuodiDomainList)
		manager.GET("/getDoorDomainList", getDoorDomainList)
		manager.GET("/getJumpDomainList", getJumpDomainList)
		manager.POST("/addLuodiDomain", addLuodiDomain)
		manager.POST("/addDoorDomain", addDoorDomain)
		manager.POST("/addJumpDomain", addJumpDomain)
		manager.GET("/delLuodiDomain", delLuodiDomain)
		manager.GET("/delDoorDomain", delDoorDomain)
		manager.GET("/delJumpDomain", delJumpDomain)
		manager.GET("/stateDrawing", stateDrawing)
		manager.POST("/aliSettlement", aliSettlement)
		manager.GET("/get_order_sales", get_order_sales)
		manager.GET("/getVisit", getVisit)
		manager.GET("/get_adminloginlog", get_adminloginlog)
		manager.GET("/getWxPayList", getWxPayList)
		manager.POST("/addWxPay", addWxPay)
		manager.POST("/editWxPay", editWxPay)
		manager.GET("/delWxPay", delWxPay)
		manager.GET("/stateWxPay", stateWxPay)
		manager.POST("/editWorkTypes", editWorkTypes)
	}
	agent := r.Group("agent", HandleUserOrigin)
	{
		agent.POST("/agentlogin", agentlogin)
		agent.POST("/addAgentLevel", addAgentLevel)
		agent.POST("/editAgentLevel", editAgentLevel)
		agent.GET("/getAgentLevel", getAgentLevel)
		agent.GET("/delAgentLevel", delAgentLevel)
		agent.POST("/agentlogout", agentlogout)
		agent.POST("/checkAgentlogin", checkAgentlogin)
		agent.GET("/getAgentOrder", getAgentOrder)
		agent.GET("/getAgentOrderLevel", getAgentOrderLevel)
		agent.GET("/getAgentDrawing", getAgentDrawing)
		agent.GET("/getAgentInfo", getAgentInfo)
		agent.POST("/addAgentDrawing", addAgentDrawing)
		agent.GET("/queryAgentSales", queryAgentSales)
		agent.GET("/getAgentLoginlog", getAgentLoginlog)
		agent.GET("/getVisitAgent", getVisitAgent)
	}
	notifys := r.Group("notify")
	{
		notifys.POST("/wxNotify/:mchid", wxNotify)
		notifys.POST("/wxNotifyH5", wxNotifyH5)
	}
	h5 := r.Group("h5")
	{
		h5.GET("/getWorksList", getWorksList)
		h5.POST("/getVideoInfo", getVideoInfo)
		h5.POST("/requestRegister", requestRegister)
		h5.POST("/requestLogin", requestLogin)
		h5.GET("/getYigouVideoList", getYigouVideoList)
		h5.POST("/getWorksInfo", getWorksInfo)
		h5.POST("/getAgentPrice", getAgentPrice)
		h5.POST("/goOrder", goOrder)
		h5.POST("/goExOrder", goExOrder)
		h5.POST("/getOrderInfo", getOrderInfo)
		h5.POST("/requestAliPay", requestAliPay)
		h5.POST("/requestWxJsapiPay", requestWxJsapiPay)
		h5.POST("/requestWxH5Pay", requestWxH5Pay)
		h5.GET("/checkOrder", checkOrder)
		h5.GET("/checkOrderData", checkOrderData)
		h5.POST("/getUserPackage", getUserPackage)
		h5.POST("/checkLoginStatus", checkLoginStatus)
		h5.POST("/insertLog", insertLog)
		h5.POST("/updateUserPasswd", updateUserPasswd)
		h5.POST("/updateUserRegister", updateUserRegister)
		h5.GET("/getDomain", getDomain)
		h5.GET("/getDoorDomain", getDoorDomain)
		h5.GET("/getJumpDomain", getJumpDomain)
		h5.GET("/getOpenid", getOpenid)
		h5.GET("/getUserInfo", getUserInfo)
		h5.GET("/addVisit", addVisit)
		h5.GET("/getIp", getIp)
		h5.GET("/getWxPayInfo", getWxPayInfo)
	}
	_ = r.Run(":9100")
	//r.RunTLS(":9100", "./zhengshu/server.pem", "./zhengshu/server.key")
}

// NotPage 当访问一个错误网站时返回
func NotPage(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status": 404,
		"error":  "找不到对应路径或没有此文件以及页面",
	})
}

// HandleAdminToken 设置中间件，判断管理员的token是否存在，如果不存在则终止请求
func HandleAdminToken(c *gin.Context) {
	origin := c.Request.Header.Get("Origin") //请求头部
	if origin != "http://boss.felitsa.cn" {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "对不起，您没有该接口访问权限，请联系管理员",
		})
		c.Abort() //终止下面的函数执行
	} else {
		path := c.Request.URL.Path
		admincode := c.Query("admincode")
		admintoken := c.Query("admintoken")
		redisdb, ctx, _ := function.CreateRedisClient()
		if path == "/adminlogin" || path == "/checkadminlogin" || path == "/adminlogout" {
			c.Next() //允许下面的函数执行
		} else {
			if admincode != "" && admintoken != "" {
				val, _ := redisdb.Get(ctx, admincode).Result() //获取admintoken
				if val == "" || val != admintoken {
					c.JSON(http.StatusOK, gin.H{
						"code": 400,
						"msg":  "对不起，您没有该接口访问权限，请联系管理员",
					})
					c.Abort() //终止下面的函数执行
				}
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 403,
					"msg":  "对不起，您没有该接口访问权限，请联系管理员",
				})
				c.Abort() //终止下面的函数执行
			}
		}
	}
}

func HandleUserOrigin(c *gin.Context) {
	origin := c.Request.Header.Get("Origin") //请求头部
	if origin != "http://agent.felitsa.cn" {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "对不起，您没有该接口访问权限，请联系管理员",
		})
		c.Abort() //终止下面的函数执行
	}
}

// HandleVisit 添加访问记录
func HandleVisit(c *gin.Context) {
	var id int64
	ip := c.ClientIP()
	code := c.Query("code")
	fmt.Println("code：", code)
	ctime := time.Now().Format("2006-01-02")
	where := fmt.Sprintf("where code='%v' && ip='%v' && ctime='%v'", code, ip, ctime)
	_ = common.Db.QueryRow("select id from `visit` " + where).Scan(&id)
	if id == 0 {
		_, _ = common.Db.Exec("insert into visit (`code`,`ctime`,`ip`)values(?,?,?)", code, ctime, ip)
	}
}

//微信支付Jsapi回调通知
func wxNotify(c *gin.Context) {
	code := notify.WxNotify(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
	})
}

//微信支付H5回调通知
func wxNotifyH5(c *gin.Context) {
	code := notify.WxNotifyH5(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
	})
}

/**
检测管理员登录状态
*/
func checkadminlogin(c *gin.Context) {
	msg, code := admin.Checkadminlogin(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

/**
管理员登录
*/
func adminlogin(c *gin.Context) {
	msg, code, admincode, admintoken, admin_id, account := admin.Adminlogin(c)
	c.JSON(http.StatusOK, gin.H{
		"code":       code,
		"msg":        msg,
		"admincode":  admincode,
		"admintoken": admintoken,
		"admin_id":   admin_id,
		"account":    account,
	})
}

//管理员退出
func adminlogout(c *gin.Context) {
	msg, code := admin.Adminlogout(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//检测代理登录状态
func checkUserlogin(c *gin.Context) {
	msg, code := agent.CheckAgentlogin(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//获取下级代理
func getAgentLevel(c *gin.Context) {
	data, count := agent.GetAgentLevel(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//代理登录
func agentlogin(c *gin.Context) {
	msg, code, agentcodeRes, agenttokenRes, aid, account, coder, pid := agent.Agentlogin(c)
	c.JSON(http.StatusOK, gin.H{
		"code":       code,
		"msg":        msg,
		"agentcode":  agentcodeRes,
		"agenttoken": agenttokenRes,
		"aid":        aid,
		"account":    account,
		"coder":      coder,
		"pid":        pid,
	})
}

//代理退登录
func agentlogout(c *gin.Context) {
	msg, code := agent.Agentlogout(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//检测用户登录状态
func checkAgentlogin(c *gin.Context) {
	msg, code := agent.CheckAgentlogin(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//读取用户登录日志
func getAgentLoginlog(c *gin.Context) {
	data, count := agent.GetAgentLoginlog(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//当前用户读取订单
func getAgentOrder(c *gin.Context) {
	data, count := agent.GetAgentOrder(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//当前下级代理订单
func getAgentOrderLevel(c *gin.Context) {
	data, count := agent.GetAgentOrderLevel(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//当前用户读取提现订单
func getAgentDrawing(c *gin.Context) {
	data, count := agent.GetAgentDrawing(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//当前用户读取账户
func getAgentInfo(c *gin.Context) {
	data := agent.GetAgentInfo(c)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
		"msg":  "",
	})
}

//当前用户提现操作
func addAgentDrawing(c *gin.Context) {
	msg, code := agent.AddAgentDrawing(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//当前用户读取销售数据
func queryAgentSales(c *gin.Context) {
	msg, code, data := agent.QueryAgentSales(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

//获取代理列表
func getAgentList(c *gin.Context) {
	data, count := admin.GetAgentList(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//获取二级代理列表
func getLevelAgent(c *gin.Context) {
	data, count := admin.GetLevelAgent(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//添加代理
func addAgent(c *gin.Context) {
	msg, code := admin.AddAgent(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//添加代理
func addAgentLevel(c *gin.Context) {
	msg, code := agent.AddAgentLevel(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//修改二级代理
func editAgentLevel(c *gin.Context) {
	msg, code := agent.EditAgentLevel(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//删除二级代理
func delAgentLevel(c *gin.Context) {
	msg, code := agent.DelAgentLevel(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//修改代理
func editAgent(c *gin.Context) {
	msg, code := admin.EditAgent(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//删除代理
func delAgent(c *gin.Context) {
	msg, code := admin.DelAgent(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//删除二级代理
func delLevelAgent(c *gin.Context) {
	msg, code := admin.DelLevelAgent(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//模特列表
func getModelList(c *gin.Context) {
	data, count := admin.GetModelList(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//获取模特数据
func getModel(c *gin.Context) {
	data := admin.GetModel(c)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": data,
	})
}

//获取作品列表
func getWorks(c *gin.Context) {
	data, count := admin.GetWorks(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//添加作品
func addWorks(c *gin.Context) {
	code, msg := admin.AddWorks(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//添加视频
func addVideo(c *gin.Context) {
	code, msg := admin.AddVideo(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//修改作品
func editWorks(c *gin.Context) {
	code, msg := admin.EditWorks(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//批量修改作品类型
func editWorkTypes(c *gin.Context) {
	code, msg := admin.EditWorkTypes(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//修改视频
func editVideo(c *gin.Context) {
	code, msg := admin.EditVideo(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//删除作品
func delWorks(c *gin.Context) {
	msg, code := admin.DelWorks(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//删除视频
func delVideo(c *gin.Context) {
	msg, code := admin.DelVideo(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//作品开关
func stateWorks(c *gin.Context) {
	msg, code := admin.StateWorks(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//视频开关
func stateVideo(c *gin.Context) {
	msg, code := admin.StateVideo(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//获取单个作品数据
func getWorksData(c *gin.Context) {
	code, msg, data := admin.GetWorksData(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

//更新作品图片
func updateWorksImage(c *gin.Context) {
	code, msg := admin.UpdateWorksImage(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//添加模特
func addModel(c *gin.Context) {
	msg, code := admin.AddModel(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//修改模特
func editModel(c *gin.Context) {
	msg, code := admin.EditModel(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//删除模特
func delModel(c *gin.Context) {
	msg, code := admin.DelModel(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//获取七牛云上传凭证
func getUptoken(c *gin.Context) {
	uptoken := function.GetUpToken()
	c.JSON(http.StatusOK, gin.H{
		"uptoken":  uptoken,
		"qiniuUrl": qiniuUrl,
	})
}

//上传图片
func uploadPhoto(c *gin.Context) {
	var code int
	var msg string
	var image string
	file, err := c.FormFile("imageFile")
	if err != nil {
		code = 400
		msg = "上传失败"
	} else {
		ok, imageName, localFile := function.SaveImageFile(c, "./images", file)
		if ok {
			imgName := function.UploadFileToQny(imageName, localFile)
			if imgName != "" {
				code = 200
				msg = "上传成功"
				image = qiniuUrl + imgName
			}
		} else {
			code = 400
			msg = "上传失败"
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   msg,
		"image": image,
	})
}

//上传文件
func uploadFile(c *gin.Context) {
	var code int
	var msg string
	var fileName string
	var filePath string
	file, err := c.FormFile("imageFile")
	if err != nil {
		code = 400
		msg = "上传失败"
	} else {
		ok, imageName, localFile := function.SaveFile(c, "./upload", file)
		if ok {
			code = 200
			msg = "上传成功"
			fileName = imageName
			filePath = localFile
		} else {
			code = 400
			msg = "上传失败"
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     code,
		"msg":      msg,
		"fileName": fileName,
		"filePath": filePath,
	})
}

//读取系统所有销量
func get_order_sales(c *gin.Context) {
	msg, code, data := admin.QueryOrderSales(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

//读取系统访问统计
func getVisit(c *gin.Context) {
	msg, code, data := admin.GetVisit(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

//读取系统访问统计
func getVisitAgent(c *gin.Context) {
	msg, code, data := agent.GetVisitAgent(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

//获取管理员登录日志
func get_adminloginlog(c *gin.Context) {
	data, count := admin.QueryAdminLoginLog(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//结算代理余额
func aliSettlement(c *gin.Context) {
	msg, code := admin.AliSettlement(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//删除视频
func delMovie(c *gin.Context) {
	msg, code := admin.DelMovie(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//获取视频列表
func getVideoList(c *gin.Context) {
	data, count := admin.GetVideoList(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "",
		"count": count,
		"data":  data,
	})
}

//获取订单列表
func getOrderList(c *gin.Context) {
	data, count := admin.GetOrderList(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//获取提现订单
func getDrawingLisst(c *gin.Context) {
	data, count := admin.GetDrawingLisst(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//获取入口域名
func getDoorDomainList(c *gin.Context) {
	data, count := admin.GetDoorDomainList(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//获取微信支付列表
func getWxPayList(c *gin.Context) {
	data, count := admin.GetWxPayList(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//获取中转域名
func getJumpDomainList(c *gin.Context) {
	data, count := admin.GetJumpDomainList(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//获取落地域名
func getLuodiDomainList(c *gin.Context) {
	data, count := admin.GetLuodiDomainList(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  data,
		"msg":   "",
	})
}

//添加落地域名
func addLuodiDomain(c *gin.Context) {
	msg, code := admin.AddLuodiDomain(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//添加入口域名
func addDoorDomain(c *gin.Context) {
	msg, code := admin.AddDoorDomain(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//添加微信支付
func addWxPay(c *gin.Context) {
	msg, code := admin.AddWxPay(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

// 修改微信支付
func editWxPay(c *gin.Context) {
	msg, code := admin.EditWxPay(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//删除微信支付
func delWxPay(c *gin.Context) {
	msg, code := admin.DelWxPay(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//微信支付开关
func stateWxPay(c *gin.Context) {
	msg, code := admin.StateWxPay(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//添加中转域名
func addJumpDomain(c *gin.Context) {
	msg, code := admin.AddJumpDomain(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//删除入口域名
func delDoorDomain(c *gin.Context) {
	msg, code := admin.DelDoorDomain(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//删除中转域名
func delJumpDomain(c *gin.Context) {
	msg, code := admin.DelJumpDomain(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//删除落地域名
func delLuodiDomain(c *gin.Context) {
	msg, code := admin.DelLuodiDomain(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//处理提现状态
func stateDrawing(c *gin.Context) {
	msg, code := admin.StateDrawing(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//获取视频列表
func getWorksList(c *gin.Context) {
	data := h5.GetWorksList(c)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"list": data,
	})
}

//获取视频信息
func getVideoInfo(c *gin.Context) {
	code, msg, playerUri, title := h5.GetVideoInfo(c)
	c.JSON(http.StatusOK, gin.H{
		"code":      code,
		"msg":       msg,
		"playerUri": playerUri,
		"title":     title,
	})
}

//获取已购视频列表
func getYigouVideoList(c *gin.Context) {
	data := h5.GetYigouVideoList(c)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"list": data,
	})
}

//请求注册
func requestRegister(c *gin.Context) {
	code, msg, uid, userName, loginTime := h5.RequestRegister(c)
	c.JSON(http.StatusOK, gin.H{
		"code":      code,
		"msg":       msg,
		"uid":       uid,
		"username":  userName,
		"loginTime": loginTime,
	})
}

//请求登录
func requestLogin(c *gin.Context) {
	code, msg, uid, userName, loginTime := h5.RequestLogin(c)
	c.JSON(http.StatusOK, gin.H{
		"code":      code,
		"msg":       msg,
		"uid":       uid,
		"username":  userName,
		"loginTime": loginTime,
	})
}

//获取作品信息
func getWorksInfo(c *gin.Context) {
	code, msg, cover, atlas := h5.GetWorksInfo(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   msg,
		"cover": cover,
		"list":  atlas,
	})
}

//获取代理销售金额
func getAgentPrice(c *gin.Context) {
	code, msg, data := h5.GetAgentPrice(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   msg,
		"price": data,
	})
}

//购买会员创建订单
func goOrder(c *gin.Context) {
	code, msg, orderNo := h5.GoOrder(c)
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"msg":     msg,
		"orderNo": orderNo,
	})
}

//推广页创建订单
func goExOrder(c *gin.Context) {
	code, msg, orderNo, mchid := h5.GoExOrder(c)
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"msg":     msg,
		"orderNo": orderNo,
		"mchid":   mchid,
	})
}

//获取订单数据
func getOrderInfo(c *gin.Context) {
	code, msg, data := h5.GetOrderInfo(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

//请求支付宝支付
func requestAliPay(c *gin.Context) {
	code, msg, payurl := h5.RequestAliPay(c)
	c.JSON(http.StatusOK, gin.H{
		"code":   code,
		"msg":    msg,
		"payurl": payurl,
	})
}

//请求微信JSAPI支付
func requestWxJsapiPay(c *gin.Context) {
	code, msg, data := h5.RequestWxJsapiPay(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

//请求微信H5支付
func requestWxH5Pay(c *gin.Context) {
	code, msg, h5Url := h5.RequestWxH5Pay(c)
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   msg,
		"h5Url": h5Url,
	})
}

//查询订单状态
func checkOrder(c *gin.Context) {
	code, msg := h5.CheckOrder(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//通过标识查询订单状态
func checkOrderData(c *gin.Context) {
	code, msg, orderNo := h5.CheckOrderData(c)
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"msg":     msg,
		"orderNo": orderNo,
	})
}

//获取会员套餐状态
func getUserPackage(c *gin.Context) {
	code, msg := h5.GetUserPackage(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//检查会员登录状态
func checkLoginStatus(c *gin.Context) {
	code, msg := h5.CheckLoginStatus(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//记录访问日志
func insertLog(c *gin.Context) {
	code, msg := h5.InsertLog(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//会员修改密码
func updateUserPasswd(c *gin.Context) {
	code, msg := h5.UpdateUserPasswd(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

//会员注册更新订单套餐
func updateUserRegister(c *gin.Context) {
	code, msg, username, uid, loginTime := h5.UpdateUserRegister(c)
	c.JSON(http.StatusOK, gin.H{
		"code":      code,
		"msg":       msg,
		"username":  username,
		"uid":       uid,
		"loginTime": loginTime,
	})
}

//获取入口域名
func getDoorDomain(c *gin.Context) {
	code, msg, domain := h5.GetDoorDomain(c)
	c.JSON(http.StatusOK, gin.H{
		"code":   code,
		"msg":    msg,
		"domain": domain,
	})
}

//获取中转域名
func getJumpDomain(c *gin.Context) {
	code, msg, domain := h5.GetJumpDomain(c)
	c.JSON(http.StatusOK, gin.H{
		"code":   code,
		"msg":    msg,
		"domain": domain,
	})
}

//获取落地域名
func getDomain(c *gin.Context) {
	code, msg, domain := h5.GetDomain(c)
	c.JSON(http.StatusOK, gin.H{
		"code":   code,
		"msg":    msg,
		"domain": domain,
	})
}

//获取微信商户数据
func getWxPayInfo(c *gin.Context) {
	code, msg, data := h5.GetWxPayInfo(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

//获取客户IP
func getIp(c *gin.Context) {
	var code int
	var msg string
	ip := c.ClientIP()
	if ip != "" {
		code = 200
		msg = "获取IP成功"
	} else {
		code = 400
		msg = "获取IP失败"
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"ip":   ip,
	})
}

//获取微信用户openid
func getOpenid(c *gin.Context) {
	code, msg, openid := h5.GetOpenid(c)
	c.JSON(http.StatusOK, gin.H{
		"code":   code,
		"msg":    msg,
		"openid": openid,
	})
}

//获取会员数据
func getUserInfo(c *gin.Context) {
	code, msg, data := h5.GetUserInfo(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

//添加访问统计
func addVisit(c *gin.Context) {
	code, msg := h5.AddVisit(c)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}
