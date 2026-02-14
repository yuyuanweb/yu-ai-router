package service

import (
	"log"
	"time"

	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/repository"
)

const defaultLogLimit = 100

type RequestLogInput struct {
	TraceID          string
	UserID           *int64
	APIKeyID         *int64
	ModelID          *int64
	RequestModel     string
	ModelName        string
	RequestType      string
	Source           string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	Duration         int
	Status           string
	ErrorMessage     string
	ErrorCode        string
	RoutingStrategy  string
	IsFallback       bool
	ClientIP         string
	UserAgent        string
}

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

func (s *RequestLogService) LogRequestAsync(input RequestLogInput) {
	go func() {
		if err := s.logRequest(input); err != nil {
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

func (s *RequestLogService) logRequest(input RequestLogInput) error {
	isFallback := 0
	if input.IsFallback {
		isFallback = 1
	}
	requestType := input.RequestType
	if requestType == "" {
		requestType = "chat"
	}
	source := input.Source
	if source == "" {
		source = "web"
	}
	modelName := input.ModelName
	if modelName == "" {
		modelName = input.RequestModel
	}

	record := &entity.RequestLog{
		TraceID:          input.TraceID,
		UserID:           input.UserID,
		APIKeyID:         input.APIKeyID,
		ModelID:          input.ModelID,
		RequestModel:     input.RequestModel,
		ModelName:        modelName,
		RequestType:      requestType,
		Source:           source,
		PromptTokens:     input.PromptTokens,
		CompletionTokens: input.CompletionTokens,
		TotalTokens:      input.TotalTokens,
		Duration:         input.Duration,
		Status:           input.Status,
		ErrorMessage:     input.ErrorMessage,
		ErrorCode:        input.ErrorCode,
		RoutingStrategy:  input.RoutingStrategy,
		IsFallback:       isFallback,
		ClientIP:         input.ClientIP,
		UserAgent:        input.UserAgent,
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	}
	if err := s.requestLogRepo.Create(record); err != nil {
		return err
	}

	if input.Status == "success" && input.APIKeyID != nil && input.TotalTokens > 0 {
		s.apiKeyService.UpdateUsageStats(*input.APIKeyID, input.TotalTokens)
	}
	return nil
}
