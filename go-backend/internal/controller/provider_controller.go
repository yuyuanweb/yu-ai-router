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

type ProviderController struct {
	providerService *service.ProviderService
}

func NewProviderController(providerService *service.ProviderService) *ProviderController {
	return &ProviderController{providerService: providerService}
}

func (p *ProviderController) AddProvider(c *gin.Context) {
	var request dto.ProviderAddRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		p.handleBindError(c, "add provider", err)
		return
	}
	id, err := p.providerService.AddProvider(request)
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, strconv.FormatInt(id, 10))
}

func (p *ProviderController) DeleteProvider(c *gin.Context) {
	var request dto.DeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		p.handleBindError(c, "delete provider", err)
		return
	}
	if request.ID == nil || request.ID.Int64() <= 0 {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	result, err := p.providerService.DeleteProvider(request.ID.Int64())
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (p *ProviderController) UpdateProvider(c *gin.Context) {
	var request dto.ProviderUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		p.handleBindError(c, "update provider", err)
		return
	}
	result, err := p.providerService.UpdateProvider(request)
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (p *ProviderController) GetProviderVOByID(c *gin.Context) {
	id, err := parsePositiveID(c.Query("id"))
	if err != nil {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	provider, err := p.providerService.GetProviderByID(id)
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, p.providerService.GetProviderVO(provider))
}

func (p *ProviderController) ListProviderVOByPage(c *gin.Context) {
	var request dto.ProviderQueryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		p.handleBindError(c, "list provider by page", err)
		return
	}
	pageResponse, err := p.providerService.ListProviderVOByPage(request)
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, pageResponse)
}

func (p *ProviderController) ListProviderVO(c *gin.Context) {
	voList, err := p.providerService.ListProviderVO()
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, voList)
}

func (p *ProviderController) ListHealthyProviders(c *gin.Context) {
	voList, err := p.providerService.ListHealthyProviderVO()
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, voList)
}

func (p *ProviderController) handleError(c *gin.Context, err error) {
	if bizErr, ok := errno.AsBusinessError(err); ok {
		log.Printf("provider business error: method=%s path=%s code=%d message=%s", c.Request.Method, c.Request.URL.Path, bizErr.Code, bizErr.Message)
		common.Error(c, bizErr.Code, bizErr.Message)
		return
	}
	log.Printf("provider system error: method=%s path=%s err=%v", c.Request.Method, c.Request.URL.Path, err)
	common.Error(c, errno.SystemError.Code, errno.SystemError.Message)
}

func (p *ProviderController) handleBindError(c *gin.Context, action string, err error) {
	log.Printf("provider bind request failed: action=%s method=%s path=%s err=%v", action, c.Request.Method, c.Request.URL.Path, err)
	common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
}
