package service

import (
	"time"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/repository"
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
	record, err := s.repo.GetByID(recordID)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if record == nil {
		return false, errno.New(errno.NotFoundError)
	}
	// 防止重复回调重复入账
	if record.Status == "success" {
		return true, nil
	}
	record.Status = "success"
	record.PaymentID = paymentID
	record.UpdateTime = time.Now()
	ok, err := s.repo.UpdateByID(record)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if !ok {
		return false, errno.New(errno.OperationError)
	}
	if err = s.balanceService.AddBalance(record.UserID, record.Amount, "充值："+paymentID); err != nil {
		return false, err
	}
	return true, nil
}
