package vo

import "time"

type ApiKeyVO struct {
	ID           int64      `json:"id,string"`
	KeyValue     string     `json:"keyValue"`
	KeyName      string     `json:"keyName"`
	Status       string     `json:"status"`
	TotalTokens  int64      `json:"totalTokens,string"`
	LastUsedTime *time.Time `json:"lastUsedTime"`
	CreateTime   time.Time  `json:"createTime"`
}
