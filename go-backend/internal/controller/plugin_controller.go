package controller

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type PluginController struct {
	pluginService *service.PluginService
	userService   *service.UserService
}

func NewPluginController(pluginService *service.PluginService, userService *service.UserService) *PluginController {
	return &PluginController{
		pluginService: pluginService,
		userService:   userService,
	}
}

func (p *PluginController) ListPlugins(c *gin.Context) {
	result, err := p.pluginService.ListAllPlugins()
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (p *PluginController) ListEnabledPlugins(c *gin.Context) {
	result, err := p.pluginService.ListEnabledPlugins()
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (p *PluginController) GetPlugin(c *gin.Context) {
	pluginKey := strings.TrimSpace(c.Query("pluginKey"))
	if pluginKey == "" {
		common.Error(c, errno.ParamsError.Code, "插件标识不能为空")
		return
	}
	result, err := p.pluginService.GetPluginByKey(pluginKey)
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (p *PluginController) UpdatePlugin(c *gin.Context) {
	var request dto.PluginUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("update plugin bind failed: %v", err)
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	result, err := p.pluginService.UpdatePlugin(request)
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (p *PluginController) EnablePlugin(c *gin.Context) {
	pluginKey := strings.TrimSpace(c.Query("pluginKey"))
	if pluginKey == "" {
		common.Error(c, errno.ParamsError.Code, "插件标识不能为空")
		return
	}
	result, err := p.pluginService.EnablePlugin(pluginKey)
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (p *PluginController) DisablePlugin(c *gin.Context) {
	pluginKey := strings.TrimSpace(c.Query("pluginKey"))
	if pluginKey == "" {
		common.Error(c, errno.ParamsError.Code, "插件标识不能为空")
		return
	}
	result, err := p.pluginService.DisablePlugin(pluginKey)
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (p *PluginController) ExecutePlugin(c *gin.Context) {
	var request dto.PluginExecuteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("execute plugin bind failed: %v", err)
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	var userID *int64
	if loginUser, err := p.userService.GetLoginUser(c); err == nil && loginUser != nil && loginUser.ID > 0 {
		userID = &loginUser.ID
	}
	result, err := p.pluginService.ExecutePlugin(request, userID)
	if err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (p *PluginController) ReloadPlugin(c *gin.Context) {
	pluginKey := strings.TrimSpace(c.Query("pluginKey"))
	if pluginKey == "" {
		common.Error(c, errno.ParamsError.Code, "插件标识不能为空")
		return
	}
	if err := p.pluginService.ReloadPlugin(pluginKey); err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, true)
}

func (p *PluginController) ReloadAllPlugins(c *gin.Context) {
	if err := p.pluginService.InitPlugins(); err != nil {
		p.handleError(c, err)
		return
	}
	common.Success(c, true)
}

func (p *PluginController) handleError(c *gin.Context, err error) {
	if bizErr, ok := errno.AsBusinessError(err); ok {
		common.Error(c, bizErr.Code, bizErr.Message)
		return
	}
	common.Error(c, errno.SystemError.Code, errno.SystemError.Message)
}
