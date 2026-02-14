package vo

import "time"

type ProviderVO struct {
	ID           int64     `json:"id,string"`
	ProviderName string    `json:"providerName"`
	DisplayName  string    `json:"displayName"`
	BaseURL      string    `json:"baseUrl"`
	Status       string    `json:"status"`
	HealthStatus string    `json:"healthStatus"`
	AvgLatency   int       `json:"avgLatency"`
	SuccessRate  float64   `json:"successRate"`
	Priority     int       `json:"priority"`
	Config       string    `json:"config"`
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
}
