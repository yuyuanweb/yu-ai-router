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
	defaultModelPriority       = 100
	defaultModelTimeoutMS      = 60000
	defaultModelContextLength  = 4096
	defaultModelStatus         = "active"
)

type ModelService struct {
	modelRepo    *repository.ModelRepository
	providerRepo *repository.ProviderRepository
}

func NewModelService(modelRepo *repository.ModelRepository, providerRepo *repository.ProviderRepository) *ModelService {
	return &ModelService{
		modelRepo:    modelRepo,
		providerRepo: providerRepo,
	}
}

func (s *ModelService) AddModel(req dto.ModelAddRequest) (int64, error) {
	if req.ProviderID == nil || req.ProviderID.Int64() <= 0 ||
		strings.TrimSpace(req.ModelKey) == "" ||
		strings.TrimSpace(req.ModelName) == "" {
		return 0, errno.New(errno.ParamsError)
	}
	priority := defaultModelPriority
	if req.Priority != nil {
		priority = *req.Priority
	}
	timeout := defaultModelTimeoutMS
	if req.DefaultTimeout != nil {
		timeout = *req.DefaultTimeout
	}
	contextLength := defaultModelContextLength
	if req.ContextLength != nil {
		contextLength = *req.ContextLength
	}
	inputPrice := float64(0)
	if req.InputPrice != nil {
		inputPrice = *req.InputPrice
	}
	outputPrice := float64(0)
	if req.OutputPrice != nil {
		outputPrice = *req.OutputPrice
	}
	modelType := req.ModelType
	if strings.TrimSpace(modelType) == "" {
		modelType = "chat"
	}

	model := &entity.Model{
		ProviderID:     req.ProviderID.Int64(),
		ModelKey:       req.ModelKey,
		ModelName:      req.ModelName,
		ModelType:      modelType,
		Description:    req.Description,
		ContextLength:  contextLength,
		InputPrice:     inputPrice,
		OutputPrice:    outputPrice,
		Status:         defaultModelStatus,
		Priority:       priority,
		DefaultTimeout: timeout,
		Capabilities:   req.Capabilities,
	}
	if err := s.modelRepo.Create(model); err != nil {
		return 0, errno.New(errno.OperationError)
	}
	return model.ID, nil
}

func (s *ModelService) DeleteModel(id int64) (bool, error) {
	ok, err := s.modelRepo.SoftDeleteByID(id)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	return ok, nil
}

func (s *ModelService) UpdateModel(req dto.ModelUpdateRequest) (bool, error) {
	if req.ID == nil || req.ID.Int64() <= 0 {
		return false, errno.New(errno.ParamsError)
	}
	fields := make(map[string]any)
	if req.ModelName != nil {
		fields["modelName"] = *req.ModelName
	}
	if req.Description != nil {
		fields["description"] = *req.Description
	}
	if req.ContextLength != nil {
		fields["contextLength"] = *req.ContextLength
	}
	if req.InputPrice != nil {
		fields["inputPrice"] = *req.InputPrice
	}
	if req.OutputPrice != nil {
		fields["outputPrice"] = *req.OutputPrice
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	if req.Priority != nil {
		fields["priority"] = *req.Priority
	}
	if req.DefaultTimeout != nil {
		fields["defaultTimeout"] = *req.DefaultTimeout
	}
	if req.Capabilities != nil {
		fields["capabilities"] = *req.Capabilities
	}
	ok, err := s.modelRepo.UpdateByID(req.ID.Int64(), fields)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	return ok, nil
}

func (s *ModelService) GetModelByID(id int64) (*entity.Model, error) {
	model, err := s.modelRepo.GetByID(id)
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	if model == nil {
		return nil, errno.New(errno.NotFoundError)
	}
	return model, nil
}

func (s *ModelService) ListModelVOByPage(req dto.ModelQueryRequest) (common.PageResponse[vo.ModelVO], error) {
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

	list, total, err := s.modelRepo.ListByPage(req)
	if err != nil {
		return common.PageResponse[vo.ModelVO]{}, errno.New(errno.SystemError)
	}
	voList, err := s.GetModelVOList(list)
	if err != nil {
		return common.PageResponse[vo.ModelVO]{}, err
	}
	return common.BuildPageResponse(voList, pageNum, pageSize, total), nil
}

func (s *ModelService) ListModelVO() ([]vo.ModelVO, error) {
	list, err := s.modelRepo.ListAll()
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return s.GetModelVOList(list)
}

func (s *ModelService) ListActiveModelVO() ([]vo.ModelVO, error) {
	list, err := s.modelRepo.ListActive()
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return s.GetModelVOList(list)
}

func (s *ModelService) ListActiveModelVOByProviderID(providerID int64) ([]vo.ModelVO, error) {
	list, err := s.modelRepo.ListActiveByProviderID(providerID)
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return s.GetModelVOList(list)
}

func (s *ModelService) ListActiveModelVOByType(modelType string) ([]vo.ModelVO, error) {
	if strings.TrimSpace(modelType) == "" {
		return make([]vo.ModelVO, 0), nil
	}
	list, err := s.modelRepo.ListActiveByType(modelType)
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return s.GetModelVOList(list)
}

func (s *ModelService) GetModelVO(model *entity.Model) (vo.ModelVO, error) {
	providerList, err := s.providerRepo.ListByIDs([]int64{model.ProviderID})
	if err != nil {
		return vo.ModelVO{}, errno.New(errno.SystemError)
	}
	providerMap := buildProviderMap(providerList)
	return buildModelVO(model, providerMap), nil
}

func (s *ModelService) GetModelVOList(list []entity.Model) ([]vo.ModelVO, error) {
	if len(list) == 0 {
		return make([]vo.ModelVO, 0), nil
	}
	providerIDSet := make(map[int64]struct{})
	for _, item := range list {
		providerIDSet[item.ProviderID] = struct{}{}
	}
	providerIDs := make([]int64, 0, len(providerIDSet))
	for providerID := range providerIDSet {
		providerIDs = append(providerIDs, providerID)
	}
	providerList, err := s.providerRepo.ListByIDs(providerIDs)
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	providerMap := buildProviderMap(providerList)

	result := make([]vo.ModelVO, 0, len(list))
	for _, item := range list {
		itemCopy := item
		result = append(result, buildModelVO(&itemCopy, providerMap))
	}
	return result, nil
}

func buildProviderMap(list []entity.ModelProvider) map[int64]entity.ModelProvider {
	result := make(map[int64]entity.ModelProvider, len(list))
	for _, item := range list {
		result[item.ID] = item
	}
	return result
}

func buildModelVO(model *entity.Model, providerMap map[int64]entity.ModelProvider) vo.ModelVO {
	modelVO := vo.ModelVO{
		ID:               model.ID,
		ProviderID:       model.ProviderID,
		ModelKey:         model.ModelKey,
		ModelName:        model.ModelName,
		ModelType:        model.ModelType,
		Description:      model.Description,
		ContextLength:    model.ContextLength,
		InputPrice:       model.InputPrice,
		OutputPrice:      model.OutputPrice,
		Status:           model.Status,
		HealthStatus:     model.HealthStatus,
		AvgLatency:       model.AvgLatency,
		SuccessRate:      model.SuccessRate,
		Priority:         model.Priority,
		DefaultTimeout:   model.DefaultTimeout,
		SupportReasoning: model.SupportReasoning,
		Capabilities:     model.Capabilities,
		CreateTime:       model.CreateTime,
		UpdateTime:       model.UpdateTime,
	}
	if provider, ok := providerMap[model.ProviderID]; ok {
		modelVO.ProviderName = provider.ProviderName
		modelVO.ProviderDisplayName = provider.DisplayName
	}
	return modelVO
}
