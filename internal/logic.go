package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/util"
	sdk "github.com/go-pay/gopay/wechat"
)

var Logic = new(logic)

type logic struct{}

func (w *logic) handle(bm gopay.BodyMap) (*sdk.UnifiedOrderResponse, error) {
	rsp, err := Wechat.Client.UnifiedOrder(context.Background(), bm)
	if err != nil {
		return nil, err
	}

	ok, err := sdk.VerifySign(Wechat.Conf.Key, sdk.SignType_MD5, rsp)
	if err != nil {
		return nil, errors.New("验签失败")
	}
	if ok && rsp.ReturnCode == "SUCCESS" && rsp.ResultCode == "SUCCESS" {
		return rsp, nil
	}
	return nil, errors.New(rsp.ReturnMsg)
}

// 构建微信统一支付所需的数据
func (w *logic) body(name, no string, amount int, tt, ip string) gopay.BodyMap {
	bm := make(gopay.BodyMap)
	bm.Set("appid", Wechat.Client.AppId).
		Set("mchid", Wechat.Client.MchId).
		Set("nonce_str", util.RandomString(32)).
		Set("body", name).
		Set("out_trade_no", no).
		Set("total_fee", amount).
		Set("notify_url", Wechat.Conf.Uri("/notice")).
		Set("trade_type", tt).
		Set("spbill_create_ip", ip)
	return bm
}

func (w *logic) Callback(bm gopay.BodyMap) {
	ok, err := sdk.VerifySign(Wechat.Conf.Key, sdk.SignType_MD5, bm)
	if err != nil {
		fmt.Println("验签出现异常")
	}
	if !ok {
		fmt.Println("验签失败")
	}

	if bm.Get("return_code") == "SUCCESS" && bm.Get("result_code") == "SUCCESS" {
		fmt.Println("支付成功")
	} else {
		fmt.Println("微信支付结果通知(支付失败)")
	}
}

// H5 微信H5支付
func (w *logic) H5(name, no string, amount int, ip string) (string, error) {
	bm := w.body(name, no, amount, "MWEB", ip).
		SetBodyMap("scene_info", func(b gopay.BodyMap) {
			b.
				Set("type", "wap").
				Set("wap_url", Wechat.Conf.Uri("/result.html")).
				Set("wap_name", "HiSeer天赋报告")
		})

	rsp, err := w.handle(bm)
	if err != nil {
		return "", err
	}
	return rsp.MwebUrl, nil
}
