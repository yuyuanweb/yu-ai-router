package repository

import (
	"strings"

	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var allowedSortFields = map[string]string{
	"id":          "id",
	"userAccount": "userAccount",
	"userName":    "userName",
	"userRole":    "userRole",
	"createTime":  "createTime",
	"updateTime":  "updateTime",
	"editTime":    "editTime",
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) DB() *gorm.DB {
	return r.db
}

func (r *UserRepository) baseQuery() *gorm.DB {
	return r.db.Model(&entity.User{})
}

func (r *UserRepository) Create(user *entity.User) (int64, error) {
	if err := r.db.
		Select("userAccount", "userPassword", "userName", "userAvatar", "userProfile", "userRole", "userStatus", "balance").
		Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *UserRepository) CountByAccount(userAccount string) (int64, error) {
	var count int64
	err := r.baseQuery().
		Where("userAccount = ?", userAccount).
		Count(&count).Error
	return count, err
}

func (r *UserRepository) GetByAccountAndPassword(userAccount, userPassword string) (*entity.User, error) {
	var user entity.User
	err := r.baseQuery().
		Where("userAccount = ? AND userPassword = ?", userAccount, userPassword).
		Take(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id int64) (*entity.User, error) {
	var user entity.User
	err := r.baseQuery().Where("id = ?", id).Take(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateByID(user *entity.User) (bool, error) {
	result := r.baseQuery().
		Where("id = ?", user.ID).
		Select("userName", "userAvatar", "userProfile", "userRole").
		Updates(user)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *UserRepository) UpdateQuota(userID, tokenQuota int64) (bool, error) {
	result := r.baseQuery().
		Where("id = ?", userID).
		Select("tokenQuota").
		Updates(&entity.User{TokenQuota: tokenQuota})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *UserRepository) ResetUsedTokens(userID int64) (bool, error) {
	result := r.baseQuery().
		Where("id = ?", userID).
		Select("usedTokens").
		Updates(&entity.User{UsedTokens: 0})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *UserRepository) UpdateStatus(userID int64, userStatus string) (bool, error) {
	result := r.baseQuery().
		Where("id = ?", userID).
		Select("userStatus").
		Updates(&entity.User{UserStatus: userStatus})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *UserRepository) UpdateBalance(userID int64, balance float64) (bool, error) {
	result := r.baseQuery().
		Where("id = ?", userID).
		Select("balance").
		Updates(&entity.User{Balance: balance})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *UserRepository) AddUsedTokens(userID int64, tokens int64) (bool, error) {
	result := r.baseQuery().
		Where("id = ?", userID).
		UpdateColumn("usedTokens", gorm.Expr("usedTokens + ?", tokens))
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *UserRepository) SoftDeleteByID(id int64) (bool, error) {
	result := r.baseQuery().
		Where("id = ?", id).
		Delete(&entity.User{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *UserRepository) ListByPage(userQueryRequest dto.UserQueryRequest) ([]entity.User, int64, error) {
	query := r.baseQuery()

	if userQueryRequest.ID != nil && *userQueryRequest.ID > 0 {
		query = query.Where("id = ?", userQueryRequest.ID.Int64())
	}
	if strings.TrimSpace(userQueryRequest.UserRole) != "" {
		query = query.Where("userRole = ?", userQueryRequest.UserRole)
	}
	if strings.TrimSpace(userQueryRequest.UserAccount) != "" {
		query = query.Where("userAccount LIKE ?", "%"+userQueryRequest.UserAccount+"%")
	}
	if strings.TrimSpace(userQueryRequest.UserName) != "" {
		query = query.Where("userName LIKE ?", "%"+userQueryRequest.UserName+"%")
	}
	if strings.TrimSpace(userQueryRequest.UserProfile) != "" {
		query = query.Where("userProfile LIKE ?", "%"+userQueryRequest.UserProfile+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sortField, ok := allowedSortFields[userQueryRequest.SortField]; ok {
		query = query.Order(clause.OrderByColumn{
			Column: clause.Column{Name: sortField},
			Desc:   userQueryRequest.SortOrder != "ascend",
		})
	}

	offset := (userQueryRequest.PageNum - 1) * userQueryRequest.PageSize
	users := make([]entity.User, 0)
	err := query.Offset(int(offset)).Limit(int(userQueryRequest.PageSize)).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
