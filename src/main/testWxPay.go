package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"log"
	"time"
)

// ExampleJsapiApiServicePrepay 预下单
func ExampleJsapiApiServicePrepay(OutTradeNo, openid string, price int64) (int, map[string]string, string) {
	var (
		mchID                      string = "1626088118"                               // 商户号
		mchCertificateSerialNumber string = "5095E93564A26B15F0F30FAC38A61DC34A6C20D2" // 商户证书序列号
		mchAPIv3Key                string = "2mJTOxiMwOVgPlJjK7ssgDv1QWLGFxb8"         // 商户APIv3密钥
	)
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("./wxpay_cert/apiclient_key.pem")
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
			Appid:       core.String("wx6348655ec475f575"),
			Mchid:       core.String("1626088118"),
			Description: core.String("VIP充值"),
			OutTradeNo:  core.String(OutTradeNo),
			Attach:      core.String("VIP充值"),
			NotifyUrl:   core.String("http://193.134.209.55:9100/notify/wxNotify/1626088118"),
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
}

// Prepay 预下单
func Prepay(OutTradeNo string, price int64) (int, string, error) {
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
			Appid:         core.String("wx6348655ec475f575"),
			Mchid:         core.String("1626088118"),
			Description:   core.String("VIP充值"),
			OutTradeNo:    core.String(OutTradeNo),
			TimeExpire:    core.Time(time.Now()),
			Attach:        core.String("VIP充值"),
			NotifyUrl:     core.String("http://45.136.15.246:9100/notify/wxNotifyH5"),
			GoodsTag:      core.String("WXG"),
			LimitPay:      []string{"LimitPay_example"},
			SupportFapiao: core.Bool(false),
			Amount: &h5.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(price),
			},
			Detail: &h5.Detail{
				CostPrice: core.Int64(608800),
				GoodsDetail: []h5.GoodsDetail{h5.GoodsDetail{
					GoodsName:        core.String("iPhoneX 256G"),
					MerchantGoodsId:  core.String("ABC"),
					Quantity:         core.Int64(1),
					UnitPrice:        core.Int64(828800),
					WechatpayGoodsId: core.String("1001"),
				}},
				InvoiceId: core.String("wx123"),
			},
			SceneInfo: &h5.SceneInfo{
				DeviceId: core.String("013467007045764"),
				H5Info: &h5.H5Info{
					AppName:     core.String("王者荣耀"),
					AppUrl:      core.String("https://pay.qq.com"),
					BundleId:    core.String("com.tencent.wzryiOS"),
					PackageName: core.String("com.tencent.tmgp.sgame"),
					Type:        core.String("iOS"),
				},
				PayerClientIp: core.String("127.0.0.1"),
				StoreInfo: &h5.StoreInfo{
					Address:  core.String("广东省深圳市南山区科技中一道10000号"),
					AreaCode: core.String("440305"),
					Id:       core.String("0001"),
					Name:     core.String("腾讯大厦分店"),
				},
			},
			SettleInfo: &h5.SettleInfo{
				ProfitSharing: core.Bool(false),
			},
		},
	)

	if err == nil {
		marshal, _ := json.Marshal(resp)
		data := make(map[string]string)
		_ = json.Unmarshal(marshal, &data)
		//log.Println(result.Response.StatusCode)
		//log.Println(data["h5_url"])
		return result.Response.StatusCode, data["h5_url"], nil
	} else {
		log.Println(err)
		return 400, "", err
	}
}

func main() {
	const openid = "oz4RR6ab6WnkN6Qf4VQHV_4oJv-4"
	/*const name = "oDzJI631Pe3Sjvms2B_LmIaheRIM"
	code, h5Url, _ := Prepay("20231234354657457435", 100)
	fmt.Println(code)
	fmt.Println(h5Url)*/

	code, data, msg := ExampleJsapiApiServicePrepay("O365675683473747457", openid, 100)
	fmt.Println(code)
	fmt.Println(msg)
	fmt.Println(data)
}
