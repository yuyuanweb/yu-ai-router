package repository

import (
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RequestLogRepository struct {
	db *gorm.DB
}

func NewRequestLogRepository(db *gorm.DB) *RequestLogRepository {
	return &RequestLogRepository{db: db}
}

func (r *RequestLogRepository) Create(log *entity.RequestLog) error {
	return r.db.
		Select(
			"userId",
			"apiKeyId",
			"modelName",
			"promptTokens",
			"completionTokens",
			"totalTokens",
			"duration",
			"status",
			"errorMessage",
		).
		Create(log).Error
}

func (r *RequestLogRepository) ListUserLogs(userID int64, limit int) ([]entity.RequestLog, error) {
	list := make([]entity.RequestLog, 0)
	err := r.db.Model(&entity.RequestLog{}).
		Where("userId = ?", userID).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "createTime"}, Desc: true}).
		Limit(limit).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *RequestLogRepository) CountUserTokens(userID int64) (int64, error) {
	type tokenSum struct {
		Total int64 `gorm:"column:total"`
	}
	var result tokenSum
	err := r.db.Model(&entity.RequestLog{}).
		Select("COALESCE(SUM(totalTokens), 0) AS total").
		Where("userId = ? AND status = ?", userID, "success").
		Take(&result).Error
	if err != nil {
		return 0, err
	}
	return result.Total, nil
}
