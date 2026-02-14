package vo

import "time"

type ModelVO struct {
	ID                  int64     `json:"id,string"`
	ProviderID          int64     `json:"providerId,string"`
	ProviderName        string    `json:"providerName"`
	ProviderDisplayName string    `json:"providerDisplayName"`
	ModelKey            string    `json:"modelKey"`
	ModelName           string    `json:"modelName"`
	ModelType           string    `json:"modelType"`
	Description         string    `json:"description"`
	ContextLength       int       `json:"contextLength"`
	InputPrice          float64   `json:"inputPrice"`
	OutputPrice         float64   `json:"outputPrice"`
	Status              string    `json:"status"`
	HealthStatus        string    `json:"healthStatus"`
	AvgLatency          int       `json:"avgLatency"`
	SuccessRate         float64   `json:"successRate"`
	Priority            int       `json:"priority"`
	DefaultTimeout      int       `json:"defaultTimeout"`
	SupportReasoning    int       `json:"supportReasoning"`
	Capabilities        string    `json:"capabilities"`
	CreateTime          time.Time `json:"createTime"`
	UpdateTime          time.Time `json:"updateTime"`
}
