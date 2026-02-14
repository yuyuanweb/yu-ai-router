package entity

import "time"

type RequestLog struct {
	ID               int64     `json:"id,string" gorm:"column:id;primaryKey;autoIncrement"`
	TraceID          string    `json:"traceId" gorm:"column:traceId"`
	UserID           *int64    `json:"userId,string,omitempty" gorm:"column:userId"`
	APIKeyID         *int64    `json:"apiKeyId,string,omitempty" gorm:"column:apiKeyId"`
	ModelID          *int64    `json:"modelId,string,omitempty" gorm:"column:modelId"`
	RequestModel     string    `json:"requestModel" gorm:"column:requestModel"`
	ModelName        string    `json:"modelName" gorm:"column:modelName"`
	RequestType      string    `json:"requestType" gorm:"column:requestType"`
	Source           string    `json:"source" gorm:"column:source"`
	PromptTokens     int       `json:"promptTokens" gorm:"column:promptTokens"`
	CompletionTokens int       `json:"completionTokens" gorm:"column:completionTokens"`
	TotalTokens      int       `json:"totalTokens" gorm:"column:totalTokens"`
	Cost             float64   `json:"cost" gorm:"column:cost"`
	Duration         int       `json:"duration" gorm:"column:duration"`
	Status           string    `json:"status" gorm:"column:status"`
	ErrorMessage     string    `json:"errorMessage" gorm:"column:errorMessage"`
	ErrorCode        string    `json:"errorCode" gorm:"column:errorCode"`
	RoutingStrategy  string    `json:"routingStrategy" gorm:"column:routingStrategy"`
	IsFallback       int       `json:"isFallback" gorm:"column:isFallback"`
	ClientIP         string    `json:"clientIp" gorm:"column:clientIp"`
	UserAgent        string    `json:"userAgent" gorm:"column:userAgent"`
	CreateTime       time.Time `json:"createTime" gorm:"column:createTime"`
	UpdateTime       time.Time `json:"updateTime" gorm:"column:updateTime"`
}

func (RequestLog) TableName() string {
	return "request_log"
}
