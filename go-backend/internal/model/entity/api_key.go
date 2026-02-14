package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type ApiKey struct {
	ID           int64                  `json:"id,string" gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int64                  `json:"userId,string" gorm:"column:userId"`
	KeyValue     string                 `json:"keyValue" gorm:"column:keyValue"`
	KeyName      string                 `json:"keyName" gorm:"column:keyName"`
	Status       string                 `json:"status" gorm:"column:status"`
	TotalTokens  int64                  `json:"totalTokens,string" gorm:"column:totalTokens"`
	LastUsedTime *time.Time             `json:"lastUsedTime" gorm:"column:lastUsedTime"`
	CreateTime   time.Time              `json:"createTime" gorm:"column:createTime"`
	UpdateTime   time.Time              `json:"updateTime" gorm:"column:updateTime"`
	IsDelete     soft_delete.DeletedAt  `json:"isDelete" gorm:"column:isDelete;softDelete:flag"`
}

func (ApiKey) TableName() string {
	return "api_key"
}
