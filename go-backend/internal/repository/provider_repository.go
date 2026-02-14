package repository

import (
	"strings"

	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var allowedProviderSortFields = map[string]string{
	"id":           "id",
	"providerName": "providerName",
	"displayName":  "displayName",
	"status":       "status",
	"healthStatus": "healthStatus",
	"priority":     "priority",
	"createTime":   "createTime",
	"updateTime":   "updateTime",
}

type ProviderRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) *ProviderRepository {
	return &ProviderRepository{db: db}
}

func (r *ProviderRepository) baseQuery() *gorm.DB {
	return r.db.Model(&entity.ModelProvider{})
}

func (r *ProviderRepository) Create(provider *entity.ModelProvider) error {
	return r.db.
		Select("providerName", "displayName", "baseUrl", "apiKey", "status", "priority", "config").
		Create(provider).Error
}

func (r *ProviderRepository) GetByID(id int64) (*entity.ModelProvider, error) {
	var provider entity.ModelProvider
	err := r.baseQuery().Where("id = ?", id).Take(&provider).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &provider, nil
}

func (r *ProviderRepository) GetByProviderName(providerName string) (*entity.ModelProvider, error) {
	var provider entity.ModelProvider
	err := r.baseQuery().Where("providerName = ?", providerName).Take(&provider).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &provider, nil
}

func (r *ProviderRepository) ListByIDs(ids []int64) ([]entity.ModelProvider, error) {
	if len(ids) == 0 {
		return make([]entity.ModelProvider, 0), nil
	}
	list := make([]entity.ModelProvider, 0, len(ids))
	err := r.baseQuery().Where("id IN ?", ids).Find(&list).Error
	return list, err
}

func (r *ProviderRepository) UpdateByID(id int64, fields map[string]any) (bool, error) {
	if len(fields) == 0 {
		return false, nil
	}
	result := r.baseQuery().Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *ProviderRepository) SoftDeleteByID(id int64) (bool, error) {
	result := r.baseQuery().Where("id = ?", id).Delete(&entity.ModelProvider{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *ProviderRepository) ListByPage(req dto.ProviderQueryRequest) ([]entity.ModelProvider, int64, error) {
	query := r.baseQuery()
	if strings.TrimSpace(req.ProviderName) != "" {
		query = query.Where("providerName LIKE ?", "%"+req.ProviderName+"%")
	}
	if strings.TrimSpace(req.DisplayName) != "" {
		query = query.Where("displayName LIKE ?", "%"+req.DisplayName+"%")
	}
	if strings.TrimSpace(req.Status) != "" {
		query = query.Where("status = ?", req.Status)
	}
	if strings.TrimSpace(req.HealthStatus) != "" {
		query = query.Where("healthStatus = ?", req.HealthStatus)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sortField, ok := allowedProviderSortFields[req.SortField]; ok {
		query = query.Order(clause.OrderByColumn{
			Column: clause.Column{Name: sortField},
			Desc:   req.SortOrder != "ascend",
		})
	} else {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: "priority"}, Desc: true}).
			Order(clause.OrderByColumn{Column: clause.Column{Name: "createTime"}, Desc: true})
	}

	offset := (req.PageNum - 1) * req.PageSize
	list := make([]entity.ModelProvider, 0)
	err := query.Offset(int(offset)).Limit(int(req.PageSize)).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *ProviderRepository) ListAll() ([]entity.ModelProvider, error) {
	list := make([]entity.ModelProvider, 0)
	err := r.baseQuery().
		Order(clause.OrderByColumn{Column: clause.Column{Name: "priority"}, Desc: true}).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "createTime"}, Desc: true}).
		Find(&list).Error
	return list, err
}

func (r *ProviderRepository) ListActive() ([]entity.ModelProvider, error) {
	list := make([]entity.ModelProvider, 0)
	err := r.baseQuery().
		Where("status = ?", "active").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "priority"}, Desc: true}).
		Find(&list).Error
	return list, err
}

func (r *ProviderRepository) ListHealthy() ([]entity.ModelProvider, error) {
	list := make([]entity.ModelProvider, 0)
	err := r.baseQuery().
		Where("status = ?", "active").
		Where("healthStatus IN ?", []string{"healthy", "degraded"}).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "priority"}, Desc: true}).
		Find(&list).Error
	return list, err
}
