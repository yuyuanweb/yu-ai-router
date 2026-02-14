package service

import (
	"time"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RechargeService struct {
	repo           *repository.RechargeRecordRepository
	balanceService *BalanceService
}

func NewRechargeService(
	repo *repository.RechargeRecordRepository,
	balanceService *BalanceService,
) *RechargeService {
	return &RechargeService{
		repo:           repo,
		balanceService: balanceService,
	}
}

func (s *RechargeService) CreateRechargeRecord(userID int64, amount float64, paymentMethod string) (*entity.RechargeRecord, error) {
	record := &entity.RechargeRecord{
		UserID:        userID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Status:        "pending",
		Description:   "账户充值",
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	if err := s.repo.Create(record); err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return record, nil
}

func (s *RechargeService) UpdateRechargeStatus(recordID int64, status, paymentID string) (bool, error) {
	record, err := s.repo.GetByID(recordID)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if record == nil {
		return false, errno.New(errno.NotFoundError)
	}
	record.Status = status
	record.PaymentID = paymentID
	record.UpdateTime = time.Now()
	return s.repo.UpdateByID(record)
}

func (s *RechargeService) GetByPaymentID(paymentID string) (*entity.RechargeRecord, error) {
	return s.repo.GetByPaymentID(paymentID)
}

func (s *RechargeService) ListUserRechargeRecords(userID, pageNum, pageSize int64) (common.PageResponse[entity.RechargeRecord], error) {
	if userID <= 0 {
		return common.PageResponse[entity.RechargeRecord]{}, errno.New(errno.ParamsError)
	}
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.ListByUserID(userID, pageNum, pageSize)
}

func (s *RechargeService) CompleteRecharge(recordID int64, paymentID string) (bool, error) {
	if paymentID != "" {
		existed, err := s.repo.GetByPaymentID(paymentID)
		if err != nil {
			return false, errno.New(errno.SystemError)
		}
		if existed != nil && existed.Status == "success" {
			return true, nil
		}
	}

	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		var record entity.RechargeRecord
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&entity.RechargeRecord{}).
			Where("id = ?", recordID).
			Take(&record).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errno.New(errno.NotFoundError)
			}
			return errno.New(errno.SystemError)
		}
		if record.Status == "success" {
			return nil
		}

		record.Status = "success"
		record.PaymentID = paymentID
		record.UpdateTime = time.Now()
		if err := tx.Model(&entity.RechargeRecord{}).
			Where("id = ?", record.ID).
			Select("status", "paymentId", "updateTime").
			Updates(&record).Error; err != nil {
			return errno.New(errno.SystemError)
		}

		var user entity.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&entity.User{}).
			Where("id = ?", record.UserID).
			Take(&user).Error; err != nil {
			return errno.New(errno.SystemError)
		}
		newBalance := user.Balance + record.Amount
		if err := tx.Model(&entity.User{}).
			Where("id = ?", record.UserID).
			Update("balance", newBalance).Error; err != nil {
			return errno.New(errno.SystemError)
		}
		billingRecord := &entity.BillingRecord{
			UserID:        record.UserID,
			Amount:        record.Amount,
			BalanceBefore: user.Balance,
			BalanceAfter:  newBalance,
			Description:   "充值：" + paymentID,
			BillingType:   "recharge",
			CreateTime:    time.Now(),
		}
		if err := tx.Create(billingRecord).Error; err != nil {
			return errno.New(errno.SystemError)
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
