package goo_wxpay

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/liqiongtao/goo"
	"github.com/liqiongtao/goo/utils"
	"strings"
	"time"
)

type UnifiedOrderRequest struct {
	Appid          string    `xml:"appid"`
	MchId          string    `xml:"mch_id"`
	NonceStr       string    `xml:"nonce_str"`
	Sign           string    `xml:"sign"`
	SignType       SignType  `xml:"sign_type"`
	Body           string    `xml:"body"`
	Detail         string    `xml:"detail"`
	Attach         string    `xml:"attach"`
	OutTradeNo     string    `xml:"out_trade_no"`
	TotalFee       int64     `xml:"total_fee"`
	SpbillCreateIp string    `xml:"spbill_create_ip"`
	GoodsTag       string    `xml:"goods_tag"`
	NotifyUrl      string    `xml:"notify_url"`
	TradeType      TradeType `xml:"trade_type"`
	ProductId      string    `xml:"product_id"`
	LimitPay       string    `xml:"limit_pay"`
	Openid         string    `xml:"openid"`
	Receipt        string    `xml:"receipt"`
	SceneInfo      string    `xml:"scene_info"`
}

func (uo *UnifiedOrderRequest) toXml(apiKey string) []byte {
	if uo.NonceStr == "" {
		uo.NonceStr = utils.NonceStr()
	}
	if uo.SignType == "" {
		uo.SignType = SIGN_TYPE_HMAC_SHA256
	}

	str := obj2querystring(uo) + fmt.Sprintf("&key=%s", apiKey)
	goo.Log.Debug("[wxpay-unified-order][req-query-string]", str)

	if uo.SignType == SIGN_TYPE_HMAC_SHA256 {
		uo.Sign = strings.ToUpper(utils.HMacSha256([]byte(str), []byte(apiKey)))
	} else if uo.SignType == SIGN_TYPE_MD5 {
		uo.Sign = strings.ToUpper(utils.MD5([]byte(str)))
	}

	return obj2xml(uo)
}

type UnifiedOrderResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`

	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`

	Appid    string `xml:"appid"`
	MchId    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`

	TradeType TradeType `xml:"trade_type"`
	PrepayId  string    `xml:"prepay_id"`
	CodeUrl   string    `xml:"code_url"`
}

func (uo *UnifiedOrderResponse) toJsApi(apiKey string, signType SignType) map[string]interface{} {
	data := map[string]interface{}{
		"appId":     uo.Appid,
		"timeStamp": fmt.Sprintf("%d", time.Now().Unix()),
		"nonceStr":  uo.NonceStr,
		"package":   fmt.Sprintf("prepay_id=%s", uo.PrepayId),
		"signType":  signType,
		"paySign":   "",
	}

	str := map2querystring(data) + fmt.Sprintf("&key=%s", apiKey)
	goo.Log.Debug("[wxpay-unified-order][jsapi-query-string]", str)

	if signType == SIGN_TYPE_HMAC_SHA256 {
		data["paySign"] = strings.ToUpper(utils.HMacSha256([]byte(str), []byte(apiKey)))
	} else if signType == SIGN_TYPE_MD5 {
		data["paySign"] = strings.ToUpper(utils.MD5([]byte(str)))
	}

	data["timestamp"] = data["timeStamp"]

	return data
}

func (uo *UnifiedOrderResponse) toApp(apiKey string, signType SignType) map[string]interface{} {
	data := map[string]interface{}{
		"appId":        uo.Appid,
		"partnerId":    uo.MchId,
		"prepayId":     uo.PrepayId,
		"packageValue": "Sign=WXPay",
		"nonceStr":     uo.NonceStr,
		"timeStamp":    fmt.Sprintf("%d", time.Now().Unix()),
		"sign":         "",
	}

	str := map2querystring(data) + fmt.Sprintf("&key=%s", apiKey)
	goo.Log.Debug("[wxpay-unified-order][app-query-string]", str)

	if signType == SIGN_TYPE_HMAC_SHA256 {
		data["sign"] = strings.ToUpper(utils.HMacSha256([]byte(str), []byte(apiKey)))
	} else if signType == SIGN_TYPE_MD5 {
		data["sign"] = strings.ToUpper(utils.MD5([]byte(str)))
	}

	return data
}

func UnifiedOrder(req *UnifiedOrderRequest, apiKey string) (map[string]interface{}, error) {
	buf := req.toXml(apiKey)
	goo.Log.Debug("[wxpay-unified-order][req-xml]", string(buf))

	rstBuf, err := goo.NewRequest().Post(URL_UNIFIED_ORDER, buf)
	if err != nil {
		goo.Log.Error("[wxpay-unified-order]", err.Error())
		return nil, err
	}

	goo.Log.Debug("[wxpay-unified-order][rsp-xml]", string(rstBuf))

	rsp := UnifiedOrderResponse{}
	if err := xml.Unmarshal(rstBuf, &rsp); err != nil {
		goo.Log.Error("[wxpay-unified-order]", err.Error())
		return nil, err
	}
	if rsp.ReturnCode == FAIL {
		goo.Log.Error("[wxpay-unified-order]", rsp.ReturnMsg)
		return nil, errors.New(rsp.ReturnMsg)
	}
	if rsp.ResultCode == FAIL {
		goo.Log.Error("[wxpay-unified-order]", rsp.ErrCodeDes)
		return nil, errors.New(rsp.ErrCodeDes)
	}

	if rsp.TradeType == TRADE_TYPE_JSAPI {
		return rsp.toJsApi(apiKey, req.SignType), nil
	}

	if rsp.TradeType == TRADE_TYPE_APP {
		return rsp.toApp(apiKey, req.SignType), nil
	}

	return nil, nil
}
