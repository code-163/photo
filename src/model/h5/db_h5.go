package h5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"photo/src/function"
	"photo/src/model/common"
	"photo/src/wlogs"
	"strconv"
	"strings"
	"time"
)

const sKey = "1850202134616888"

// GetWorksList 获取作品列表
func GetWorksList(c *gin.Context) []map[string]string {
	var where, wheremap, mapres string
	categoryName := c.Query("category_name")
	keywork := c.Query("keywork")
	if keywork != "" {
		where += "POSITION('" + keywork + "' IN title) && "
	} else {
		if categoryName == "free" {
			where += "types=2 && "
		} else if categoryName == "vip" {
			where += "types=0 && "
		} else if categoryName == "svip" {
			where += "types=1 && "
		} else if categoryName == "" {
			where += "types=2 && "
		}
	}
	where += "status=0 && "
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	limit := common.Pagecut(c)
	var rows *sql.Rows
	var result []map[string]string
	if categoryName == "video" {
		rows, _ = common.Db.Query("SELECT * FROM video " + wheremap + " order by id desc " + limit)
		result = common.AssemblyData(rows)
		data := make([]map[string]string, 0)
		for _, val := range result {
			val["num"] = "1"
			data = append(data, val)
		}
		result = data
	} else {
		rows, _ = common.Db.Query("SELECT * FROM works " + wheremap + " order by id desc " + limit)
		result = common.AssemblyData(rows)
		data := make([]map[string]string, 0)
		for _, val := range result {
			arr := strings.Split(val["atlas"], ",")
			val["num"] = strconv.Itoa(len(arr))
			data = append(data, val)
		}
		result = data
	}
	defer rows.Close()
	return result
}

// GetYigouVideoList 获取已购视频
func GetYigouVideoList(c *gin.Context) string {
	var where, wheremap, mapres string
	uid := c.Query("uid")
	ip := c.ClientIP()
	if uid != "" {
		where += "uid='" + uid + "' && "
	} else {
		where += "ip='" + ip + "' && "
	}
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	queryStr := "SELECT m.* FROM `shopping_record` as sr LEFT JOIN movie as m on sr.mid=m.id " + wheremap + " order by sr.id desc " + limit
	rows, _ := common.Db.Query(queryStr)
	defer rows.Close()
	result := common.AssemblyData(rows)
	res, _ := json.Marshal(result)
	resAesEncrypt, _ := function.AesEncrypt(string(res), sKey)
	return resAesEncrypt
}

// RequestRegister 请求注册
func RequestRegister(c *gin.Context) (int, string, int64, string, string) {
	var code int
	var msg string
	var uid int64
	var userName string
	var loginTime string
	username := c.PostForm("username")
	password := c.PostForm("password")
	//ip := c.ClientIP()
	var id int64
	_ = common.Db.QueryRow("select id from user where username=?", username).Scan(&id)
	if id > 0 {
		code = 400
		msg = "账号已被注册"
	} else {
		ctime := time.Now().Format("2006-01-02 15:04:05")
		result, _ := common.Db.Exec("insert into user (`username`,`password`,`ctime`)values(?,?,?)", username, password, ctime)
		resId, _ := result.LastInsertId()
		rowCount, _ := result.RowsAffected()
		if rowCount > 0 {
			loginTime = time.Now().Format("2006-01-02 15:04:05")
			_, _ = common.Db.Exec("update user set loginTime=? where id=?", loginTime, resId)
			//_, _ = common.Db.Exec("update shopping_record set uid=? where ip=? && uid=0", resId, ip)
			code = 200
			msg = "注册成功"
			uid = resId
			userName = username
		} else {
			code = 400
			msg = "注册失败，请联系客服"
		}
	}
	return code, msg, uid, userName, loginTime
}

// RequestLogin 请求登录
func RequestLogin(c *gin.Context) (int, string, int64, string, string) {
	var code int
	var msg string
	var uid int64
	var userName string
	var loginTime string
	username := c.PostForm("username")
	password := c.PostForm("password")
	//ip := c.ClientIP()
	var id int64
	_ = common.Db.QueryRow("select id from user where username=? && password=?", username, password).Scan(&id)
	if id > 0 {
		loginTime = time.Now().Format("2006-01-02 15:04:05")
		_, _ = common.Db.Exec("update user set loginTime=? where id=?", loginTime, id)
		//_, _ = common.Db.Exec("update shopping_record set uid=? where ip=? && uid=0", id, ip)
		code = 200
		msg = "登录成功"
		uid = id
		userName = username
	} else {
		code = 400
		msg = "账号或密码错误"
	}
	return code, msg, uid, userName, loginTime
}

// GetWorksInfo 获取作品信息
func GetWorksInfo(c *gin.Context) (int, string, string, []string) {
	var code int
	var msg string
	var cover string
	var atlas string
	var types int
	var where, wheremap, mapres string
	var level int
	var member_end_time string
	wid := c.PostForm("wid")
	uid := c.PostForm("uid")
	//ip := c.ClientIP()
	where += "id = " + wid + " && "
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	_ = common.Db.QueryRow("SELECT cover,atlas,types FROM `works` "+wheremap).Scan(&cover, &atlas, &types)
	atlasArr := strings.Split(atlas, ",")
	num := strconv.Itoa(len(atlasArr))
	if types != 2 {
		if uid == "" {
			atlasArr = nil //atlasArr[3:7]
			code = 300
			msg = "会员未登录，全部共" + num + "张"
		} else {
			timer := time.Now().Format("2006-01-02 15:04:05")
			_ = common.Db.QueryRow("SELECT level,member_end_time FROM `user` where id=?", uid).Scan(&level, &member_end_time)
			//VIP
			if types == 0 {
				if level == 0 {
					atlasArr = nil //atlasArr[3:7]
					code = 301
					msg = "未购买会员，全部共" + num + "张，需VIP观看"
				} else {
					if timer > member_end_time {
						atlasArr = nil //atlasArr[3:7]
						code = 302
						msg = "会员已过期，全部共" + num + "张"
					} else {
						code = 200
						msg = "获取数据成功"
					}
				}
			}
			//SVIP
			if types == 1 {
				if level < 2 {
					atlasArr = nil //atlasArr[3:7]
					code = 303
					msg = "会员等级不够，全部共" + num + "张，需SVIP观看"
				} else {
					if timer > member_end_time {
						atlasArr = nil //atlasArr[3:7]
						code = 302
						msg = "会员已过期，全部共" + num + "张"
					} else {
						code = 200
						msg = "获取数据成功"
					}
				}
			}
		}
	} else {
		code = 200
		msg = "获取数据成功"
	}
	if cover == "" && len(atlasArr) == 0 {
		code = 400
		msg = "找不到相关内容"
	}
	return code, msg, cover, atlasArr
}

// GetVideoInfo 获取视频信息
func GetVideoInfo(c *gin.Context) (int, string, string, string) {
	var code int
	var msg string
	var playerUrl string
	var playerUri string
	var title string
	var level int
	var member_end_time string
	uid := c.PostForm("uid")
	id := c.PostForm("id")
	if id != "" {
		_ = common.Db.QueryRow("SELECT player_url,title FROM `video` where id=?", id).Scan(&playerUrl, &title)
	}
	if uid == "" {
		code = 300
		msg = "会员未登录，请先登录"
	} else {
		timer := time.Now().Format("2006-01-02 15:04:05")
		_ = common.Db.QueryRow("SELECT level,member_end_time FROM `user` where id=?", uid).Scan(&level, &member_end_time)
		//SVIP
		if level < 2 {
			code = 301
			msg = "尊敬的客户，视频区为SVIP专享，升级后观看"
		} else {
			if timer > member_end_time {
				code = 302
				msg = "会员已过期，请购买SVIP"
			} else {
				code = 200
				msg = "获取数据成功"
				playerUri = playerUrl
			}
		}
	}
	return code, msg, playerUri, title
}

// GetAgentPrice 获取代理销售金额
func GetAgentPrice(c *gin.Context) (int, string, map[string]string) {
	var code int
	var msg string
	var single string
	var day string
	var week string
	var month string

	dataConfig := make(map[string]string)
	_ = common.Db.QueryRow("select single_amount,day_amount,week_amount,month_amount from `config`").Scan(&single, &day, &week, &month)
	dataConfig["single_amount"] = single
	dataConfig["day_amount"] = day
	dataConfig["week_amount"] = week
	dataConfig["month_amount"] = month

	acode := c.PostForm("code")
	data := make(map[string]string)
	if acode != "" {
		var single_amount string
		var day_amount string
		var week_amount string
		var month_amount string
		var id int64
		_ = common.Db.QueryRow("select id,single_amount,day_amount,week_amount,month_amount from agent where code=?", acode).Scan(&id, &single_amount, &day_amount, &week_amount, &month_amount)
		if id > 0 {
			data["single_amount"] = single_amount
			data["day_amount"] = day_amount
			data["week_amount"] = week_amount
			data["month_amount"] = month_amount
		} else {
			data = dataConfig
		}
	} else {
		data = dataConfig
	}
	code = 200
	msg = "获取金额成功"
	return code, msg, data
}

// GoOrder 购买会员创建订单
func GoOrder(c *gin.Context) (int, string, string) {
	var code int
	var msg string
	var orderNo string
	var acode string
	var recode string
	coder, message, resData := GetWxPayData()
	if coder == 200 {
		vip := c.PostForm("vip")
		uid := c.PostForm("uid")
		price := c.PostForm("price")
		ip := c.ClientIP()
		if c.PostForm("code") != "" {
			acode = c.PostForm("code")
		} else {
			if uid != "" {
				_ = common.Db.QueryRow("select code from user where id=?", uid).Scan(&recode)
				acode = recode
			}
		}
		ctime := time.Now().Format("2006-01-02 15:04:05")
		endTime := time.Now().Add(time.Second * 90).Format("2006-01-02 15:04:05")
		orderNo = "A" + time.Now().Format("20060102150405") + function.CreateRandomNumber(5)
		result, _ := common.Db.Exec("insert into orderbuy (`code`,`order_no`,`price`,`ctime`,`end_time`,`vip`,`uid`,`ip`,`mchid`)values(?,?,?,?,?,?,?,?,?)", acode, orderNo, price, ctime, endTime, vip, uid, ip, resData["mchid"])
		rowCount, _ := result.RowsAffected()
		if rowCount > 0 {
			code = 200
			msg = "订单创建成功"
		} else {
			code = 400
			msg = "订单创建失败，请联系客服"
		}
	} else {
		code = 400
		msg = message
	}
	return code, msg, orderNo
}

// GoExOrder 推广页创建订单
func GoExOrder(c *gin.Context) (int, string, string, string) {
	var code int
	var msg string
	var orderNo string
	var mchid string
	coder, message, resData := GetWxPayData()
	if coder == 200 {
		vip := c.PostForm("vip")
		acode := c.PostForm("code")
		uuid := c.PostForm("uuid")
		ip := c.ClientIP()
		if vip == "1" {
			price := 38
			ctime := time.Now().Format("2006-01-02 15:04:05")
			endTime := time.Now().Add(time.Second * 90).Format("2006-01-02 15:04:05")
			orderNo = "E" + time.Now().Format("20060102150405") + function.CreateRandomNumber(5)
			result, _ := common.Db.Exec("insert into orderbuy (`code`,`order_no`,`price`,`ctime`,`end_time`,`vip`,`ip`,`uuid`,`mchid`)values(?,?,?,?,?,?,?,?,?)", acode, orderNo, price, ctime, endTime, vip, ip, uuid, resData["mchid"])
			rowCount, _ := result.RowsAffected()
			if rowCount > 0 {
				code = 200
				msg = "订单创建成功"
				mchid = resData["mchid"]
			} else {
				code = 400
				msg = "订单创建失败，请联系客服"
			}
		} else {
			code = 200
			msg = "订单创建失败，参数异常"
		}
	} else {
		code = 400
		msg = message
	}
	return code, msg, orderNo, mchid
}

// GetOrderInfo 获取订单数据
func GetOrderInfo(c *gin.Context) (int, string, map[string]interface{}) {
	var code int
	var msg string
	var data = make(map[string]interface{})
	var id int64
	var price float64
	var order_no string
	var username string
	var vip int
	var mchid string
	orderNo := c.PostForm("orderNo")
	if orderNo != "" {
		_ = common.Db.QueryRow("SELECT o.id,o.price,o.order_no,o.vip,o.mchid,u.username FROM `orderbuy` as o left join user as u on o.uid=u.id where o.order_no=?", orderNo).Scan(&id, &price, &order_no, &vip, &mchid, &username)
		if id > 0 {
			data["price"] = price
			data["order_no"] = order_no
			data["vip"] = vip
			data["username"] = username
			data["mchid"] = mchid
			code = 200
			msg = "获取订单数据成功"
		} else {
			code = 400
			msg = "订单数据不存在"
		}
	} else {
		code = 400
		msg = "参数缺失，请求失败"
	}

	return code, msg, data
}

// CheckOrder 查询订单状态
func CheckOrder(c *gin.Context) (int, string) {
	var code int
	var msg string
	var status int
	orderNo := c.Query("orderNo")
	if orderNo != "" {
		_ = common.Db.QueryRow("select status from orderbuy where order_no=?", orderNo).Scan(&status)
		if status == 1 {
			code = 200
			msg = "订单支付成功"
		} else if status == 2 {
			code = 500
			msg = "订单已超时"
		} else {
			code = 400
			msg = "订单未支付"
		}
	} else {
		code = 400
		msg = "缺少订单号"
	}
	return code, msg
}

// CheckOrderData 通过ip查询订单状态
func CheckOrderData(c *gin.Context) (int, string, string) {
	var code int
	var msg string
	var uid string
	var orderNo string
	var resOrderNo string
	uuid := c.Query("uuid")
	if uuid != "" {
		_ = common.Db.QueryRow("select uid,order_no from orderbuy where status=1 && uuid=?", uuid).Scan(&uid, &orderNo)
		if orderNo != "" {
			if uid != "" {
				code = 300
				msg = "订单支付成功"
				resOrderNo = orderNo
			} else {
				code = 200
				msg = "订单支付成功"
				resOrderNo = orderNo
			}
		} else {
			code = 400
			msg = "订单不存在"
		}
	} else {
		code = 400
		msg = "缺少参数"
	}
	return code, msg, resOrderNo
}

// RequestWxJsapiPay 请求微信JSAPI支付
func RequestWxJsapiPay(c *gin.Context) (int, string, map[string]string) {
	var code int
	var msg string
	var price float64
	openid := c.PostForm("openid")
	orderNo := c.PostForm("orderNo")
	mchid := c.PostForm("mchid")
	_ = common.Db.QueryRow("select price from orderbuy where order_no=?", orderNo).Scan(&price)
	coder, data, message := ExampleJsapiApiServicePrepay(orderNo, openid, mchid, int64(price)*100)
	if coder == 200 && data != nil {
		code = coder
		msg = message
	}
	return code, msg, data
}

// RequestWxH5Pay 请求微信H5支付
func RequestWxH5Pay(c *gin.Context) (int, string, string) {
	var code int
	var msg string
	var h5Url string
	var price float64
	orderNo := c.PostForm("orderNo")
	_ = common.Db.QueryRow("select price from orderbuy where order_no=?", orderNo).Scan(&price)
	code, h5Url, err := ExampleH5ApiServicePrepay(orderNo, int64(price)*100)
	if code == 200 {
		msg = "获取支付地址成功"
	} else {
		msg = err.Error()
	}
	return code, msg, h5Url
}

// RequestAliPay 请求支付宝支付
func RequestAliPay(c *gin.Context) (int, string, string) {
	var code int
	var msg string
	var payurl string
	price := c.PostForm("price")
	order_no := c.PostForm("orderNo")
	return_url := c.PostForm("return_url")
	usercode := "dpJlV4Rk"
	notify_url := "http://45.144.138.61:9100/notify/order"
	payUri := "http://124.248.67.93:7200/order"
	sign := function.Md5(usercode + "|" + price + "|" + order_no + "|" + notify_url + "|" + return_url)
	var uri url.URL
	data := uri.Query()
	data.Add("usercode", usercode)
	data.Add("price", price)
	data.Add("order_no", order_no)
	data.Add("return_url", return_url)
	data.Add("notify_url", notify_url)
	data.Add("types", "1")
	data.Add("sign", sign)
	queryStr := data.Encode()
	resp, err := http.Post(payUri, "application/x-www-form-urlencoded", strings.NewReader(queryStr))
	resData := make(map[string]interface{})
	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			_ = json.Unmarshal(body, &resData)
			if resData["code"].(float64) == 200 {
				code = 200
				msg = "请求支付成功"
				payurl = resData["payurl"].(string)
			} else {
				code = 400
				msg = resData["msg"].(string)
			}
		} else {
			code = 400
			msg = "支付网关请求失败"
		}
	}
	return code, msg, payurl
}

// GetUserPackage 获取会员套餐状态
func GetUserPackage(c *gin.Context) (int, string) {
	var code int
	var msg string
	uid := c.PostForm("uid")
	if uid != "" {
		var member_end_time string
		_ = common.Db.QueryRow("select member_end_time from user where id=?", uid).Scan(&member_end_time)
		timeNow := time.Now().Format("2006-01-02 15:04:05")
		if member_end_time > timeNow {
			code = 200
			msg = "会员套餐正常"
		} else {
			code = 400
			msg = "会员套餐到期"
		}
	} else {
		code = 400
		msg = "会员未登录"
	}
	return code, msg
}

// CheckLoginStatus 检查会员登录状态
func CheckLoginStatus(c *gin.Context) (int, string) {
	var code int
	var msg string
	uid := c.PostForm("uid")
	loginTime := c.PostForm("loginTime")
	var LoginTime string
	_ = common.Db.QueryRow("select loginTime from user where id=?", uid).Scan(&LoginTime)
	if LoginTime == loginTime {
		code = 200
		msg = "会员已登录"
	} else {
		code = 400
		msg = "会员已退出，被其他设备登上"
	}
	return code, msg
}

// InsertLog 记录访问日志
func InsertLog(c *gin.Context) (int, string) {
	var code int
	var msg string
	types := c.PostForm("types")
	loginTime := time.Now().Format("2006-01-02 15:04:05")
	ip := c.ClientIP()
	_ = wlogs.WriteLog("COS访问日志", "设备类型："+types+" => 访问时间："+loginTime+" => ip："+ip)
	code = 200
	msg = "testLog"
	return code, msg
}

// GetDoorDomain 获取入口域名
func GetDoorDomain(c *gin.Context) (int, string, string) {
	var code int
	var msg string
	var domain string
	rows, _ := common.Db.Query("select domain from door_domain where status=0")
	defer rows.Close()
	result := common.AssemblyData(rows)
	num := len(result)
	if num > 0 {
		n := function.RandNum(0, num-1)
		domain = result[n]["domain"]
		code = 200
		msg = "获取成功"
	} else {
		code = 400
		msg = "无可用域名"
	}
	return code, msg, domain
}

// GetJumpDomain 获取中转域名
func GetJumpDomain(c *gin.Context) (int, string, string) {
	var code int
	var msg string
	var domain string
	rows, _ := common.Db.Query("select domain from jump_domain where status=0")
	defer rows.Close()
	result := common.AssemblyData(rows)
	num := len(result)
	if num > 0 {
		n := function.RandNum(0, num-1)
		domain = result[n]["domain"]
		code = 200
		msg = "获取成功"
	} else {
		code = 400
		msg = "无可用域名"
	}
	return code, msg, domain
}

// GetDomain 获取落地域名
func GetDomain(c *gin.Context) (int, string, string) {
	var code int
	var msg string
	var domain string
	rows, _ := common.Db.Query("select domain from luodi_domain where status=0")
	defer rows.Close()
	result := common.AssemblyData(rows)
	num := len(result)
	if num > 0 {
		n := function.RandNum(0, num-1)
		domain = result[n]["domain"]
		code = 200
		msg = "获取成功"
	} else {
		code = 400
		msg = "无可用域名"
	}
	return code, msg, domain
}

// GetOpenid 获取微信用户openid
func GetOpenid(c *gin.Context) (int, string, string) {
	var code int
	var msg string
	var openid string
	coder := c.Query("code")
	mchid := c.Query("mchid")
	codes, msg, result := getWxPayInfo(mchid)
	if codes == 200 {
		if coder != "" {
			uri := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + result["appid"] + "&secret=" + result["secret"] + "&code=" + coder + "&grant_type=authorization_code"
			resp, err := http.Get(uri)
			if err != nil {
				fmt.Println(err)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode == 200 {
					body, _ := io.ReadAll(resp.Body)
					resData := make(map[string]interface{})
					_ = json.Unmarshal(body, &resData)
					fmt.Println(string(body))
					fmt.Println("openid：", resData["openid"])
					if resData["openid"] != nil {
						code = 200
						msg = "获取OPENID成功"
						openid = resData["openid"].(string)
					} else {
						code = 400
						msg = "获取OPENID失败"
					}
				} else {
					code = 400
					msg = "获取OPENID请求失败"
				}
			}
		} else {
			code = 400
			msg = "缺少参数，code"
		}
	} else {
		code = 400
		msg = "获取商户数据失败"
	}
	return code, msg, openid
}

// UpdateUserPasswd 会员修改密码
func UpdateUserPasswd(c *gin.Context) (int, string) {
	var code int
	var msg string
	uid := c.PostForm("uid")
	passwd := c.PostForm("passwd")
	if uid == "" {
		code = 400
		msg = "会员未登录，请先登录"
	} else {
		result, _ := common.Db.Exec("update user set password=? where id=?", passwd, uid)
		affected, _ := result.RowsAffected()
		if affected > 0 {
			code = 200
			msg = "修改密码成功"
		} else {
			code = 400
			msg = "修改密码失败，请联系管理员"
		}
	}
	return code, msg
}

// GetUserInfo 获取会员数据
func GetUserInfo(c *gin.Context) (int, string, map[string]string) {
	var code int
	var msg string
	var data map[string]string
	uid := c.Query("uid")
	if uid != "" {
		rows, _ := common.Db.Query("select * from user where id=?", uid)
		defer rows.Close()
		result := common.AssemblyData(rows)
		data = result[0]
		code = 200
		msg = "数据获取成功"
	} else {
		code = 400
		msg = "会员未登录，请先登录"
	}
	return code, msg, data
}

// UpdateUserRegister 会员注册更新订单套餐
func UpdateUserRegister(c *gin.Context) (int, string, string, int64, string) {
	var code int
	var msg string
	var newMemberEndTime string
	var resId int64
	var userId int64
	var resUsername string
	username := c.PostForm("username")
	passwd := c.PostForm("passwd")
	orderNo := c.PostForm("orderNo")
	coder := c.PostForm("code")
	ctime := time.Now().Format("2006-01-02 15:04:05")
	if username == "" || passwd == "" {
		code = 400
		msg = "账号和密码不能为空"
	} else {
		_ = common.Db.QueryRow("select id from user where username=?", username).Scan(&resId)
		if orderNo != "" {
			var vip int
			var uid string
			var orderId int64
			_ = common.Db.QueryRow("select id,vip,uid from orderbuy where order_no=?", orderNo).Scan(&orderId, &vip, &uid)
			if orderId > 0 {
				if uid != "" {
					code = 400
					msg = "注册失败，此订单号已被使用"
				} else {
					if resId > 0 {
						code = 400
						msg = "此账号已被注册"
					} else {
						result, _ := common.Db.Exec("insert into user (`username`,`password`,`ctime`,`code`)values(?,?,?,?)", username, passwd, ctime, coder)
						id, _ := result.LastInsertId()
						if id > 0 {
							if vip == 1 {
								//VIP包月
								newMemberEndTime = time.Now().AddDate(0, 1, 0).Format("2006-01-02 15:04:05")
							} else if vip == 2 {
								//SVIP1包一年
								newMemberEndTime = time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05")
							} else if vip == 3 {
								//SVIP2包两年
								newMemberEndTime = time.Now().AddDate(2, 0, 0).Format("2006-01-02 15:04:05")
							}
							exec, _ := common.Db.Exec("update user set level=?,member_end_time=?,loginTime=? where id=?", vip, newMemberEndTime, ctime, id)
							affected, _ := exec.RowsAffected()
							if affected > 0 {
								_, _ = common.Db.Exec("update orderbuy set uid=? where order_no=?", id, orderNo)
								code = 200
								msg = "注册会员成功"
								resUsername = username
								userId = id
							} else {
								code = 400
								msg = "注册成功，但更新会员套餐失败"
							}
						} else {
							code = 400
							msg = "注册会员失败"
						}
					}
				}
			} else {
				code = 400
				msg = "订单不存在，无法注册"
			}
		} else {
			if resId > 0 {
				code = 400
				msg = "此账号已被注册"
			} else {
				result, _ := common.Db.Exec("insert into user (`username`,`password`,`ctime`,`loginTime`,`code`)values(?,?,?,?,?)", username, passwd, ctime, ctime, coder)
				id, _ := result.LastInsertId()
				if id > 0 {
					code = 200
					msg = "注册会员成功"
					resUsername = username
					userId = id
				}
			}
		}
	}
	return code, msg, resUsername, userId, ctime
}

// AddVisit 添加访问统计
func AddVisit(c *gin.Context) (int, string) {
	var code int
	var msg string
	types := c.Query("types")
	acode := c.Query("code")
	ip := c.ClientIP()
	ctime := time.Now().Format("2006-01-02")
	exec, _ := common.Db.Exec("insert into visit (`code`,`ctime`,`ip`,`types`)values (?,?,?,?)", acode, ctime, ip, types)
	affected, _ := exec.RowsAffected()
	if affected > 0 {
		code = 200
		msg = "添加统计成功"
	} else {
		code = 400
		msg = "添加统计失败"
	}
	return code, msg
}
