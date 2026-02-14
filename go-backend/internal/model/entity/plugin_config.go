package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type PluginConfig struct {
	ID          int64                 `json:"id,string" gorm:"column:id;primaryKey;autoIncrement"`
	PluginKey   string                `json:"pluginKey" gorm:"column:pluginKey"`
	PluginName  string                `json:"pluginName" gorm:"column:pluginName"`
	PluginType  string                `json:"pluginType" gorm:"column:pluginType"`
	Description string                `json:"description" gorm:"column:description"`
	Config      string                `json:"config" gorm:"column:config"`
	Status      string                `json:"status" gorm:"column:status"`
	Priority    int                   `json:"priority" gorm:"column:priority"`
	CreateTime  time.Time             `json:"createTime" gorm:"column:createTime"`
	UpdateTime  time.Time             `json:"updateTime" gorm:"column:updateTime"`
	IsDelete    soft_delete.DeletedAt `json:"isDelete" gorm:"column:isDelete;softDelete:flag"`
}

func (PluginConfig) TableName() string {
	return "plugin_config"
}
