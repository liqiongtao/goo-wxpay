package goo_wxpay

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/liqiongtao/goo"
	"github.com/liqiongtao/goo/utils"
	"log"
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
		co.NonceStr = utils.NonceStr()
	}
	if co.SignType == "" {
		co.SignType = SIGN_TYPE_HMAC_SHA256
	}

	str := obj2querystring(co) + fmt.Sprintf("&key=%s", apiKey)
	log.Println("[UnifiedOrderRequest.querystring]", str)

	if co.SignType == SIGN_TYPE_HMAC_SHA256 {
		co.Sign = strings.ToUpper(utils.HMacSha256([]byte(str), []byte(apiKey)))
	} else if co.SignType == SIGN_TYPE_MD5 {
		co.Sign = strings.ToUpper(utils.MD5([]byte(str)))
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
	log.Println("[CloseOrderRequest.xml]", string(buf))

	rstBuf, err := goo.NewRequest().Post(URL_ORDER_QUERY, buf)
	if err != nil {
		return err
	}

	log.Println("[CloseOrderResponse.xml]", string(rstBuf))

	rsp := &CloseOrderResponse{}
	if err := xml.Unmarshal(rstBuf, rsp); err != nil {
		return err
	}
	if rsp.ReturnCode == FAIL {
		return errors.New(rsp.ReturnMsg)
	}
	if rsp.ResultCode == FAIL {
		return errors.New(rsp.ErrCodeDes)
	}

	return nil
}
