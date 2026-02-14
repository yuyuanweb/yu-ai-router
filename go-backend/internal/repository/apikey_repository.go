package repository

import (
	"time"

	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ApiKeyRepository struct {
	db *gorm.DB
}

func NewApiKeyRepository(db *gorm.DB) *ApiKeyRepository {
	return &ApiKeyRepository{db: db}
}

func (r *ApiKeyRepository) Create(apiKey *entity.ApiKey) error {
	return r.db.
		Select("userId", "keyValue", "keyName", "status", "totalTokens", "lastUsedTime").
		Create(apiKey).Error
}

func (r *ApiKeyRepository) ListByUser(userID, pageNum, pageSize int64) ([]entity.ApiKey, int64, error) {
	query := r.db.Model(&entity.ApiKey{}).Where("userId = ?", userID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pageNum - 1) * pageSize
	list := make([]entity.ApiKey, 0)
	err := query.
		Order(clause.OrderByColumn{Column: clause.Column{Name: "createTime"}, Desc: true}).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *ApiKeyRepository) GetByID(id int64) (*entity.ApiKey, error) {
	var apiKey entity.ApiKey
	err := r.db.Where("id = ?", id).Take(&apiKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &apiKey, nil
}

func (r *ApiKeyRepository) GetByKeyValueActive(keyValue string) (*entity.ApiKey, error) {
	var apiKey entity.ApiKey
	err := r.db.
		Where("keyValue = ? AND status = ?", keyValue, "active").
		Take(&apiKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &apiKey, nil
}

func (r *ApiKeyRepository) RevokeByIDAndUser(id, userID int64) (bool, error) {
	result := r.db.Model(&entity.ApiKey{}).
		Where("id = ? AND userId = ?", id, userID).
		Updates(map[string]any{
			"status":     "revoked",
			"updateTime": time.Now(),
		})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *ApiKeyRepository) AddUsageStats(apiKeyID int64, tokens int) error {
	now := time.Now()
	return r.db.Model(&entity.ApiKey{}).
		Where("id = ?", apiKeyID).
		Updates(map[string]any{
			"totalTokens":  gorm.Expr("totalTokens + ?", tokens),
			"lastUsedTime": &now,
			"updateTime":   now,
		}).Error
}
