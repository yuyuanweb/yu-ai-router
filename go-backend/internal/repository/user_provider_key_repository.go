package repository

import (
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"gorm.io/gorm"
)

type UserProviderKeyRepository struct {
	db *gorm.DB
}

func NewUserProviderKeyRepository(db *gorm.DB) *UserProviderKeyRepository {
	return &UserProviderKeyRepository{db: db}
}

func (r *UserProviderKeyRepository) baseQuery() *gorm.DB {
	return r.db.Model(&entity.UserProviderKey{})
}

func (r *UserProviderKeyRepository) Create(record *entity.UserProviderKey) error {
	return r.db.
		Select("userId", "providerId", "providerName", "apiKey", "status").
		Create(record).Error
}

func (r *UserProviderKeyRepository) GetByID(id int64) (*entity.UserProviderKey, error) {
	var record entity.UserProviderKey
	err := r.baseQuery().Where("id = ?", id).Take(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *UserProviderKeyRepository) GetByUserAndProvider(userID, providerID int64) (*entity.UserProviderKey, error) {
	var record entity.UserProviderKey
	err := r.baseQuery().
		Where("userId = ? AND providerId = ?", userID, providerID).
		Take(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *UserProviderKeyRepository) GetActiveByUserAndProvider(userID, providerID int64) (*entity.UserProviderKey, error) {
	var record entity.UserProviderKey
	err := r.baseQuery().
		Where("userId = ? AND providerId = ? AND status = ?", userID, providerID, "active").
		Take(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *UserProviderKeyRepository) ListByUser(userID int64) ([]entity.UserProviderKey, error) {
	list := make([]entity.UserProviderKey, 0)
	err := r.baseQuery().
		Where("userId = ?", userID).
		Order("createTime DESC").
		Find(&list).Error
	return list, err
}

func (r *UserProviderKeyRepository) UpdateByID(id int64, fields map[string]any) (bool, error) {
	if len(fields) == 0 {
		return false, nil
	}
	result := r.baseQuery().Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *UserProviderKeyRepository) SoftDeleteByID(id int64) (bool, error) {
	result := r.baseQuery().Where("id = ?", id).Delete(&entity.UserProviderKey{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
