package service

import (
	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/repository"
)

type BillingRecordService struct {
	repo *repository.BillingRecordRepository
}

func NewBillingRecordService(repo *repository.BillingRecordRepository) *BillingRecordService {
	return &BillingRecordService{repo: repo}
}

func (s *BillingRecordService) ListUserBillingRecords(userID, pageNum, pageSize int64) (common.PageResponse[entity.BillingRecord], error) {
	if userID <= 0 {
		return common.PageResponse[entity.BillingRecord]{}, errno.New(errno.ParamsError)
	}
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.ListByUserID(userID, pageNum, pageSize)
}

func (s *BillingRecordService) GetUserTotalSpending(userID int64) (float64, error) {
	if userID <= 0 {
		return 0, errno.New(errno.ParamsError)
	}
	return s.repo.SumAmountByType(userID, "api_call")
}

func (s *BillingRecordService) GetUserTotalRecharge(userID int64) (float64, error) {
	if userID <= 0 {
		return 0, errno.New(errno.ParamsError)
	}
	return s.repo.SumAmountByType(userID, "recharge")
}

func (s *BillingRecordService) CreateRecord(record *entity.BillingRecord) error {
	return s.repo.Create(record)
}
