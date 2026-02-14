package entity

import "time"

type BillingRecord struct {
	ID            int64     `json:"id,string" gorm:"column:id;primaryKey;autoIncrement"`
	UserID        int64     `json:"userId,string" gorm:"column:userId"`
	RequestLogID  *int64    `json:"requestLogId,string,omitempty" gorm:"column:requestLogId"`
	Amount        float64   `json:"amount" gorm:"column:amount"`
	BalanceBefore float64   `json:"balanceBefore" gorm:"column:balanceBefore"`
	BalanceAfter  float64   `json:"balanceAfter" gorm:"column:balanceAfter"`
	Description   string    `json:"description" gorm:"column:description"`
	BillingType   string    `json:"billingType" gorm:"column:billingType"`
	CreateTime    time.Time `json:"createTime" gorm:"column:createTime"`
}

func (BillingRecord) TableName() string {
	return "billing_record"
}
