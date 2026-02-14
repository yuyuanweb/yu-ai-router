package repository

import (
	"math"
	"strings"

	"github.com/yupi/airouter/go-backend/internal/constant"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var allowedModelSortFields = map[string]string{
	"id":             "id",
	"providerId":     "providerId",
	"modelKey":       "modelKey",
	"modelName":      "modelName",
	"modelType":      "modelType",
	"status":         "status",
	"healthStatus":   "healthStatus",
	"priority":       "priority",
	"defaultTimeout": "defaultTimeout",
	"createTime":     "createTime",
	"updateTime":     "updateTime",
}

type ModelRepository struct {
	db *gorm.DB
}

func NewModelRepository(db *gorm.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

func (r *ModelRepository) baseQuery() *gorm.DB {
	return r.db.Model(&entity.Model{})
}

func (r *ModelRepository) Create(model *entity.Model) error {
	return r.db.
		Select(
			"providerId",
			"modelKey",
			"modelName",
			"modelType",
			"description",
			"contextLength",
			"inputPrice",
			"outputPrice",
			"status",
			"priority",
			"defaultTimeout",
			"capabilities",
		).
		Create(model).Error
}

func (r *ModelRepository) GetByID(id int64) (*entity.Model, error) {
	var model entity.Model
	err := r.baseQuery().Where("id = ?", id).Take(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

func (r *ModelRepository) GetByModelKey(modelKey string) (*entity.Model, error) {
	var model entity.Model
	err := r.baseQuery().Where("modelKey = ?", modelKey).Take(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

func (r *ModelRepository) GetRoutableByModelKey(modelKey, modelType string) (*entity.Model, error) {
	query := r.baseQuery().
		Where("modelKey = ?", modelKey).
		Where("status = ?", constant.ModelStatusActive).
		Where("healthStatus IN ?", []string{constant.HealthStatusHealthy, constant.HealthStatusDegraded, constant.HealthStatusUnknown})
	if strings.TrimSpace(modelType) != "" {
		query = query.Where("modelType = ?", modelType)
	}

	var model entity.Model
	err := query.Take(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

func (r *ModelRepository) ListRoutable(modelType string) ([]entity.Model, error) {
	list := make([]entity.Model, 0)
	query := r.baseQuery().
		Where("status = ?", constant.ModelStatusActive).
		Where("healthStatus IN ?", []string{constant.HealthStatusHealthy, constant.HealthStatusDegraded, constant.HealthStatusUnknown})
	if strings.TrimSpace(modelType) != "" {
		query = query.Where("modelType = ?", modelType)
	}
	err := query.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ModelRepository) UpdateByID(id int64, fields map[string]any) (bool, error) {
	if len(fields) == 0 {
		return false, nil
	}
	result := r.baseQuery().Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *ModelRepository) SoftDeleteByID(id int64) (bool, error) {
	result := r.baseQuery().Where("id = ?", id).Delete(&entity.Model{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *ModelRepository) ListByPage(req dto.ModelQueryRequest) ([]entity.Model, int64, error) {
	query := r.baseQuery()
	if req.ProviderID != nil && req.ProviderID.Int64() > 0 {
		query = query.Where("providerId = ?", req.ProviderID.Int64())
	}
	if strings.TrimSpace(req.ModelKey) != "" {
		query = query.Where("modelKey LIKE ?", "%"+req.ModelKey+"%")
	}
	if strings.TrimSpace(req.ModelName) != "" {
		query = query.Where("modelName LIKE ?", "%"+req.ModelName+"%")
	}
	if strings.TrimSpace(req.ModelType) != "" {
		query = query.Where("modelType = ?", req.ModelType)
	}
	if strings.TrimSpace(req.Status) != "" {
		query = query.Where("status = ?", req.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sortField, ok := allowedModelSortFields[req.SortField]; ok {
		query = query.Order(clause.OrderByColumn{
			Column: clause.Column{Name: sortField},
			Desc:   req.SortOrder != "ascend",
		})
	} else {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: "priority"}, Desc: true}).
			Order(clause.OrderByColumn{Column: clause.Column{Name: "createTime"}, Desc: true})
	}

	offset := (req.PageNum - 1) * req.PageSize
	list := make([]entity.Model, 0)
	err := query.Offset(int(offset)).Limit(int(req.PageSize)).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *ModelRepository) ListAll() ([]entity.Model, error) {
	list := make([]entity.Model, 0)
	err := r.baseQuery().
		Order(clause.OrderByColumn{Column: clause.Column{Name: "priority"}, Desc: true}).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "createTime"}, Desc: true}).
		Find(&list).Error
	return list, err
}

func (r *ModelRepository) ListActive() ([]entity.Model, error) {
	list := make([]entity.Model, 0)
	err := r.baseQuery().
		Where("status = ?", "active").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "priority"}, Desc: true}).
		Find(&list).Error
	return list, err
}

func (r *ModelRepository) ListActiveByProviderID(providerID int64) ([]entity.Model, error) {
	list := make([]entity.Model, 0)
	err := r.baseQuery().
		Where("providerId = ? AND status = ?", providerID, "active").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "priority"}, Desc: true}).
		Find(&list).Error
	return list, err
}

func (r *ModelRepository) ListActiveByType(modelType string) ([]entity.Model, error) {
	list := make([]entity.Model, 0)
	err := r.baseQuery().
		Where("modelType = ? AND status = ?", modelType, "active").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "priority"}, Desc: true}).
		Find(&list).Error
	return list, err
}

func CostValue(model entity.Model) float64 {
	return model.InputPrice + model.OutputPrice
}

func LatencyOrderValue(model entity.Model) float64 {
	if model.AvgLatency <= 0 {
		return 999999
	}
	return float64(model.AvgLatency)
}

func ScoreOrderValue(model entity.Model) float64 {
	if model.Score <= 0 {
		return math.MaxFloat64
	}
	return model.Score
}
