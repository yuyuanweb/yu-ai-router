package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"

	"github.com/yupi/airouter/go-backend/internal/config"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/model/vo"
	"github.com/yupi/airouter/go-backend/internal/repository"
)

const (
	pluginStatusActive = "active"
	pluginKeyWebSearch = "web_search"
	pluginKeyPDFParser = "pdf_parser"
	pluginKeyImageRec  = "image_recognition"

	defaultWebSearchTimeoutMs = 15000
	defaultWebSearchMaxResult = 5
	defaultPDFMaxPages        = 50
	defaultPDFMaxTextLength   = 50000
	defaultImageMaxBytes      = 4 * 1024 * 1024
	defaultImageModelKey      = "qwen-vl-plus"
)

type PluginService struct {
	pluginRepo      *repository.PluginRepository
	providerService *ProviderService
	cfg             *config.Config
	httpClient      *http.Client
}

func NewPluginService(pluginRepo *repository.PluginRepository, providerService *ProviderService, cfg *config.Config) *PluginService {
	return &PluginService{
		pluginRepo:      pluginRepo,
		providerService: providerService,
		cfg:             cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (s *PluginService) ListAllPlugins() ([]vo.PluginConfigVO, error) {
	list, err := s.pluginRepo.ListAll()
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return s.toPluginConfigVOList(list), nil
}

func (s *PluginService) ListEnabledPlugins() ([]vo.PluginConfigVO, error) {
	list, err := s.pluginRepo.ListEnabled()
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	return s.toPluginConfigVOList(list), nil
}

func (s *PluginService) GetPluginByKey(pluginKey string) (*vo.PluginConfigVO, error) {
	plugin, err := s.pluginRepo.GetByKey(strings.TrimSpace(pluginKey))
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	if plugin == nil {
		return nil, errno.New(errno.NotFoundError)
	}
	result := s.toPluginConfigVO(plugin)
	return &result, nil
}

func (s *PluginService) UpdatePlugin(request dto.PluginUpdateRequest) (bool, error) {
	if request.ID == nil || request.ID.Int64() <= 0 {
		return false, errno.New(errno.ParamsError)
	}
	fields := make(map[string]any)
	if request.PluginName != nil {
		fields["pluginName"] = *request.PluginName
	}
	if request.Description != nil {
		fields["description"] = *request.Description
	}
	if request.Config != nil {
		fields["config"] = *request.Config
	}
	if request.Status != nil {
		fields["status"] = *request.Status
	}
	if request.Priority != nil {
		fields["priority"] = *request.Priority
	}
	ok, err := s.pluginRepo.UpdateByID(request.ID.Int64(), fields)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	return ok, nil
}

func (s *PluginService) EnablePlugin(pluginKey string) (bool, error) {
	plugin, err := s.pluginRepo.GetByKey(strings.TrimSpace(pluginKey))
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if plugin == nil {
		return false, errno.New(errno.NotFoundError)
	}
	return s.pluginRepo.UpdateByID(plugin.ID, map[string]any{"status": pluginStatusActive})
}

func (s *PluginService) DisablePlugin(pluginKey string) (bool, error) {
	plugin, err := s.pluginRepo.GetByKey(strings.TrimSpace(pluginKey))
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if plugin == nil {
		return false, errno.New(errno.NotFoundError)
	}
	return s.pluginRepo.UpdateByID(plugin.ID, map[string]any{"status": "inactive"})
}

func (s *PluginService) ReloadPlugin(pluginKey string) error {
	_, err := s.GetPluginByKey(pluginKey)
	return err
}

func (s *PluginService) InitPlugins() error {
	_, err := s.pluginRepo.ListAll()
	if err != nil {
		return errno.New(errno.SystemError)
	}
	return nil
}

func (s *PluginService) ExecutePlugin(request dto.PluginExecuteRequest, userID *int64) (vo.PluginExecuteVO, error) {
	start := time.Now()
	pluginKey := strings.TrimSpace(request.PluginKey)
	if pluginKey == "" {
		return vo.PluginExecuteVO{}, errno.NewWithMessage(errno.ParamsError, "插件标识不能为空")
	}
	pluginConfig, err := s.pluginRepo.GetByKey(pluginKey)
	if err != nil {
		return vo.PluginExecuteVO{}, errno.New(errno.SystemError)
	}
	if pluginConfig == nil {
		return vo.PluginExecuteVO{}, errno.NewWithMessage(errno.NotFoundError, "插件不存在")
	}
	if pluginConfig.Status != pluginStatusActive {
		return vo.PluginExecuteVO{}, errno.NewWithMessage(errno.OperationError, "插件未启用")
	}

	var result vo.PluginExecuteVO
	switch pluginKey {
	case pluginKeyWebSearch:
		result = s.executeWebSearch(request)
	case pluginKeyPDFParser:
		result = s.executePDFParser(request)
	case pluginKeyImageRec:
		result = s.executeImageRecognition(request)
	default:
		result = vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKey,
			ErrorMessage: "不支持的插件类型",
		}
	}

	result.Duration = time.Since(start).Milliseconds()
	if !result.Success {
		log.Printf("plugin execute failed: plugin=%s userId=%v err=%s", pluginKey, userID, result.ErrorMessage)
	}
	return result, nil
}

func (s *PluginService) executeWebSearch(request dto.PluginExecuteRequest) vo.PluginExecuteVO {
	apiKey := strings.TrimSpace(s.cfg.PluginSerpAPIKey)
	if apiKey == "" {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyWebSearch,
			ErrorMessage: "PLUGIN_SERPAPI_API_KEY 未配置",
		}
	}
	query := strings.TrimSpace(request.Input)
	if query == "" {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyWebSearch,
			ErrorMessage: "搜索关键词不能为空",
		}
	}

	maxResult := defaultWebSearchMaxResult
	timeoutMs := defaultWebSearchTimeoutMs
	engine := "google"
	if request.Params != nil {
		if v, ok := request.Params["maxResults"]; ok {
			if parsed := parsePositiveInt(v); parsed > 0 {
				maxResult = parsed
			}
		}
		if v, ok := request.Params["timeout"]; ok {
			if parsed := parsePositiveInt(v); parsed > 0 {
				timeoutMs = parsed
			}
		}
		if v, ok := request.Params["searchEngine"]; ok {
			if text := strings.TrimSpace(fmt.Sprintf("%v", v)); text != "" {
				engine = text
			}
		}
	}

	client := &http.Client{Timeout: time.Duration(timeoutMs) * time.Millisecond}
	apiURL := "https://serpapi.com/search.json?api_key=" + url.QueryEscape(apiKey) +
		"&q=" + url.QueryEscape(query) +
		"&engine=" + url.QueryEscape(engine) +
		"&num=" + strconv.Itoa(maxResult) +
		"&hl=zh-CN&gl=cn"
	resp, err := client.Get(apiURL)
	if err != nil {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyWebSearch,
			ErrorMessage: "搜索请求失败: " + err.Error(),
		}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyWebSearch,
			ErrorMessage: "搜索请求失败，状态码: " + strconv.Itoa(resp.StatusCode),
		}
	}

	type organicResult struct {
		Title   string `json:"title"`
		Link    string `json:"link"`
		Snippet string `json:"snippet"`
	}
	var parsed struct {
		Error          string          `json:"error"`
		OrganicResults []organicResult `json:"organic_results"`
	}
	if err = json.Unmarshal(body, &parsed); err != nil {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyWebSearch,
			ErrorMessage: "搜索结果解析失败: " + err.Error(),
		}
	}
	if strings.TrimSpace(parsed.Error) != "" {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyWebSearch,
			ErrorMessage: parsed.Error,
		}
	}
	if len(parsed.OrganicResults) == 0 {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyWebSearch,
			ErrorMessage: "未找到相关搜索结果",
		}
	}

	builder := strings.Builder{}
	builder.WriteString("## 搜索结果\n\n")
	for i, item := range parsed.OrganicResults {
		builder.WriteString(strconv.Itoa(i + 1))
		builder.WriteString(". ")
		builder.WriteString(strings.TrimSpace(item.Title))
		builder.WriteString("\n")
		builder.WriteString(strings.TrimSpace(item.Snippet))
		builder.WriteString("\n来源: ")
		builder.WriteString(strings.TrimSpace(item.Link))
		builder.WriteString("\n\n")
	}

	return vo.PluginExecuteVO{
		Success:   true,
		PluginKey: pluginKeyWebSearch,
		Content:   builder.String(),
		Data: map[string]any{
			"query":       query,
			"resultCount": len(parsed.OrganicResults),
			"searchEngine": engine,
		},
	}
}

func (s *PluginService) executePDFParser(request dto.PluginExecuteRequest) vo.PluginExecuteVO {
	pdfBytes, err := s.loadFileBytes(request.FileURL, request.FileBytes)
	if err != nil {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyPDFParser,
			ErrorMessage: err.Error(),
		}
	}
	text, pages, parseErr := parsePDFText(pdfBytes, defaultPDFMaxPages, defaultPDFMaxTextLength)
	if parseErr != nil {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyPDFParser,
			ErrorMessage: "PDF文档解析失败: " + parseErr.Error(),
		}
	}
	preview, totalLines := firstNLines(text, 100)
	log.Printf("pdf_parser parsed content preview: pages=%d textLength=%d totalLines=%d previewLines=%d content=\n%s\n",
		pages, len(text), totalLines, 100, preview)
	return vo.PluginExecuteVO{
		Success:   true,
		PluginKey: pluginKeyPDFParser,
		Content:   "## PDF文档内容\n\n```\n" + text + "\n```",
		Data: map[string]any{
			"textLength": len(text),
			"pages":      pages,
			"maxPages":   defaultPDFMaxPages,
		},
	}
}

func (s *PluginService) executeImageRecognition(request dto.PluginExecuteRequest) vo.PluginExecuteVO {
	imageBytes, err := s.loadFileBytes(request.FileURL, request.FileBytes)
	if err != nil {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyImageRec,
			ErrorMessage: err.Error(),
		}
	}
	if len(imageBytes) > defaultImageMaxBytes {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyImageRec,
			ErrorMessage: "图片大小超过限制（最大 4MB）",
		}
	}

	provider, providerErr := s.providerService.GetProviderByName("qwen")
	if providerErr != nil {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyImageRec,
			ErrorMessage: "未找到通义千问提供者配置",
		}
	}

	prompt := strings.TrimSpace(request.Input)
	if prompt == "" {
		prompt = "请描述这张图片的内容，包括主要对象、场景、颜色、文字等信息。"
	}
	mimeType := strings.TrimSpace(request.FileType)
	if mimeType == "" {
		mimeType = "image/jpeg"
	}
	imageDataURL := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(imageBytes)
	endpoint := strings.TrimRight(provider.BaseURL, "/") + "/v1/chat/completions"
	payload := map[string]any{
		"model": defaultImageModelKey,
		"messages": []any{
			map[string]any{
				"role": "user",
				"content": []any{
					map[string]any{
						"type": "image_url",
						"image_url": map[string]any{
							"url": imageDataURL,
						},
					},
					map[string]any{
						"type": "text",
						"text": prompt,
					},
				},
			},
		},
	}
	payloadBytes, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payloadBytes))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+provider.APIKey)
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyImageRec,
			ErrorMessage: "调用视觉模型失败: " + err.Error(),
		}
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("image_recognition upstream failed: status=%d body=%s", resp.StatusCode, string(respBody))
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyImageRec,
			ErrorMessage: "视觉模型调用失败，状态码: " + strconv.Itoa(resp.StatusCode),
		}
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err = json.Unmarshal(respBody, &response); err != nil {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyImageRec,
			ErrorMessage: "视觉模型结果解析失败: " + err.Error(),
		}
	}
	if len(response.Choices) == 0 || strings.TrimSpace(response.Choices[0].Message.Content) == "" {
		return vo.PluginExecuteVO{
			Success:      false,
			PluginKey:    pluginKeyImageRec,
			ErrorMessage: "视觉模型返回为空",
		}
	}
	return vo.PluginExecuteVO{
		Success:   true,
		PluginKey: pluginKeyImageRec,
		Content:   "## 图片识别结果\n\n" + response.Choices[0].Message.Content,
	}
}

func (s *PluginService) loadFileBytes(fileURL string, fileBytes []byte) ([]byte, error) {
	if len(fileBytes) > 0 {
		return fileBytes, nil
	}
	fileURL = strings.TrimSpace(fileURL)
	if fileURL == "" {
		return nil, fmt.Errorf("请提供文件内容或文件 URL")
	}
	resp, err := s.httpClient.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("下载文件失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载文件失败，状态码: %d", resp.StatusCode)
	}
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("读取文件失败: %w", readErr)
	}
	return body, nil
}

func (s *PluginService) toPluginConfigVO(config *entity.PluginConfig) vo.PluginConfigVO {
	return vo.PluginConfigVO{
		ID:          config.ID,
		PluginKey:   config.PluginKey,
		PluginName:  config.PluginName,
		PluginType:  config.PluginType,
		Description: config.Description,
		Config:      config.Config,
		Status:      config.Status,
		Priority:    config.Priority,
		CreateTime:  config.CreateTime,
		UpdateTime:  config.UpdateTime,
	}
}

func (s *PluginService) toPluginConfigVOList(list []entity.PluginConfig) []vo.PluginConfigVO {
	if len(list) == 0 {
		return make([]vo.PluginConfigVO, 0)
	}
	result := make([]vo.PluginConfigVO, 0, len(list))
	for _, item := range list {
		itemCopy := item
		result = append(result, s.toPluginConfigVO(&itemCopy))
	}
	return result
}

func parsePositiveInt(v any) int {
	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case int:
		if n > 0 {
			return n
		}
	case int64:
		if n > 0 {
			return int(n)
		}
	case float64:
		if n > 0 {
			return int(n)
		}
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(n))
		if err == nil && parsed > 0 {
			return parsed
		}
	}
	return 0
}

func firstNLines(text string, n int) (string, int) {
	if n <= 0 {
		return "", 0
	}
	lines := strings.Split(text, "\n")
	total := len(lines)
	if total <= n {
		return text, total
	}
	return strings.Join(lines[:n], "\n"), total
}

func parsePDFText(data []byte, maxPages, maxTextLength int) (string, int, error) {
	if len(data) == 0 {
		return "", 0, fmt.Errorf("文件内容为空")
	}

	reader := bytes.NewReader(data)
	pdfReader, err := pdf.NewReader(reader, int64(len(data)))
	if err != nil {
		return "", 0, fmt.Errorf("PDF格式无效或已损坏")
	}

	totalPages := pdfReader.NumPage()
	if totalPages <= 0 {
		return "", 0, fmt.Errorf("PDF页数异常")
	}
	usedPages := totalPages
	if usedPages > maxPages {
		usedPages = maxPages
	}

	builder := strings.Builder{}
	for pageNum := 1; pageNum <= usedPages; pageNum++ {
		page := pdfReader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}
		pageText, contentErr := page.GetPlainText(nil)
		if contentErr != nil {
			continue
		}
		if strings.TrimSpace(pageText) == "" {
			continue
		}
		builder.WriteString(pageText)
		builder.WriteString("\n")
		if builder.Len() >= maxTextLength {
			break
		}
	}

	text := strings.TrimSpace(builder.String())
	if text == "" {
		return "", usedPages, fmt.Errorf("未提取到文本（可能是扫描件或纯图片PDF）")
	}
	if len(text) > maxTextLength {
		text = text[:maxTextLength] + "\n\n[文档内容过长，已截断...]"
	}
	return text, usedPages, nil
}
