package entity

import "time"

type ImageGenerationRecord struct {
	ID            int64     `json:"id,string" gorm:"column:id;primaryKey;autoIncrement"`
	UserID        int64     `json:"userId,string" gorm:"column:userId"`
	APIKeyID      *int64    `json:"apiKeyId,string,omitempty" gorm:"column:apiKeyId"`
	ModelID       int64     `json:"modelId,string" gorm:"column:modelId"`
	ModelKey      string    `json:"modelKey" gorm:"column:modelKey"`
	Prompt        string    `json:"prompt" gorm:"column:prompt"`
	RevisedPrompt string    `json:"revisedPrompt" gorm:"column:revisedPrompt"`
	ImageURL      string    `json:"imageUrl" gorm:"column:imageUrl"`
	ImageData     string    `json:"imageData" gorm:"column:imageData"`
	Size          string    `json:"size" gorm:"column:size"`
	Quality       string    `json:"quality" gorm:"column:quality"`
	Status        string    `json:"status" gorm:"column:status"`
	Cost          float64   `json:"cost" gorm:"column:cost"`
	Duration      int       `json:"duration" gorm:"column:duration"`
	ErrorMessage  string    `json:"errorMessage" gorm:"column:errorMessage"`
	ClientIP      string    `json:"clientIp" gorm:"column:clientIp"`
	CreateTime    time.Time `json:"createTime" gorm:"column:createTime"`
}

func (ImageGenerationRecord) TableName() string {
	return "image_generation_record"
}
