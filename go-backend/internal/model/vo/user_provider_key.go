package vo

import "time"

type UserProviderKeyVO struct {
	ID           int64     `json:"id,string"`
	ProviderID   int64     `json:"providerId,string"`
	ProviderName string    `json:"providerName"`
	APIKey       string    `json:"apiKey"`
	Status       string    `json:"status"`
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
}
