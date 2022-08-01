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

type QueryOrderRequest struct {
	Appid         string   `xml:"appid"`
	MchId         string   `xml:"mch_id"`
	TransactionId string   `xml:"transaction_id"`
	OutTradeNo    string   `xml:"out_trade_no"`
	NonceStr      string   `xml:"nonce_str"`
	Sign          string   `xml:"sign"`
	SignType      SignType `xml:"sign_type"`
}

func (qo *QueryOrderRequest) toXml(apiKey string) []byte {
	if qo.NonceStr == "" {
		qo.NonceStr = goo_utils.NonceStr()
	}
	if qo.SignType == "" {
		qo.SignType = SIGN_TYPE_HMAC_SHA256
	}

	str := obj2querystring(qo) + fmt.Sprintf("&key=%s", apiKey)
	goo_log.WithTag("unified-order").WithField("query-string", str).Debug()

	if qo.SignType == SIGN_TYPE_HMAC_SHA256 {
		qo.Sign = strings.ToUpper(goo_utils.HMacSha256([]byte(str), []byte(apiKey)))
	} else if qo.SignType == SIGN_TYPE_MD5 {
		qo.Sign = strings.ToUpper(goo_utils.MD5([]byte(str)))
	}

	return obj2xml(qo)
}

type QueryOrderResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`

	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`

	Appid    string `xml:"appid"`
	MchId    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`

	DeviceInfo         string    `xml:"device_info"`
	Openid             string    `xml:"openid"`
	IsSubscribe        string    `xml:"is_subscribe"`
	TradeType          TradeType `xml:"trade_type"`
	TradeState         string    `xml:"trade_state"`
	BankType           string    `xml:"bank_type"`
	TotalFee           int64     `xml:"total_fee"`
	SettlementTotalFee int64     `xml:"settlement_total_fee"`
	CashFee            int64     `xml:"cash_fee"`
	CouponFee          int64     `xml:"coupon_fee"`
	CouponCount        int64     `xml:"coupon_count"`
	OutTradeNo         string    `xml:"out_trade_no"`
	Attach             string    `xml:"attach"`
	TimeEnd            string    `xml:"time_end"`
	TradeStateDesc     string    `xml:"trade_state_desc"`
}

func QueryOrder(req *QueryOrderRequest, apiKey string) (*QueryOrderResponse, error) {
	buf := req.toXml(apiKey)
	goo_log.WithTag("wxpay-order-query").WithField("req-xml", string(buf)).Debug()

	rstBuf, err := goo_http_request.Post(URL_ORDER_QUERY, buf)
	if err != nil {
		goo_log.WithTag("wxpay-order-query").Error(err.Error())
		return nil, err
	}

	goo_log.WithTag("wxpay-order-query").WithField("res-xml", string(rstBuf)).Debug()

	rsp := &QueryOrderResponse{}
	if err := xml.Unmarshal(rstBuf, rsp); err != nil {
		goo_log.WithTag("wxpay-order-query").Error(err.Error())
		return nil, err
	}
	if rsp.ReturnCode == FAIL {
		goo_log.WithTag("wxpay-order-query").Error(rsp.ReturnMsg)
		return nil, errors.New(rsp.ReturnMsg)
	}
	if rsp.ResultCode == FAIL {
		goo_log.WithTag("wxpay-order-query").Error(rsp.ErrCodeDes)
		return nil, errors.New(rsp.ErrCodeDes)
	}
	if rsp.TradeState != SUCCESS {
		goo_log.WithTag("wxpay-order-query").Error(tradeStateMsg[rsp.TradeState])
		return nil, errors.New(tradeStateMsg[rsp.TradeState])
	}

	return rsp, nil
}
