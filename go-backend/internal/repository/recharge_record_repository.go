package repository

import (
	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"gorm.io/gorm"
)

type RechargeRecordRepository struct {
	db *gorm.DB
}

func NewRechargeRecordRepository(db *gorm.DB) *RechargeRecordRepository {
	return &RechargeRecordRepository{db: db}
}

func (r *RechargeRecordRepository) Create(record *entity.RechargeRecord) error {
	return r.db.Create(record).Error
}

func (r *RechargeRecordRepository) GetByID(id int64) (*entity.RechargeRecord, error) {
	var record entity.RechargeRecord
	err := r.db.Where("id = ?", id).Take(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *RechargeRecordRepository) GetByPaymentID(paymentID string) (*entity.RechargeRecord, error) {
	var record entity.RechargeRecord
	err := r.db.Where("paymentId = ?", paymentID).Order("id desc").Take(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *RechargeRecordRepository) UpdateByID(record *entity.RechargeRecord) (bool, error) {
	result := r.db.Model(&entity.RechargeRecord{}).
		Where("id = ?", record.ID).
		Select("status", "paymentId", "updateTime").
		Updates(record)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *RechargeRecordRepository) ListByUserID(userID, pageNum, pageSize int64) (common.PageResponse[entity.RechargeRecord], error) {
	query := r.db.Model(&entity.RechargeRecord{}).Where("userId = ?", userID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return common.PageResponse[entity.RechargeRecord]{}, err
	}
	records := make([]entity.RechargeRecord, 0)
	offset := (pageNum - 1) * pageSize
	if err := query.Order("createTime desc").Offset(int(offset)).Limit(int(pageSize)).Find(&records).Error; err != nil {
		return common.PageResponse[entity.RechargeRecord]{}, err
	}
	return common.BuildPageResponse(records, pageNum, pageSize, total), nil
}
