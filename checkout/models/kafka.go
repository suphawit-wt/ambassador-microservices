package models

type EmailKafkaMessage struct {
	Id                uint    `json:"id"`
	Code              string  `json:"code"`
	AmbassadorEmail   string  `json:"ambassador_email"`
	AdminRevenue      float64 `json:"admin_revenue"`
	AmbassadorRevenue float64 `json:"ambassador_revenue"`
}
