package service

import "time"

type BillingService struct {
	requestLogService *RequestLogService
}

func NewBillingService(requestLogService *RequestLogService) *BillingService {
	return &BillingService{requestLogService: requestLogService}
}

func (s *BillingService) GetUserTotalCost(userID int64) (float64, error) {
	return s.requestLogService.SumUserCost(userID, nil, nil)
}

func (s *BillingService) GetUserTodayCost(userID int64) (float64, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24*time.Hour - time.Nanosecond)
	return s.requestLogService.SumUserCost(userID, &start, &end)
}
