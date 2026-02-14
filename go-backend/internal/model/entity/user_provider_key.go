package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type UserProviderKey struct {
	ID           int64                 `json:"id,string" gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int64                 `json:"userId,string" gorm:"column:userId"`
	ProviderID   int64                 `json:"providerId,string" gorm:"column:providerId"`
	ProviderName string                `json:"providerName" gorm:"column:providerName"`
	APIKey       string                `json:"apiKey" gorm:"column:apiKey"`
	Status       string                `json:"status" gorm:"column:status"`
	CreateTime   time.Time             `json:"createTime" gorm:"column:createTime"`
	UpdateTime   time.Time             `json:"updateTime" gorm:"column:updateTime"`
	IsDelete     soft_delete.DeletedAt `json:"isDelete" gorm:"column:isDelete;softDelete:flag"`
}

func (UserProviderKey) TableName() string {
	return "user_provider_key"
}
