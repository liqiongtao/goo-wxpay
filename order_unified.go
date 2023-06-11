package goo_wxpay

import (
	"encoding/xml"
	"errors"
	"fmt"
	goo_http_request "github.com/liqiongtao/googo.io/goo-http-request"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
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
	Attach         string    `xml:"attach"`           // 附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	OutTradeNo     string    `xml:"out_trade_no"`     // 商户系统内部订单号，要求32个字符内（最少6个字符），只能是数字、大小写字母_-|*且在同一个商户号下唯一
	TotalFee       int64     `xml:"total_fee"`        // 订单总金额，单位为分
	SpbillCreateIp string    `xml:"spbill_create_ip"` // 必须传正确的用户端IP,支持ipv4、ipv6格式
	GoodsTag       string    `xml:"goods_tag"`        // 商品标记，代金券或立减优惠功能的参数
	NotifyUrl      string    `xml:"notify_url"`       // 接收微信支付异步通知回调地址，通知url必须为直接可访问的url，不能携带参数。
	TradeType      TradeType `xml:"trade_type"`       // 调用接口提交的交易类型，取值如下：JSAPI，NATIVE，APP，,H5支付固定传MWEB
	ProductId      string    `xml:"product_id"`       // trade_type=NATIVE，此参数必传。此id为二维码中包含的商品ID，商户自行定义。
	LimitPay       string    `xml:"limit_pay"`        // no_credit--指定不能使用信用卡支付
	Openid         string    `xml:"openid"`           // trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识
	Receipt        string    `xml:"receipt"`          // Y，传入Y时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效
	SceneInfo      string    `xml:"scene_info"`       // 该字段用于上报支付的场景信息,针对H5支付有以下三种场景,请根据对应场景上报,H5支付不建议在APP端使用，针对场景1，2请接入APP支付，不然可能会出现兼容性问题
}

func (uo *UnifiedOrderRequest) toXml(apiKey string) []byte {
	if uo.NonceStr == "" {
		uo.NonceStr = goo_utils.NonceStr()
	}
	if uo.SignType == "" {
		uo.SignType = SIGN_TYPE_HMAC_SHA256
	}

	str := obj2querystring(uo) + fmt.Sprintf("&key=%s", apiKey)
	goo_log.WithTag("wxpay-unified-order").WithField("query-string", str).Debug()

	if uo.SignType == SIGN_TYPE_HMAC_SHA256 {
		uo.Sign = strings.ToUpper(goo_utils.HMacSha256([]byte(str), []byte(apiKey)))
	} else if uo.SignType == SIGN_TYPE_MD5 {
		uo.Sign = strings.ToUpper(goo_utils.MD5([]byte(str)))
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

	TradeType TradeType `xml:"trade_type"` // 调用接口提交的交易类型，取值如下：JSAPI，NATIVE，APP，,H5支付固定传MWEB
	PrepayId  string    `xml:"prepay_id"`  // 微信生成的预支付会话标识，用于后续接口调用中使用，该值有效期为2小时,针对H5支付此参数无特殊用途
	CodeUrl   string    `xml:"code_url"`
	MWebUrl   string    `xml:"mweb_url"` // mweb_url为拉起微信支付收银台的中间页面，可通过访问该url来拉起微信客户端，完成支付,mweb_url的有效期为5分钟。
}

func (uo *UnifiedOrderResponse) JsApi(apiKey string, signType SignType) map[string]interface{} {
	data := map[string]interface{}{
		"appId":     uo.Appid,
		"timeStamp": fmt.Sprintf("%d", time.Now().Unix()),
		"nonceStr":  uo.NonceStr,
		"package":   fmt.Sprintf("prepay_id=%s", uo.PrepayId),
		"signType":  signType,
		"paySign":   "",
	}

	str := map2querystring(data) + fmt.Sprintf("&key=%s", apiKey)
	goo_log.WithTag("wxpay-unified-order").WithField("query-string", str).Debug()

	if signType == SIGN_TYPE_HMAC_SHA256 {
		data["paySign"] = strings.ToUpper(goo_utils.HMacSha256([]byte(str), []byte(apiKey)))
	} else if signType == SIGN_TYPE_MD5 {
		data["paySign"] = strings.ToUpper(goo_utils.MD5([]byte(str)))
	}

	data["timestamp"] = data["timeStamp"]

	return data
}

func (uo *UnifiedOrderResponse) App(apiKey string, signType SignType) map[string]interface{} {
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
	goo_log.WithTag("wxpay-unified-order").WithField("query-string", str).Debug()

	if signType == SIGN_TYPE_HMAC_SHA256 {
		data["sign"] = strings.ToUpper(goo_utils.HMacSha256([]byte(str), []byte(apiKey)))
	} else if signType == SIGN_TYPE_MD5 {
		data["sign"] = strings.ToUpper(goo_utils.MD5([]byte(str)))
	}

	return data
}

func UnifiedOrder(req *UnifiedOrderRequest, apiKey string) (resp UnifiedOrderResponse, err error) {
	resp = UnifiedOrderResponse{}

	buf := req.toXml(apiKey)
	goo_log.WithTag("wxpay-unified-order").WithField("req-xml", string(buf)).Debug()

	var rstBuf []byte
	rstBuf, err = goo_http_request.New().Debug().Post(URL_UNIFIED_ORDER, buf)
	if err != nil {
		goo_log.WithTag("wxpay-unified-order").Error(err.Error())
		return
	}

	goo_log.WithTag("wxpay-unified-order").WithField("res-xml", string(rstBuf)).Debug()

	if err = xml.Unmarshal(rstBuf, &resp); err != nil {
		goo_log.WithTag("wxpay-unified-order").Error(err.Error())
		return
	}
	if resp.ReturnCode == FAIL {
		goo_log.WithTag("wxpay-unified-order").Error(resp.ReturnMsg)
		err = errors.New(resp.ReturnMsg)
		return
	}
	if resp.ResultCode == FAIL {
		goo_log.WithTag("wxpay-unified-order").Error(resp.ErrCodeDes)
		err = errors.New(resp.ErrCodeDes)
		return
	}

	return
}
