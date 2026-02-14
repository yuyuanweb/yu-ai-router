package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/repository"
)

const (
	defaultImageModel   = "qwen-image-plus"
	defaultImageSize    = "1024*1024"
	defaultImageCount   = 1
	imagePollInterval   = 2 * time.Second
	imagePollMaxRetries = 60
	imageTokensPerItem  = 1000
)

type ImageGenerationService struct {
	imageRecordRepo *repository.ImageGenerationRecordRepository
	modelService    *ModelService
	providerService *ProviderService
	userService     *UserService
	balanceService  *BalanceService
}

func NewImageGenerationService(
	imageRecordRepo *repository.ImageGenerationRecordRepository,
	modelService *ModelService,
	providerService *ProviderService,
	userService *UserService,
	balanceService *BalanceService,
) *ImageGenerationService {
	return &ImageGenerationService{
		imageRecordRepo: imageRecordRepo,
		modelService:    modelService,
		providerService: providerService,
		userService:     userService,
		balanceService:  balanceService,
	}
}

func (s *ImageGenerationService) GenerateImage(
	request dto.ImageGenerationRequest,
	userID *int64,
	apiKeyID *int64,
	clientIP string,
) (*dto.ImageGenerationResponse, error) {
	start := time.Now()
	prompt := strings.TrimSpace(request.Prompt)
	if prompt == "" {
		return nil, errno.NewWithMessage(errno.ParamsError, "提示词不能为空")
	}

	modelKey := strings.TrimSpace(request.Model)
	if modelKey == "" {
		modelKey = defaultImageModel
	}
	size := strings.TrimSpace(request.Size)
	if size == "" {
		size = defaultImageSize
	}
	n := defaultImageCount
	if request.N != nil && *request.N > 0 {
		n = *request.N
	}

	if userID != nil && *userID > 0 {
		disabled, disabledErr := s.userService.IsUserDisabled(*userID)
		if disabledErr != nil {
			return nil, disabledErr
		}
		if disabled {
			return nil, errno.NewWithMessage(errno.ForbiddenError, "账号已被禁用，无法使用服务")
		}
		quotaEnough, quotaErr := s.userService.CheckQuota(*userID)
		if quotaErr != nil {
			return nil, quotaErr
		}
		if !quotaEnough {
			return nil, errno.NewWithMessage(errno.OperationError, "Token配额已用尽，请联系管理员增加配额")
		}
	}

	model, err := s.modelService.GetByModelKey(modelKey)
	if err != nil {
		return nil, err
	}
	if model.ModelType != "image" {
		return nil, errno.NewWithMessage(errno.ParamsError, "模型不存在或不是绘图模型")
	}
	provider, err := s.providerService.GetProviderByID(model.ProviderID)
	if err != nil {
		return nil, err
	}

	estimatedCost := model.InputPrice * float64(n)
	if userID != nil && *userID > 0 {
		ok, checkErr := s.balanceService.CheckBalance(*userID, estimatedCost)
		if checkErr != nil {
			return nil, checkErr
		}
		if !ok {
			return nil, errno.NewWithMessage(errno.OperationError, fmt.Sprintf("账户余额不足，生成%d张图片预计需要¥%.4f，请先充值", n, estimatedCost))
		}
	}

	response, genErr := s.callImageModel(provider, modelKey, prompt, size, n, request.ResponseFormat)
	duration := int(time.Since(start).Milliseconds())
	if genErr != nil {
		if userID != nil && *userID > 0 {
			_ = s.imageRecordRepo.Create(&entity.ImageGenerationRecord{
				UserID:       *userID,
				APIKeyID:     apiKeyID,
				ModelID:      model.ID,
				ModelKey:     modelKey,
				Prompt:       prompt,
				Size:         size,
				Quality:      request.Quality,
				Status:       "failed",
				Cost:         0,
				Duration:     duration,
				ErrorMessage: genErr.Error(),
				ClientIP:     clientIP,
				CreateTime:   time.Now(),
			})
		}
		return nil, genErr
	}

	actualCount := len(response.Data)
	if actualCount == 0 {
		actualCount = n
	}
	actualCost := model.InputPrice * float64(actualCount)
	for _, item := range response.Data {
		record := &entity.ImageGenerationRecord{
			UserID:        zeroIfNil(userID),
			APIKeyID:      apiKeyID,
			ModelID:       model.ID,
			ModelKey:      modelKey,
			Prompt:        prompt,
			RevisedPrompt: item.RevisedPrompt,
			ImageURL:      item.URL,
			ImageData:     item.B64JSON,
			Size:          size,
			Quality:       request.Quality,
			Status:        "success",
			Cost:          model.InputPrice,
			Duration:      duration,
			ClientIP:      clientIP,
			CreateTime:    time.Now(),
		}
		_ = s.imageRecordRepo.Create(record)
	}

	if userID != nil && *userID > 0 {
		_ = s.userService.DeductTokens(*userID, int64(actualCount*imageTokensPerItem))
		if actualCost > 0 {
			desc := fmt.Sprintf("网页图片生成 - %s x%d", modelKey, actualCount)
			if apiKeyID != nil && *apiKeyID > 0 {
				desc = fmt.Sprintf("API图片生成 - %s x%d", modelKey, actualCount)
			}
			if deductErr := s.balanceService.DeductBalance(*userID, actualCost, nil, desc); deductErr != nil {
				log.Printf("image deduct balance failed: userId=%d model=%s cost=%.4f err=%v", *userID, modelKey, actualCost, deductErr)
				return nil, deductErr
			}
		}
	}

	return response, nil
}

func (s *ImageGenerationService) ListUserRecords(userID, pageNum, pageSize int64) (common.PageResponse[entity.ImageGenerationRecord], error) {
	if userID <= 0 {
		return common.PageResponse[entity.ImageGenerationRecord]{}, errno.New(errno.ParamsError)
	}
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.imageRecordRepo.ListByUserID(userID, pageNum, pageSize)
}

func (s *ImageGenerationService) callImageModel(
	provider *entity.ModelProvider,
	modelKey, prompt, size string,
	n int,
	responseFormat string,
) (*dto.ImageGenerationResponse, error) {
	if strings.Contains(provider.BaseURL, "dashscope.aliyuncs.com") && strings.HasPrefix(modelKey, "qwen-image") {
		return s.callQwenAsync(provider, modelKey, prompt, size, n)
	}
	return s.callOpenAIImage(provider, modelKey, prompt, size, n, responseFormat)
}

func (s *ImageGenerationService) callOpenAIImage(
	provider *entity.ModelProvider,
	modelKey, prompt, size string,
	n int,
	responseFormat string,
) (*dto.ImageGenerationResponse, error) {
	baseURL := strings.TrimRight(provider.BaseURL, "/")
	url := baseURL + "/v1/images/generations"
	body := map[string]any{
		"model":  modelKey,
		"prompt": prompt,
		"n":      n,
		"size":   size,
	}
	if strings.TrimSpace(responseFormat) != "" {
		body["response_format"] = responseFormat
	}
	payload, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+provider.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errno.NewWithMessage(errno.SystemError, "调用图片生成API失败: "+err.Error())
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errno.NewWithMessage(errno.SystemError, "调用图片生成API失败: "+string(respBody))
	}
	var result struct {
		Created int64 `json:"created"`
		Data    []struct {
			URL           string `json:"url"`
			B64JSON       string `json:"b64_json"`
			RevisedPrompt string `json:"revised_prompt"`
		} `json:"data"`
	}
	if err = json.Unmarshal(respBody, &result); err != nil {
		return nil, errno.NewWithMessage(errno.SystemError, "解析图片生成响应失败")
	}
	items := make([]dto.ImageDataEntry, 0, len(result.Data))
	for _, item := range result.Data {
		items = append(items, dto.ImageDataEntry{
			URL:           item.URL,
			B64JSON:       item.B64JSON,
			RevisedPrompt: item.RevisedPrompt,
		})
	}
	return &dto.ImageGenerationResponse{
		Created: result.Created,
		Data:    items,
	}, nil
}

func (s *ImageGenerationService) callQwenAsync(
	provider *entity.ModelProvider,
	modelKey, prompt, size string,
	n int,
) (*dto.ImageGenerationResponse, error) {
	base := strings.Replace(strings.TrimRight(provider.BaseURL, "/"), "/compatible-mode", "", 1)
	createURL := base + "/api/v1/services/aigc/text2image/image-synthesis"
	createBody := map[string]any{
		"model": modelKey,
		"input": map[string]any{
			"prompt": prompt,
		},
		"parameters": map[string]any{
			"size": size,
			"n":    n,
		},
	}
	payload, _ := json.Marshal(createBody)
	req, _ := http.NewRequest(http.MethodPost, createURL, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+provider.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-Async", "enable")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errno.NewWithMessage(errno.SystemError, "创建图片生成任务失败: "+err.Error())
	}
	defer resp.Body.Close()
	createRespBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errno.NewWithMessage(errno.SystemError, "创建图片生成任务失败: "+string(createRespBody))
	}
	var createResp struct {
		Output struct {
			TaskID string `json:"task_id"`
		} `json:"output"`
		Message string `json:"message"`
	}
	if err = json.Unmarshal(createRespBody, &createResp); err != nil || createResp.Output.TaskID == "" {
		return nil, errno.NewWithMessage(errno.SystemError, "创建图片生成任务失败")
	}

	taskURL := base + "/api/v1/tasks/" + createResp.Output.TaskID
	for i := 0; i < imagePollMaxRetries; i++ {
		time.Sleep(imagePollInterval)
		taskReq, _ := http.NewRequest(http.MethodGet, taskURL, nil)
		taskReq.Header.Set("Authorization", "Bearer "+provider.APIKey)
		taskResp, taskErr := client.Do(taskReq)
		if taskErr != nil {
			continue
		}
		taskBody, _ := io.ReadAll(taskResp.Body)
		taskResp.Body.Close()
		var taskResult struct {
			Output struct {
				TaskStatus string `json:"task_status"`
				Message    string `json:"message"`
				Results    []struct {
					URL string `json:"url"`
				} `json:"results"`
			} `json:"output"`
		}
		if err = json.Unmarshal(taskBody, &taskResult); err != nil {
			continue
		}
		switch taskResult.Output.TaskStatus {
		case "SUCCEEDED":
			items := make([]dto.ImageDataEntry, 0, len(taskResult.Output.Results))
			for _, item := range taskResult.Output.Results {
				if item.URL != "" {
					items = append(items, dto.ImageDataEntry{
						URL:           item.URL,
						RevisedPrompt: prompt,
					})
				}
			}
			if len(items) == 0 {
				return nil, errno.NewWithMessage(errno.SystemError, "任务成功但未返回图片")
			}
			return &dto.ImageGenerationResponse{
				Created: time.Now().Unix(),
				Data:    items,
			}, nil
		case "FAILED":
			return nil, errno.NewWithMessage(errno.SystemError, "图片生成失败: "+taskResult.Output.Message)
		}
	}
	return nil, errno.NewWithMessage(errno.SystemError, "图片生成超时")
}

func zeroIfNil(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}
