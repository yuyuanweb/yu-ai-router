package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type User struct {
	ID           int64     `json:"id,string" gorm:"column:id;primaryKey;autoIncrement"`
	UserAccount  string    `json:"userAccount" gorm:"column:userAccount"`
	UserPassword string    `json:"userPassword" gorm:"column:userPassword"`
	UserName     string    `json:"userName" gorm:"column:userName"`
	UserAvatar   string    `json:"userAvatar" gorm:"column:userAvatar"`
	UserProfile  string    `json:"userProfile" gorm:"column:userProfile"`
	UserRole     string    `json:"userRole" gorm:"column:userRole"`
	EditTime     time.Time `json:"editTime" gorm:"column:editTime"`
	CreateTime   time.Time `json:"createTime" gorm:"column:createTime"`
	UpdateTime   time.Time `json:"updateTime" gorm:"column:updateTime"`
	IsDelete     soft_delete.DeletedAt `json:"isDelete" gorm:"column:isDelete;softDelete:flag"`
}

func (User) TableName() string {
	return "user"
}
