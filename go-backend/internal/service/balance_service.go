package service

import (
	"fmt"
	"time"

	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BalanceService struct {
	userRepo             *repository.UserRepository
	billingRecordService *BillingRecordService
}

func NewBalanceService(
	userRepo *repository.UserRepository,
	billingRecordService *BillingRecordService,
) *BalanceService {
	return &BalanceService{
		userRepo:             userRepo,
		billingRecordService: billingRecordService,
	}
}

func (s *BalanceService) CheckBalance(userID int64, amount float64) (bool, error) {
	if userID <= 0 || amount <= 0 {
		return false, errno.New(errno.ParamsError)
	}
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if user == nil {
		return false, errno.New(errno.NotFoundError)
	}
	return user.Balance >= amount, nil
}

func (s *BalanceService) DeductBalance(userID int64, amount float64, requestLogID *int64, description string) error {
	if userID <= 0 || amount <= 0 {
		return errno.New(errno.ParamsError)
	}
	return s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		var user entity.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&entity.User{}).
			Where("id = ?", userID).
			Take(&user).Error; err != nil {
			return errno.New(errno.SystemError)
		}
		if user.Balance < amount {
			return errno.NewWithMessage(errno.ForbiddenError, fmt.Sprintf("余额不足，当前余额：%.4f，需要：%.4f", user.Balance, amount))
		}
		newBalance := user.Balance - amount
		if err := tx.Model(&entity.User{}).
			Where("id = ?", userID).
			Update("balance", newBalance).Error; err != nil {
			return errno.New(errno.SystemError)
		}
		if description == "" {
			description = "API调用消费"
		}
		record := &entity.BillingRecord{
			UserID:        userID,
			RequestLogID:  requestLogID,
			Amount:        amount,
			BalanceBefore: user.Balance,
			BalanceAfter:  newBalance,
			Description:   description,
			BillingType:   "api_call",
			CreateTime:    time.Now(),
		}
		if err := tx.Create(record).Error; err != nil {
			return errno.New(errno.SystemError)
		}
		return nil
	})
}

func (s *BalanceService) AddBalance(userID int64, amount float64, description string) error {
	if userID <= 0 || amount <= 0 {
		return errno.New(errno.ParamsError)
	}
	return s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		var user entity.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&entity.User{}).
			Where("id = ?", userID).
			Take(&user).Error; err != nil {
			return errno.New(errno.SystemError)
		}
		newBalance := user.Balance + amount
		if err := tx.Model(&entity.User{}).
			Where("id = ?", userID).
			Update("balance", newBalance).Error; err != nil {
			return errno.New(errno.SystemError)
		}
		if description == "" {
			description = "账户充值"
		}
		record := &entity.BillingRecord{
			UserID:        userID,
			Amount:        amount,
			BalanceBefore: user.Balance,
			BalanceAfter:  newBalance,
			Description:   description,
			BillingType:   "recharge",
			CreateTime:    time.Now(),
		}
		if err := tx.Create(record).Error; err != nil {
			return errno.New(errno.SystemError)
		}
		return nil
	})
}

func (s *BalanceService) GetUserBalance(userID int64) (float64, error) {
	if userID <= 0 {
		return 0, errno.New(errno.ParamsError)
	}
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return 0, errno.New(errno.SystemError)
	}
	if user == nil {
		return 0, errno.New(errno.NotFoundError)
	}
	return user.Balance, nil
}
