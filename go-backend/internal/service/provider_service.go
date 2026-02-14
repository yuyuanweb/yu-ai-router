package service

import (
	"strings"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/model/vo"
	"github.com/yupi/airouter/go-backend/internal/repository"
)

const (
	defaultProviderPriority = 100
	statusActive            = "active"
)

type ProviderService struct {
	providerRepo *repository.ProviderRepository
}

func NewProviderService(providerRepo *repository.ProviderRepository) *ProviderService {
	return &ProviderService{providerRepo: providerRepo}
}

func (s *ProviderService) AddProvider(req dto.ProviderAddRequest) (int64, error) {
	if strings.TrimSpace(req.ProviderName) == "" ||
		strings.TrimSpace(req.DisplayName) == "" ||
		strings.TrimSpace(req.BaseURL) == "" ||
		strings.TrimSpace(req.APIKey) == "" {
		return 0, errno.New(errno.ParamsError)
	}
	priority := defaultProviderPriority
	if req.Priority != nil {
		priority = *req.Priority
	}
	provider := &entity.ModelProvider{
		ProviderName: req.ProviderName,
		DisplayName:  req.DisplayName,
		BaseURL:      req.BaseURL,
		APIKey:       req.APIKey,
		Status:       statusActive,
		Priority:     priority,
		Config:       req.Config,
	}
	if err := s.providerRepo.Create(provider); err != nil {
		return 0, errno.New(errno.OperationError)
	}
	return provider.ID, nil
}

func (s *ProviderService) DeleteProvider(id int64) (bool, error) {
	ok, err := s.providerRepo.SoftDeleteByID(id)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	return ok, nil
}

func (s *ProviderService) UpdateProvider(req dto.ProviderUpdateRequest) (bool, error) {
	if req.ID == nil || req.ID.Int64() <= 0 {
		return false, errno.New(errno.ParamsError)
	}
	fields := make(map[string]any)
	if req.DisplayName != nil {
		fields["displayName"] = *req.DisplayName
	}
	if req.BaseURL != nil {
		fields["baseUrl"] = *req.BaseURL
	}
	if req.APIKey != nil {
		fields["apiKey"] = *req.APIKey
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	if req.Priority != nil {
		fields["priority"] = *req.Priority
	}
	if req.Config != nil {
		fields["config"] = *req.Config
	}
	ok, err := s.providerRepo.UpdateByID(req.ID.Int64(), fields)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	return ok, nil
}

func (s *ProviderService) GetProviderByID(id int64) (*entity.ModelProvider, error) {
	provider, err := s.providerRepo.GetByID(id)
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	if provider == nil {
		return nil, errno.New(errno.NotFoundError)
	}
	return provider, nil
}

func (s *ProviderService) ListProviderVOByPage(req dto.ProviderQueryRequest) (common.PageResponse[vo.ProviderVO], error) {
	pageNum := req.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	sortOrder := req.SortOrder
	if sortOrder == "" {
		sortOrder = "descend"
	}
	req.PageNum = pageNum
	req.PageSize = pageSize
	req.SortOrder = sortOrder

	list, total, err := s.providerRepo.ListByPage(req)
	if err != nil {
		return common.PageResponse[vo.ProviderVO]{}, errno.New(errno.SystemError)
	}
	return common.BuildPageResponse(s.GetProviderVOList(list), pageNum, pageSize, total), nil
}

func (s *ProviderService) ListProviderVO() ([]vo.ProviderVO, error) {
	list, err := s.providerRepo.ListAll()
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return s.GetProviderVOList(list), nil
}

func (s *ProviderService) ListHealthyProviderVO() ([]vo.ProviderVO, error) {
	list, err := s.providerRepo.ListHealthy()
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return s.GetProviderVOList(list), nil
}

func (s *ProviderService) GetProviderVO(provider *entity.ModelProvider) vo.ProviderVO {
	return vo.ProviderVO{
		ID:           provider.ID,
		ProviderName: provider.ProviderName,
		DisplayName:  provider.DisplayName,
		BaseURL:      provider.BaseURL,
		Status:       provider.Status,
		HealthStatus: provider.HealthStatus,
		AvgLatency:   provider.AvgLatency,
		SuccessRate:  provider.SuccessRate,
		Priority:     provider.Priority,
		Config:       provider.Config,
		CreateTime:   provider.CreateTime,
		UpdateTime:   provider.UpdateTime,
	}
}

func (s *ProviderService) GetProviderVOList(list []entity.ModelProvider) []vo.ProviderVO {
	if len(list) == 0 {
		return make([]vo.ProviderVO, 0)
	}
	result := make([]vo.ProviderVO, 0, len(list))
	for _, item := range list {
		itemCopy := item
		result = append(result, s.GetProviderVO(&itemCopy))
	}
	return result
}
