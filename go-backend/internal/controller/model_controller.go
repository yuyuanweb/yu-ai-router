package controller

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type ModelController struct {
	modelService *service.ModelService
}

func NewModelController(modelService *service.ModelService) *ModelController {
	return &ModelController{modelService: modelService}
}

func (m *ModelController) AddModel(c *gin.Context) {
	var request dto.ModelAddRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		m.handleBindError(c, "add model", err)
		return
	}
	id, err := m.modelService.AddModel(request)
	if err != nil {
		m.handleError(c, err)
		return
	}
	common.Success(c, strconv.FormatInt(id, 10))
}

func (m *ModelController) DeleteModel(c *gin.Context) {
	var request dto.DeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		m.handleBindError(c, "delete model", err)
		return
	}
	if request.ID == nil || request.ID.Int64() <= 0 {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	result, err := m.modelService.DeleteModel(request.ID.Int64())
	if err != nil {
		m.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (m *ModelController) UpdateModel(c *gin.Context) {
	var request dto.ModelUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		m.handleBindError(c, "update model", err)
		return
	}
	result, err := m.modelService.UpdateModel(request)
	if err != nil {
		m.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (m *ModelController) GetModelVOByID(c *gin.Context) {
	id, err := parsePositiveID(c.Query("id"))
	if err != nil {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	model, err := m.modelService.GetModelByID(id)
	if err != nil {
		m.handleError(c, err)
		return
	}
	modelVO, err := m.modelService.GetModelVO(model)
	if err != nil {
		m.handleError(c, err)
		return
	}
	common.Success(c, modelVO)
}

func (m *ModelController) ListModelVOByPage(c *gin.Context) {
	var request dto.ModelQueryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		m.handleBindError(c, "list model by page", err)
		return
	}
	pageResponse, err := m.modelService.ListModelVOByPage(request)
	if err != nil {
		m.handleError(c, err)
		return
	}
	common.Success(c, pageResponse)
}

func (m *ModelController) ListModelVO(c *gin.Context) {
	voList, err := m.modelService.ListModelVO()
	if err != nil {
		m.handleError(c, err)
		return
	}
	common.Success(c, voList)
}

func (m *ModelController) ListActiveModels(c *gin.Context) {
	voList, err := m.modelService.ListActiveModelVO()
	if err != nil {
		m.handleError(c, err)
		return
	}
	common.Success(c, voList)
}

func (m *ModelController) ListActiveModelsByProvider(c *gin.Context) {
	providerID, err := parsePositiveID(c.Param("providerId"))
	if err != nil {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	voList, err := m.modelService.ListActiveModelVOByProviderID(providerID)
	if err != nil {
		m.handleError(c, err)
		return
	}
	common.Success(c, voList)
}

func (m *ModelController) ListActiveModelsByType(c *gin.Context) {
	modelType := c.Param("modelType")
	if modelType == "" {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	voList, err := m.modelService.ListActiveModelVOByType(modelType)
	if err != nil {
		m.handleError(c, err)
		return
	}
	common.Success(c, voList)
}

func (m *ModelController) handleError(c *gin.Context, err error) {
	if bizErr, ok := errno.AsBusinessError(err); ok {
		log.Printf("model business error: method=%s path=%s code=%d message=%s", c.Request.Method, c.Request.URL.Path, bizErr.Code, bizErr.Message)
		common.Error(c, bizErr.Code, bizErr.Message)
		return
	}
	log.Printf("model system error: method=%s path=%s err=%v", c.Request.Method, c.Request.URL.Path, err)
	common.Error(c, errno.SystemError.Code, errno.SystemError.Message)
}

func (m *ModelController) handleBindError(c *gin.Context, action string, err error) {
	log.Printf("model bind request failed: action=%s method=%s path=%s err=%v", action, c.Request.Method, c.Request.URL.Path, err)
	common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
}
