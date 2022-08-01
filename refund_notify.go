package goo_wxpay

import (
	"encoding/xml"
	"errors"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"strings"
)

type RefundData struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	Appid      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	ReqInfo    string `xml:"req_info"`
}

type RefundReqInfo struct {
	TransactionId       string `xml:"transaction_id"`
	OutTradeNo          string `xml:"out_trade_no"`
	RefundId            string `xml:"refund_id"`
	OutRefundNo         string `xml:"out_refund_no"`
	TotalFee            int64  `xml:"total_fee"`
	SettlementTotalFee  int64  `xml:"settlement_total_fee"`
	RefundFee           int64  `xml:"refund_fee"`
	SettlementRefundFee int64  `xml:"settlement_refund_fee"`
	RefundStatus        string `xml:"refund_status"`
	SuccessTime         string `xml:"success_time"`
	RefundRecvAccout    string `xml:"refund_recv_accout"`
	RefundAccount       string `xml:"refund_account"`
	RefundRequestSource string `xml:"refund_request_source"`
}

func RefundNotify(buf []byte, apiKey string) (*RefundReqInfo, error) {
	data := &RefundData{}

	if err := xml.Unmarshal(buf, data); err != nil {
		goo_log.WithTag("wxpay-refund-notify").Error(err.Error())
		return nil, err
	}
	if data.ReturnCode == FAIL {
		goo_log.WithTag("wxpay-refund-notify").Error(data.ReturnMsg)
		return nil, errors.New(data.ReturnMsg)
	}

	base64buf := goo_utils.Base64Decode(data.ReqInfo)
	key := strings.ToUpper(goo_utils.MD5([]byte(apiKey)))
	goo_log.WithTag("wxpay-refund-notify").WithField("base64buf", string(base64buf)).WithField("key", key).Debug()
	infoBuf, err := goo_utils.AESECBDecrypt(base64buf, []byte(key))
	if err != nil {
		goo_log.WithTag("wxpay-refund-notify").Error(err.Error())
		return nil, err
	}

	info := &RefundReqInfo{}
	if err := xml.Unmarshal(infoBuf, info); err != nil {
		goo_log.WithTag("wxpay-refund-notify").Error(err.Error())
		return nil, err
	}
	return info, nil
}
