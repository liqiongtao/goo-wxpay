package goo_wxpay

type Config struct {
	MchId           string `yaml:"mch_id" json:"mch_id"`
	ApiKey          string `yaml:"api_key" json:"api_key"`
	NotifyUrl       string `yaml:"notify_url" json:"notify_url"`
	NotifyRefundUrl string `yaml:"notify_refund_url" json:"notify_refund_url"`
}
