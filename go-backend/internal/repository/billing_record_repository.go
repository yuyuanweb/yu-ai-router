package repository

import (
	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"gorm.io/gorm"
)

type BillingRecordRepository struct {
	db *gorm.DB
}

func NewBillingRecordRepository(db *gorm.DB) *BillingRecordRepository {
	return &BillingRecordRepository{db: db}
}

func (r *BillingRecordRepository) Create(record *entity.BillingRecord) error {
	return r.db.Create(record).Error
}

func (r *BillingRecordRepository) ListByUserID(userID, pageNum, pageSize int64) (common.PageResponse[entity.BillingRecord], error) {
	query := r.db.Model(&entity.BillingRecord{}).Where("userId = ?", userID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return common.PageResponse[entity.BillingRecord]{}, err
	}
	records := make([]entity.BillingRecord, 0)
	offset := (pageNum - 1) * pageSize
	if err := query.Order("createTime desc").Offset(int(offset)).Limit(int(pageSize)).Find(&records).Error; err != nil {
		return common.PageResponse[entity.BillingRecord]{}, err
	}
	return common.BuildPageResponse(records, pageNum, pageSize, total), nil
}

func (r *BillingRecordRepository) SumAmountByType(userID int64, billingType string) (float64, error) {
	type resultRow struct {
		Total float64 `gorm:"column:total"`
	}
	var row resultRow
	if err := r.db.Model(&entity.BillingRecord{}).
		Select("COALESCE(SUM(amount), 0) AS total").
		Where("userId = ? AND billingType = ?", userID, billingType).
		Take(&row).Error; err != nil {
		return 0, err
	}
	return row.Total, nil
}
