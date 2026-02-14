package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Model struct {
	ID               int64                 `json:"id,string" gorm:"column:id;primaryKey;autoIncrement"`
	ProviderID       int64                 `json:"providerId,string" gorm:"column:providerId"`
	ModelKey         string                `json:"modelKey" gorm:"column:modelKey"`
	ModelName        string                `json:"modelName" gorm:"column:modelName"`
	ModelType        string                `json:"modelType" gorm:"column:modelType"`
	Description      string                `json:"description" gorm:"column:description"`
	ContextLength    int                   `json:"contextLength" gorm:"column:contextLength"`
	InputPrice       float64               `json:"inputPrice" gorm:"column:inputPrice"`
	OutputPrice      float64               `json:"outputPrice" gorm:"column:outputPrice"`
	Status           string                `json:"status" gorm:"column:status"`
	HealthStatus     string                `json:"healthStatus" gorm:"column:healthStatus"`
	AvgLatency       int                   `json:"avgLatency" gorm:"column:avgLatency"`
	SuccessRate      float64               `json:"successRate" gorm:"column:successRate"`
	Score            float64               `json:"score" gorm:"column:score"`
	Priority         int                   `json:"priority" gorm:"column:priority"`
	DefaultTimeout   int                   `json:"defaultTimeout" gorm:"column:defaultTimeout"`
	SupportReasoning int                   `json:"supportReasoning" gorm:"column:supportReasoning"`
	Capabilities     string                `json:"capabilities" gorm:"column:capabilities"`
	CreateTime       time.Time             `json:"createTime" gorm:"column:createTime"`
	UpdateTime       time.Time             `json:"updateTime" gorm:"column:updateTime"`
	IsDelete         soft_delete.DeletedAt `json:"isDelete" gorm:"column:isDelete;softDelete:flag"`
}

func (Model) TableName() string {
	return "model"
}
