package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"photo/src/function"
	"photo/src/model/common"
	"photo/src/wlogs"
	"sync"
	"time"
)

var locks sync.Mutex

// WxNotifyH5 微信支付H5回调通知
func WxNotifyH5(c *gin.Context) string {
	var code string
	var (
		mchID                      string = "1626088118"                               // 商户号
		mchCertificateSerialNumber string = "5095E93564A26B15F0F30FAC38A61DC34A6C20D2" // 商户证书序列号
		mchAPIv3Key                string = "2mJTOxiMwOVgPlJjK7ssgDv1QWLGFxb8"         // 商户APIv3密钥
	)
	ctx := context.Background()
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	mchPrivateKey, _ := utils.LoadPrivateKeyWithPath("./wxpay_cert/apiclient_key.pem")
	_ = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, mchPrivateKey, mchCertificateSerialNumber, mchID, mchAPIv3Key)
	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mchID)
	// 3. 使用证书访问器初始化 `notify.Handler`
	handler, _ := notify.NewRSANotifyHandler(mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	content := make(map[string]interface{})
	notifyReq, err := handler.ParseNotifyRequest(context.Background(), c.Request, content)
	// 如果验签未通过，或者解密失败
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// 处理通知内容
	resData := make(map[string]interface{})
	_ = json.Unmarshal([]byte(notifyReq.Resource.Plaintext), &resData)
	_ = wlogs.WriteLog("微信H5回调", notifyReq.Resource.Plaintext)
	if resData["out_trade_no"].(string) != "" && resData["trade_state"].(string) == "SUCCESS" {
		code = "SUCCESS"
		go updateOrder(resData["out_trade_no"].(string), "")
	} else {
		code = "FAIL"
	}
	return code
}

// WxNotify 微信支付Jsapi回调通知
func WxNotify(c *gin.Context) string {
	var code string
	mchid := c.Param("mchid")
	if mchid != "" {
		rows, _ := common.Db.Query("select * from wxPayList where mchid=?", mchid)
		defer rows.Close()
		result := common.AssemblyData(rows)
		data := result[0]
		var (
			mchID                      string = data["mchid"]       // 商户号
			mchCertificateSerialNumber string = data["mchnumber"]   // 商户证书序列号
			mchAPIv3Key                string = data["mchapiv3key"] // 商户APIv3密钥
		)
		ctx := context.Background()
		// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
		mchPrivateKey, _ := utils.LoadPrivateKeyWithPath(data["mchprivatekey"])
		_ = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, mchPrivateKey, mchCertificateSerialNumber, mchID, mchAPIv3Key)
		// 2. 获取商户号对应的微信支付平台证书访问器
		certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mchID)
		// 3. 使用证书访问器初始化 `notify.Handler`
		handler, _ := notify.NewRSANotifyHandler(mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
		var content interface{}
		notifyReq, err := handler.ParseNotifyRequest(context.Background(), c.Request, content)
		// 如果验签未通过，或者解密失败
		if err != nil {
			fmt.Println(err)
			return ""
		}
		// 处理通知内容
		resData := make(map[string]interface{})
		_ = json.Unmarshal([]byte(notifyReq.Resource.Plaintext), &resData)
		_ = wlogs.WriteLog("微信Jsapi回调", mchid+" => "+notifyReq.Resource.Plaintext)
		if resData["out_trade_no"].(string) != "" && resData["trade_state"].(string) == "SUCCESS" {
			code = "SUCCESS"
			go updateOrder(resData["out_trade_no"].(string), data["mchid"])
		} else {
			code = "FAIL"
		}
	} else {
		_ = wlogs.WriteLog("Jsapi回调错误", "缺少mchid")
		code = "FAIL"
	}
	return code
}

//更新订单状态
func updateOrder(orderNo, mchid string) {
	pay_time := time.Now().Format("2006-01-02 15:04:05")
	result, _ := common.Db.Exec("update orderbuy set status=?,pay_time=?,mchid=? where order_no=? && status != 1", 1, pay_time, mchid, orderNo)
	rowCount, err := result.RowsAffected()
	function.CheckErr("回调更新订单失败", err)
	if rowCount > 0 {
		go cutUserOrder(orderNo)
		go updateAgentPrice(orderNo)
	}
}

//订单用户套餐处理
func cutUserOrder(orderNo string) {
	var ip string
	var vip int
	var uid string
	var memberEndTime string
	var newMemberEndTime string
	common.Db.QueryRow("select ip,vip,uid from orderbuy where status=1 && order_no=?", orderNo).Scan(&ip, &vip, &uid)
	common.Db.QueryRow("select member_end_time from user where id=?", uid).Scan(&memberEndTime)
	timer := time.Now().Format("2006-01-02 15:04:05")
	if uid != "" {
		if vip == 1 {
			//VIP包月
			if memberEndTime > timer {
				loc, _ := time.LoadLocation("Local")
				theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", memberEndTime, loc)
				newMemberEndTime = theTime.AddDate(0, 1, 0).Format("2006-01-02 15:04:05")
			} else {
				newMemberEndTime = time.Now().AddDate(0, 1, 0).Format("2006-01-02 15:04:05")
			}
			_, _ = common.Db.Exec("update user set member_end_time=?,level=? where id=?", newMemberEndTime, vip, uid)
		} else if vip == 2 {
			//SVIP1包一年
			if memberEndTime > timer {
				loc, _ := time.LoadLocation("Local")
				theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", memberEndTime, loc)
				newMemberEndTime = theTime.AddDate(1, 0, 0).Format("2006-01-02 15:04:05")
			} else {
				newMemberEndTime = time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05")
			}
			_, _ = common.Db.Exec("update user set member_end_time=?,level=? where id=?", newMemberEndTime, vip, uid)
		} else if vip == 3 {
			//SVIP2包两年
			if memberEndTime > timer {
				loc, _ := time.LoadLocation("Local")
				theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", memberEndTime, loc)
				newMemberEndTime = theTime.AddDate(2, 0, 0).Format("2006-01-02 15:04:05")
			} else {
				newMemberEndTime = time.Now().AddDate(2, 0, 0).Format("2006-01-02 15:04:05")
			}
			_, _ = common.Db.Exec("update user set member_end_time=?,level=? where id=?", newMemberEndTime, vip, uid)
		}
	}
}

//更新代理钱包余额
func updateAgentPrice(orderNo string) {
	locks.Lock()
	var price float64
	var code string
	var point float64
	var pid int64
	var money float64
	common.Db.QueryRow("select price,code from orderbuy where status=1 && order_no=?", orderNo).Scan(&price, &code)
	if code != "" {
		common.Db.QueryRow("select point,pid from agent where code=?", code).Scan(&point, &pid)
		if pid == 0 {
			//一级代理结算
			money = price * (point * 0.01)
			_, _ = common.Db.Exec("update agent set wallet=(wallet+?) where code=?", money, code)
		} else {
			//一级代理结算
			var points float64
			common.Db.QueryRow("select point from agent where id=?", pid).Scan(&points)
			money = price * ((points - point) * 0.01)
			_, _ = common.Db.Exec("update agent set wallet=(wallet+?) where id=?", money, pid)

			//二级代理结算
			amount := price * (point * 0.01)
			_, _ = common.Db.Exec("update agent set wallet=(wallet+?) where code=?", amount, code)
		}
	}
	locks.Unlock()
}
