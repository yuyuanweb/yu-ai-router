package service

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/model/vo"
	"github.com/yupi/airouter/go-backend/internal/repository"
)

const (
	apiKeyPrefix             = "sk-"
	apiKeyRandomHexBytes     = 16
	apiKeyMaskMinLength      = 12
	apiKeyMaskPrefixLength   = 8
	apiKeyMaskSuffixLength   = 4
)

type ApiKeyService struct {
	apiKeyRepo *repository.ApiKeyRepository
}

func NewApiKeyService(apiKeyRepo *repository.ApiKeyRepository) *ApiKeyService {
	return &ApiKeyService{apiKeyRepo: apiKeyRepo}
}

func (s *ApiKeyService) CreateApiKey(keyName string, loginUser *entity.User) (*entity.ApiKey, error) {
	keyValue, err := generateAPIKey()
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	apiKey := &entity.ApiKey{
		UserID:      loginUser.ID,
		KeyValue:    keyValue,
		KeyName:     keyName,
		Status:      "active",
		TotalTokens: 0,
	}
	if err = s.apiKeyRepo.Create(apiKey); err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return apiKey, nil
}

func (s *ApiKeyService) ListMyApiKeys(userID, pageNum, pageSize int64) (common.PageResponse[vo.ApiKeyVO], error) {
	list, total, err := s.apiKeyRepo.ListByUser(userID, pageNum, pageSize)
	if err != nil {
		return common.PageResponse[vo.ApiKeyVO]{}, errno.New(errno.SystemError)
	}
	voList := make([]vo.ApiKeyVO, 0, len(list))
	for _, item := range list {
		voList = append(voList, vo.ApiKeyVO{
			ID:           item.ID,
			KeyValue:     maskAPIKey(item.KeyValue),
			KeyName:      item.KeyName,
			Status:       item.Status,
			TotalTokens:  item.TotalTokens,
			LastUsedTime: item.LastUsedTime,
			CreateTime:   item.CreateTime,
		})
	}
	return common.BuildPageResponse(voList, pageNum, pageSize, total), nil
}

func (s *ApiKeyService) RevokeApiKey(id, userID int64) (bool, error) {
	ok, err := s.apiKeyRepo.RevokeByIDAndUser(id, userID)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if !ok {
		return false, errno.NewWithMessage(errno.NotFoundError, "API Key 不存在")
	}
	return true, nil
}

func (s *ApiKeyService) GetByID(id int64) (*entity.ApiKey, error) {
	apiKey, err := s.apiKeyRepo.GetByID(id)
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return apiKey, nil
}

func (s *ApiKeyService) GetByKeyValue(keyValue string) (*entity.ApiKey, error) {
	apiKey, err := s.apiKeyRepo.GetByKeyValueActive(keyValue)
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return apiKey, nil
}

func (s *ApiKeyService) UpdateUsageStats(apiKeyID int64, tokens int) {
	_ = s.apiKeyRepo.AddUsageStats(apiKeyID, tokens)
}

func generateAPIKey() (string, error) {
	randomBytes := make([]byte, apiKeyRandomHexBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return apiKeyPrefix + hex.EncodeToString(randomBytes), nil
}

func maskAPIKey(key string) string {
	if len(key) <= apiKeyMaskMinLength {
		return key
	}
	prefix := key[:apiKeyMaskPrefixLength]
	suffix := key[len(key)-apiKeyMaskSuffixLength:]
	return strings.Join([]string{prefix, "****", suffix}, "")
}
