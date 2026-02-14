package service

import (
	"log"

	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/repository"
)

const defaultLogLimit = 100

type RequestLogService struct {
	requestLogRepo *repository.RequestLogRepository
	apiKeyService  *ApiKeyService
}

func NewRequestLogService(requestLogRepo *repository.RequestLogRepository, apiKeyService *ApiKeyService) *RequestLogService {
	return &RequestLogService{
		requestLogRepo: requestLogRepo,
		apiKeyService:  apiKeyService,
	}
}

func (s *RequestLogService) LogRequestAsync(
	userID *int64,
	apiKeyID *int64,
	modelName string,
	promptTokens int,
	completionTokens int,
	totalTokens int,
	duration int,
	status string,
	errorMessage string,
) {
	go func() {
		if err := s.logRequest(userID, apiKeyID, modelName, promptTokens, completionTokens, totalTokens, duration, status, errorMessage); err != nil {
			log.Printf("log request failed: %v", err)
		}
	}()
}

func (s *RequestLogService) ListUserLogs(userID int64, limit int) ([]entity.RequestLog, error) {
	if limit <= 0 {
		limit = defaultLogLimit
	}
	list, err := s.requestLogRepo.ListUserLogs(userID, limit)
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return list, nil
}

func (s *RequestLogService) CountUserTokens(userID int64) (int64, error) {
	total, err := s.requestLogRepo.CountUserTokens(userID)
	if err != nil {
		return 0, errno.New(errno.SystemError)
	}
	return total, nil
}

func (s *RequestLogService) logRequest(
	userID *int64,
	apiKeyID *int64,
	modelName string,
	promptTokens int,
	completionTokens int,
	totalTokens int,
	duration int,
	status string,
	errorMessage string,
) error {
	record := &entity.RequestLog{
		UserID:           userID,
		APIKeyID:         apiKeyID,
		ModelName:        modelName,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      totalTokens,
		Duration:         duration,
		Status:           status,
		ErrorMessage:     errorMessage,
	}
	if err := s.requestLogRepo.Create(record); err != nil {
		return err
	}

	if status == "success" && apiKeyID != nil && totalTokens > 0 {
		s.apiKeyService.UpdateUsageStats(*apiKeyID, totalTokens)
	}
	return nil
}
