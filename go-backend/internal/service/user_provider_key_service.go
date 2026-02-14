package service

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"strings"

	"github.com/yupi/airouter/go-backend/internal/config"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/model/vo"
	"github.com/yupi/airouter/go-backend/internal/repository"
)

const (
	userProviderKeyStatusActive   = "active"
	userProviderKeyStatusInactive = "inactive"
)

type UserProviderKeyService struct {
	repo            *repository.UserProviderKeyRepository
	providerService *ProviderService
	secretKey       string
}

func NewUserProviderKeyService(
	repo *repository.UserProviderKeyRepository,
	providerService *ProviderService,
	cfg *config.Config,
) *UserProviderKeyService {
	secret := strings.TrimSpace(cfg.EncryptionSecretKey)
	if secret == "" {
		secret = strings.TrimSpace(cfg.SessionSecret)
	}
	return &UserProviderKeyService{
		repo:            repo,
		providerService: providerService,
		secretKey:       secret,
	}
}

func (s *UserProviderKeyService) AddUserProviderKey(request dto.UserProviderKeyAddRequest, userID int64) (bool, error) {
	if request.ProviderID == nil || request.ProviderID.Int64() <= 0 || strings.TrimSpace(request.APIKey) == "" {
		return false, errno.New(errno.ParamsError)
	}
	provider, err := s.providerService.GetProviderByID(request.ProviderID.Int64())
	if err != nil {
		return false, err
	}
	if provider == nil {
		return false, errno.NewWithMessage(errno.ParamsError, "提供者不存在")
	}
	existing, err := s.repo.GetByUserAndProvider(userID, request.ProviderID.Int64())
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if existing != nil {
		return false, errno.NewWithMessage(errno.ParamsError, "已配置过该提供者的密钥，请使用更新功能")
	}

	encryptedKey, err := s.encrypt(strings.TrimSpace(request.APIKey))
	if err != nil {
		return false, err
	}
	record := &entity.UserProviderKey{
		UserID:       userID,
		ProviderID:   request.ProviderID.Int64(),
		ProviderName: provider.ProviderName,
		APIKey:       encryptedKey,
		Status:       userProviderKeyStatusActive,
	}
	if err = s.repo.Create(record); err != nil {
		return false, errno.New(errno.SystemError)
	}
	return true, nil
}

func (s *UserProviderKeyService) UpdateUserProviderKey(request dto.UserProviderKeyUpdateRequest, userID int64) (bool, error) {
	if request.ID == nil || request.ID.Int64() <= 0 {
		return false, errno.New(errno.ParamsError)
	}
	record, err := s.repo.GetByID(request.ID.Int64())
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if record == nil {
		return false, errno.New(errno.NotFoundError)
	}
	if record.UserID != userID {
		return false, errno.New(errno.NoAuthError)
	}

	updates := make(map[string]any)
	if strings.TrimSpace(request.APIKey) != "" {
		encryptedKey, encryptErr := s.encrypt(strings.TrimSpace(request.APIKey))
		if encryptErr != nil {
			return false, encryptErr
		}
		updates["apiKey"] = encryptedKey
	}
	status := strings.TrimSpace(strings.ToLower(request.Status))
	if status != "" {
		if status != userProviderKeyStatusActive && status != userProviderKeyStatusInactive {
			return false, errno.NewWithMessage(errno.ParamsError, "状态不合法")
		}
		updates["status"] = status
	}
	if len(updates) == 0 {
		return false, errno.New(errno.ParamsError)
	}
	ok, updateErr := s.repo.UpdateByID(request.ID.Int64(), updates)
	if updateErr != nil {
		return false, errno.New(errno.SystemError)
	}
	if !ok {
		return false, errno.New(errno.OperationError)
	}
	return true, nil
}

func (s *UserProviderKeyService) DeleteUserProviderKey(id, userID int64) (bool, error) {
	record, err := s.repo.GetByID(id)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if record == nil {
		return false, errno.New(errno.NotFoundError)
	}
	if record.UserID != userID {
		return false, errno.New(errno.NoAuthError)
	}
	ok, deleteErr := s.repo.SoftDeleteByID(id)
	if deleteErr != nil {
		return false, errno.New(errno.SystemError)
	}
	if !ok {
		return false, errno.New(errno.OperationError)
	}
	return true, nil
}

func (s *UserProviderKeyService) ListUserProviderKeys(userID int64) ([]vo.UserProviderKeyVO, error) {
	list, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	result := make([]vo.UserProviderKeyVO, 0, len(list))
	for _, item := range list {
		result = append(result, vo.UserProviderKeyVO{
			ID:           item.ID,
			ProviderID:   item.ProviderID,
			ProviderName: item.ProviderName,
			APIKey:       maskUserProviderKey(item.APIKey),
			Status:       item.Status,
			CreateTime:   item.CreateTime,
			UpdateTime:   item.UpdateTime,
		})
	}
	return result, nil
}

func (s *UserProviderKeyService) GetUserProviderAPIKey(userID, providerID int64) (string, error) {
	record, err := s.repo.GetActiveByUserAndProvider(userID, providerID)
	if err != nil {
		return "", errno.New(errno.SystemError)
	}
	if record == nil {
		return "", nil
	}
	decrypted, decryptErr := s.decrypt(record.APIKey)
	if decryptErr != nil {
		return "", nil
	}
	return decrypted, nil
}

func (s *UserProviderKeyService) HasUserProviderKey(userID, providerID int64) (bool, error) {
	key, err := s.GetUserProviderAPIKey(userID, providerID)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(key) != "", nil
}

func (s *UserProviderKeyService) encrypt(plainText string) (string, error) {
	if strings.TrimSpace(plainText) == "" {
		return "", errno.NewWithMessage(errno.ParamsError, "API Key 不能为空")
	}
	key := fixedKey32Bytes(s.secretKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errno.NewWithMessage(errno.SystemError, "加密失败")
	}
	data := pkcs7Pad([]byte(plainText), aes.BlockSize)
	encrypted := make([]byte, len(data))
	for start := 0; start < len(data); start += aes.BlockSize {
		block.Encrypt(encrypted[start:start+aes.BlockSize], data[start:start+aes.BlockSize])
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (s *UserProviderKeyService) decrypt(cipherText string) (string, error) {
	encoded := strings.TrimSpace(cipherText)
	if encoded == "" {
		return "", errno.NewWithMessage(errno.SystemError, "解密失败")
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", errno.NewWithMessage(errno.SystemError, "解密失败")
	}
	if len(raw)%aes.BlockSize != 0 {
		return "", errno.NewWithMessage(errno.SystemError, "解密失败")
	}
	key := fixedKey32Bytes(s.secretKey)
	block, cipherErr := aes.NewCipher(key)
	if cipherErr != nil {
		return "", errno.NewWithMessage(errno.SystemError, "解密失败")
	}
	decrypted := make([]byte, len(raw))
	for start := 0; start < len(raw); start += aes.BlockSize {
		block.Decrypt(decrypted[start:start+aes.BlockSize], raw[start:start+aes.BlockSize])
	}
	unpadded, unpadErr := pkcs7Unpad(decrypted, aes.BlockSize)
	if unpadErr != nil {
		return "", errno.NewWithMessage(errno.SystemError, "解密失败")
	}
	return string(unpadded), nil
}

func fixedKey32Bytes(secret string) []byte {
	key := make([]byte, 32)
	source := []byte(strings.TrimSpace(secret))
	copy(key, source)
	return key
}

func pkcs7Pad(src []byte, blockSize int) []byte {
	padding := blockSize - (len(src) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func pkcs7Unpad(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	if length == 0 || length%blockSize != 0 {
		return nil, errno.NewWithMessage(errno.SystemError, "解密失败")
	}
	padding := int(src[length-1])
	if padding <= 0 || padding > blockSize || padding > length {
		return nil, errno.NewWithMessage(errno.SystemError, "解密失败")
	}
	for i := 0; i < padding; i++ {
		if src[length-1-i] != byte(padding) {
			return nil, errno.NewWithMessage(errno.SystemError, "解密失败")
		}
	}
	return src[:length-padding], nil
}

func maskUserProviderKey(value string) string {
	const (
		minLength = 12
		prefixLen = 8
		suffixLen = 4
	)
	if len(value) <= minLength {
		return "****"
	}
	return value[:prefixLen] + "****" + value[len(value)-suffixLen:]
}
