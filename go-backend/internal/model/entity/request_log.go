package entity

import "time"

type RequestLog struct {
	ID               int64     `json:"id,string" gorm:"column:id;primaryKey;autoIncrement"`
	UserID           *int64    `json:"userId,string,omitempty" gorm:"column:userId"`
	APIKeyID         *int64    `json:"apiKeyId,string,omitempty" gorm:"column:apiKeyId"`
	ModelName        string    `json:"modelName" gorm:"column:modelName"`
	PromptTokens     int       `json:"promptTokens" gorm:"column:promptTokens"`
	CompletionTokens int       `json:"completionTokens" gorm:"column:completionTokens"`
	TotalTokens      int       `json:"totalTokens" gorm:"column:totalTokens"`
	Duration         int       `json:"duration" gorm:"column:duration"`
	Status           string    `json:"status" gorm:"column:status"`
	ErrorMessage     string    `json:"errorMessage" gorm:"column:errorMessage"`
	CreateTime       time.Time `json:"createTime" gorm:"column:createTime"`
	UpdateTime       time.Time `json:"updateTime" gorm:"column:updateTime"`
}

func (RequestLog) TableName() string {
	return "request_log"
}
