package entity

import "time"

type RechargeRecord struct {
	ID            int64     `json:"id,string" gorm:"column:id;primaryKey;autoIncrement"`
	UserID        int64     `json:"userId,string" gorm:"column:userId"`
	Amount        float64   `json:"amount" gorm:"column:amount"`
	PaymentMethod string    `json:"paymentMethod" gorm:"column:paymentMethod"`
	PaymentID     string    `json:"paymentId" gorm:"column:paymentId"`
	Status        string    `json:"status" gorm:"column:status"`
	Description   string    `json:"description" gorm:"column:description"`
	CreateTime    time.Time `json:"createTime" gorm:"column:createTime"`
	UpdateTime    time.Time `json:"updateTime" gorm:"column:updateTime"`
}

func (RechargeRecord) TableName() string {
	return "recharge_record"
}
