package goo_wxpay

import (
	"encoding/xml"
	"errors"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"strings"
)

type OrderPayData struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`

	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`

	Appid      string   `xml:"appid"`
	MchId      string   `xml:"mch_id"`
	DeviceInfo string   `xml:"device_info"`
	NonceStr   string   `xml:"nonce_str"`
	Sign       string   `xml:"sign"`
	SignType   SignType `xml:"sign_type"`

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
	TransactionId      string    `xml:"transaction_id"`
	OutTradeNo         string    `xml:"out_trade_no"`
	Attach             string    `xml:"attach"`
	TimeEnd            string    `xml:"time_end"`
}

func OrderPayNotifyVerify(buf []byte, apiKey string) (data *OrderPayData, err error) {
	data = &OrderPayData{}

	if err = xml.Unmarshal(buf, data); err != nil {
		return
	}
	if data.ReturnCode == FAIL {
		err = errors.New(data.ReturnMsg)
		return
	}
	if data.ResultCode == FAIL {
		err = errors.New(data.ErrCodeDes)
		return
	}

	params := Xml2Map(buf)
	str := map2querystring(params) + "&key=" + apiKey

	var signStr string
	if data.SignType == SIGN_TYPE_MD5 {
		signStr = goo_utils.MD5([]byte(str))
	} else {
		signStr = goo_utils.HMacSha256([]byte(str), []byte(apiKey))
	}

	if strings.ToLower(signStr) != strings.ToLower(data.Sign) {
		err = errors.New("签名验证失败")
		return
	}

	return
}
