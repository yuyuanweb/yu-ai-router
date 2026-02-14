package repository

import (
	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"gorm.io/gorm"
)

type ImageGenerationRecordRepository struct {
	db *gorm.DB
}

func NewImageGenerationRecordRepository(db *gorm.DB) *ImageGenerationRecordRepository {
	return &ImageGenerationRecordRepository{db: db}
}

func (r *ImageGenerationRecordRepository) Create(record *entity.ImageGenerationRecord) error {
	return r.db.Create(record).Error
}

func (r *ImageGenerationRecordRepository) ListByUserID(userID, pageNum, pageSize int64) (common.PageResponse[entity.ImageGenerationRecord], error) {
	query := r.db.Model(&entity.ImageGenerationRecord{}).Where("userId = ?", userID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return common.PageResponse[entity.ImageGenerationRecord]{}, err
	}
	records := make([]entity.ImageGenerationRecord, 0)
	offset := (pageNum - 1) * pageSize
	if err := query.Order("createTime desc").Offset(int(offset)).Limit(int(pageSize)).Find(&records).Error; err != nil {
		return common.PageResponse[entity.ImageGenerationRecord]{}, err
	}
	return common.BuildPageResponse(records, pageNum, pageSize, total), nil
}
