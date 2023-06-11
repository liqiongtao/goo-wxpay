package goo_wxpay

import (
	"fmt"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"strconv"
	"testing"
)

func TestUnifiedOrder(t *testing.T) {
	req := &UnifiedOrderRequest{
		Appid:          "",
		MchId:          "",
		Body:           "积分充值",
		OutTradeNo:     strconv.FormatInt(goo_utils.GenId(), 10),
		TotalFee:       1,
		SpbillCreateIp: "192.168.2.100",
		NotifyUrl:      "",
		TradeType:      TRADE_TYPE_NATIVE,
	}

	resp, err := UnifiedOrder(req, "")

	fmt.Println(resp.CodeUrl, err)
}
