package repository

import (
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PluginRepository struct {
	db *gorm.DB
}

func NewPluginRepository(db *gorm.DB) *PluginRepository {
	return &PluginRepository{db: db}
}

func (r *PluginRepository) baseQuery() *gorm.DB {
	return r.db.Model(&entity.PluginConfig{})
}

func (r *PluginRepository) ListAll() ([]entity.PluginConfig, error) {
	list := make([]entity.PluginConfig, 0)
	err := r.baseQuery().
		Order(clause.OrderByColumn{Column: clause.Column{Name: "priority"}, Desc: true}).
		Find(&list).Error
	return list, err
}

func (r *PluginRepository) ListEnabled() ([]entity.PluginConfig, error) {
	list := make([]entity.PluginConfig, 0)
	err := r.baseQuery().
		Where("status = ?", "active").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "priority"}, Desc: true}).
		Find(&list).Error
	return list, err
}

func (r *PluginRepository) GetByID(id int64) (*entity.PluginConfig, error) {
	var plugin entity.PluginConfig
	err := r.baseQuery().Where("id = ?", id).Take(&plugin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &plugin, nil
}

func (r *PluginRepository) GetByKey(pluginKey string) (*entity.PluginConfig, error) {
	var plugin entity.PluginConfig
	err := r.baseQuery().Where("pluginKey = ?", pluginKey).Take(&plugin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &plugin, nil
}

func (r *PluginRepository) UpdateByID(id int64, fields map[string]any) (bool, error) {
	if len(fields) == 0 {
		return false, nil
	}
	result := r.baseQuery().Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
