package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type ModelProvider struct {
	ID           int64                 `json:"id,string" gorm:"column:id;primaryKey;autoIncrement"`
	ProviderName string                `json:"providerName" gorm:"column:providerName"`
	DisplayName  string                `json:"displayName" gorm:"column:displayName"`
	BaseURL      string                `json:"baseUrl" gorm:"column:baseUrl"`
	APIKey       string                `json:"apiKey" gorm:"column:apiKey"`
	Status       string                `json:"status" gorm:"column:status"`
	HealthStatus string                `json:"healthStatus" gorm:"column:healthStatus"`
	AvgLatency   int                   `json:"avgLatency" gorm:"column:avgLatency"`
	SuccessRate  float64               `json:"successRate" gorm:"column:successRate"`
	Priority     int                   `json:"priority" gorm:"column:priority"`
	Config       string                `json:"config" gorm:"column:config"`
	CreateTime   time.Time             `json:"createTime" gorm:"column:createTime"`
	UpdateTime   time.Time             `json:"updateTime" gorm:"column:updateTime"`
	IsDelete     soft_delete.DeletedAt `json:"isDelete" gorm:"column:isDelete;softDelete:flag"`
}

func (ModelProvider) TableName() string {
	return "model_provider"
}
