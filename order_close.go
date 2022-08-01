package goo_wxpay

import (
	"encoding/xml"
	"errors"
	"fmt"
	goo_http_request "github.com/liqiongtao/googo.io/goo-http-request"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"strings"
)

type CloseOrderRequest struct {
	Appid      string   `xml:"appid"`
	MchId      string   `xml:"mch_id"`
	OutTradeNo string   `xml:"out_trade_no"`
	NonceStr   string   `xml:"nonce_str"`
	Sign       string   `xml:"sign"`
	SignType   SignType `xml:"sign_type"`
}

func (co *CloseOrderRequest) toXml(apiKey string) []byte {
	if co.NonceStr == "" {
		co.NonceStr = goo_utils.NonceStr()
	}
	if co.SignType == "" {
		co.SignType = SIGN_TYPE_HMAC_SHA256
	}

	str := obj2querystring(co) + fmt.Sprintf("&key=%s", apiKey)
	goo_log.WithTag("wxpay-order-close").WithField("query-string", str).Debug()

	if co.SignType == SIGN_TYPE_HMAC_SHA256 {
		co.Sign = strings.ToUpper(goo_utils.HMacSha256([]byte(str), []byte(apiKey)))
	} else if co.SignType == SIGN_TYPE_MD5 {
		co.Sign = strings.ToUpper(goo_utils.MD5([]byte(str)))
	}

	return obj2xml(co)
}

type CloseOrderResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`

	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`

	Appid    string `xml:"appid"`
	MchId    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`
}

func CloseOrder(req *CloseOrderRequest, apiKey string) error {
	buf := req.toXml(apiKey)
	goo_log.WithTag("wxpay-close-order").WithField("req-xml", string(buf)).Debug()

	rstBuf, err := goo_http_request.Post(URL_ORDER_QUERY, buf)
	if err != nil {
		goo_log.WithTag("wxpay-close-order").Error(err.Error())
		return err
	}

	goo_log.WithTag("wxpay-close-order").WithField("res-xml", string(rstBuf)).Debug()

	rsp := &CloseOrderResponse{}
	if err := xml.Unmarshal(rstBuf, rsp); err != nil {
		goo_log.WithTag("wxpay-close-order").Error(err.Error())
		return err
	}
	if rsp.ReturnCode == FAIL {
		goo_log.WithTag("wxpay-close-order").Error(rsp.ReturnMsg)
		return errors.New(rsp.ReturnMsg)
	}
	if rsp.ResultCode == FAIL {
		goo_log.WithTag("wxpay-close-order").Error(rsp.ErrCodeDes)
		return errors.New(rsp.ErrCodeDes)
	}

	return nil
}
