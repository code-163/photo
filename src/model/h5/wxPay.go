package h5

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"log"
	"photo/src/function"
	"photo/src/model/common"
)

// GetWxPayData 获取微信支付商户数据
func GetWxPayData() (int, string, map[string]string) {
	var code int
	var msg string
	var data map[string]string
	rows, _ := common.Db.Query("select * from wxPayList where status=1")
	defer rows.Close()
	result := common.AssemblyData(rows)
	if len(result) > 0 {
		n := function.RandNum(0, len(result)-1)
		code = 200
		msg = "获取商户成功"
		data = result[n]
	} else {
		code = 400
		msg = "支付维护中"
	}
	return code, msg, data
}

// GetWxPayInfo 获取微信商户数据
func GetWxPayInfo(c *gin.Context) (int, string, map[string]string) {
	var code int
	var msg string
	var data map[string]string
	mchid := c.Query("mchid")
	query, _ := common.Db.Query("select * from wxPayList where mchid=?", mchid)
	defer query.Close()
	result := common.AssemblyData(query)
	fmt.Println(result)
	if len(result) > 0 {
		data = result[0]
		code = 200
		msg = "获取商户成功"
	} else {
		code = 400
		msg = "获取商户失败"
	}
	return code, msg, data
}

// 获取微信商户数据
func getWxPayInfo(mchid string) (int, string, map[string]string) {
	var code int
	var msg string
	var data map[string]string
	query, _ := common.Db.Query("select * from wxPayList where mchid=?", mchid)
	defer query.Close()
	result := common.AssemblyData(query)
	if len(result) > 0 {
		data = result[0]
		code = 200
		msg = "获取商户成功"
	} else {
		code = 400
		msg = "获取商户失败"
	}
	return code, msg, data
}

// ExampleJsapiApiServicePrepay 预下单
func ExampleJsapiApiServicePrepay(OutTradeNo, openid, mchid string, price int64) (int, map[string]string, string) {
	code, msg, result := getWxPayInfo(mchid)
	if code == 200 {
		var (
			mchID                      string = result["mchid"]       // 商户号
			mchCertificateSerialNumber string = result["mchnumber"]   // 商户证书序列号
			mchAPIv3Key                string = result["mchapiv3key"] // 商户APIv3密钥
		)
		// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
		mchPrivateKey, err := utils.LoadPrivateKeyWithPath(result["mchprivatekey"])
		if err != nil {
			log.Fatal("load merchant private key error")
		}
		ctx := context.Background()
		// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
		opts := []core.ClientOption{
			option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
		}
		client, err := core.NewClient(ctx, opts...)
		if err != nil {
			log.Fatalf("new wechat pay client err:%s", err)
		}
		svc := jsapi.JsapiApiService{Client: client}
		// 得到prepay_id，以及调起支付所需的参数和签名
		resp, result, err := svc.PrepayWithRequestPayment(ctx,
			jsapi.PrepayRequest{
				Appid:       core.String(result["appid"]),
				Mchid:       core.String(result["mchid"]),
				Description: core.String("VIP充值"),
				OutTradeNo:  core.String(OutTradeNo),
				Attach:      core.String("VIP充值"),
				NotifyUrl:   core.String("http://193.134.209.55:9100/notify/wxNotify/" + result["mchid"]),
				Amount: &jsapi.Amount{
					Total: core.Int64(price),
				},
				Payer: &jsapi.Payer{
					Openid: core.String(openid),
				},
			},
		)
		if err == nil {
			marshal, _ := json.Marshal(resp)
			data := make(map[string]string)
			_ = json.Unmarshal(marshal, &data)
			log.Println(result.Response.StatusCode)
			log.Println(data["appId"])
			log.Println(data["nonceStr"])
			log.Println(data["package"])
			log.Println(data["paySign"])
			log.Println(data["timeStamp"])
			log.Println(data["prepay_id"])
			return result.Response.StatusCode, data, "微信下单成功"
		} else {
			log.Println(err)
			return result.Response.StatusCode, nil, "微信下单失败"
		}
	} else {
		return 400, nil, msg
	}
}

// ExampleH5ApiServicePrepay 微信支付H5下单
func ExampleH5ApiServicePrepay(OutTradeNo string, price int64) (int, string, error) {
	var (
		mchID                      string = "1626088118"                               // 商户号
		mchCertificateSerialNumber string = "5095E93564A26B15F0F30FAC38A61DC34A6C20D2" // 商户证书序列号
		mchAPIv3Key                string = "2mJTOxiMwOVgPlJjK7ssgDv1QWLGFxb8"         // 商户APIv3密钥
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("./wxpay_cert/apiclient_key.pem")
	if err != nil {
		log.Print("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}
	svc := h5.H5ApiService{Client: client}
	resp, result, err := svc.Prepay(ctx,
		h5.PrepayRequest{
			Appid:       core.String("wx6348655ec475f575"),
			Mchid:       core.String("1626088118"),
			Description: core.String("VIP充值"),
			OutTradeNo:  core.String(OutTradeNo),
			NotifyUrl:   core.String("http://193.134.209.55:9100/notify/wxNotifyH5"),
			Amount: &h5.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(price),
			},
			SceneInfo: &h5.SceneInfo{
				H5Info: &h5.H5Info{
					Type: core.String("Wap"),
				},
				PayerClientIp: core.String("193.134.209.55"),
			},
		},
	)

	if err == nil {
		marshal, _ := json.Marshal(resp)
		data := make(map[string]string)
		_ = json.Unmarshal(marshal, &data)
		fmt.Println(result.Response.StatusCode)
		fmt.Println(data["h5_url"])
		return result.Response.StatusCode, data["h5_url"], nil
	} else {
		fmt.Println(err)
		return 400, "", err
	}
}
