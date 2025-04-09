package model

type AltreRequest struct {
	Host         string `json:"host"`
	Ip           string `json:"ip"`
	Type         string `json:"type"` // 告警类型  5 1 0
	LastType     string `json:"last_type"`
	MetricCnName string `json:"metric_cn_name"`
	MetricValue  string `json:"metric_value"`
	MetricUnit   string `json:"metric_unit"`
	Threshold5   string `json:"threshold5"`
	Threshold4   string `json:"threshold4"`
}
