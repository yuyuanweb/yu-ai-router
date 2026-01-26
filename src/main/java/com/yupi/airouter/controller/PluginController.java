/**
 * 插件管理控制器
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.controller;

import com.yupi.airouter.annotation.AuthCheck;
import com.yupi.airouter.common.BaseResponse;
import com.yupi.airouter.common.ResultUtils;
import com.yupi.airouter.constant.UserConstant;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.model.dto.plugin.PluginExecuteRequest;
import com.yupi.airouter.model.dto.plugin.PluginUpdateRequest;
import com.yupi.airouter.model.entity.User;
import com.yupi.airouter.model.vo.PluginConfigVO;
import com.yupi.airouter.model.vo.PluginExecuteVO;
import com.yupi.airouter.service.PluginService;
import com.yupi.airouter.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import jakarta.annotation.Resource;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Slf4j
@RestController
@RequestMapping("/plugin")
public class PluginController {

    @Resource
    private PluginService pluginService;

    @Resource
    private UserService userService;

    /**
     * 获取所有插件列表（管理员）
     */
    @GetMapping("/list")
    @Operation(summary = "获取所有插件列表")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<List<PluginConfigVO>> listPlugins() {
        List<PluginConfigVO> pluginList = pluginService.listAllPlugins();
        return ResultUtils.success(pluginList);
    }

    /**
     * 获取启用的插件列表
     */
    @GetMapping("/list/enabled")
    @Operation(summary = "获取启用的插件列表")
    public BaseResponse<List<PluginConfigVO>> listEnabledPlugins() {
        List<PluginConfigVO> pluginList = pluginService.listEnabledPlugins();
        return ResultUtils.success(pluginList);
    }

    /**
     * 获取插件详情
     */
    @GetMapping("/get")
    @Operation(summary = "获取插件详情")
    public BaseResponse<PluginConfigVO> getPlugin(@RequestParam String pluginKey) {
        if (pluginKey == null || pluginKey.isEmpty()) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "插件标识不能为空");
        }
        PluginConfigVO pluginConfig = pluginService.getPluginByKey(pluginKey);
        if (pluginConfig == null) {
            throw new BusinessException(ErrorCode.NOT_FOUND_ERROR, "插件不存在");
        }
        return ResultUtils.success(pluginConfig);
    }

    /**
     * 更新插件配置（管理员）
     */
    @PostMapping("/update")
    @Operation(summary = "更新插件配置")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Boolean> updatePlugin(@RequestBody PluginUpdateRequest request) {
        if (request == null) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR);
        }
        boolean result = pluginService.updatePlugin(request);
        return ResultUtils.success(result);
    }

    /**
     * 启用插件（管理员）
     */
    @PostMapping("/enable")
    @Operation(summary = "启用插件")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Boolean> enablePlugin(@RequestParam String pluginKey) {
        if (pluginKey == null || pluginKey.isEmpty()) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "插件标识不能为空");
        }
        boolean result = pluginService.enablePlugin(pluginKey);
        return ResultUtils.success(result);
    }

    /**
     * 禁用插件（管理员）
     */
    @PostMapping("/disable")
    @Operation(summary = "禁用插件")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Boolean> disablePlugin(@RequestParam String pluginKey) {
        if (pluginKey == null || pluginKey.isEmpty()) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "插件标识不能为空");
        }
        boolean result = pluginService.disablePlugin(pluginKey);
        return ResultUtils.success(result);
    }

    /**
     * 执行插件
     */
    @PostMapping("/execute")
    @Operation(summary = "执行插件")
    public BaseResponse<PluginExecuteVO> executePlugin(@RequestBody PluginExecuteRequest request,
                                                       HttpServletRequest httpRequest) {
        if (request == null) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR);
        }

        // 获取当前用户
        User loginUser = userService.getLoginUser(httpRequest);
        Long userId = loginUser != null ? loginUser.getId() : null;

        PluginExecuteVO result = pluginService.executePlugin(request, userId);
        return ResultUtils.success(result);
    }

    /**
     * 重新加载插件（管理员）
     */
    @PostMapping("/reload")
    @Operation(summary = "重新加载插件")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Boolean> reloadPlugin(@RequestParam String pluginKey) {
        if (pluginKey == null || pluginKey.isEmpty()) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "插件标识不能为空");
        }
        pluginService.reloadPlugin(pluginKey);
        return ResultUtils.success(true);
    }

    /**
     * 重新加载所有插件（管理员）
     */
    @PostMapping("/reload/all")
    @Operation(summary = "重新加载所有插件")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Boolean> reloadAllPlugins() {
        pluginService.initPlugins();
        return ResultUtils.success(true);
    }
}
