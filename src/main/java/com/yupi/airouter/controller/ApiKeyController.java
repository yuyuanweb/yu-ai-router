package com.yupi.airouter.controller;

import cn.hutool.core.bean.BeanUtil;
import com.mybatisflex.core.paginate.Page;
import com.mybatisflex.core.query.QueryWrapper;
import com.yupi.airouter.annotation.AuthCheck;
import com.yupi.airouter.common.BaseResponse;
import com.yupi.airouter.common.DeleteRequest;
import com.yupi.airouter.common.ResultUtils;
import com.yupi.airouter.constant.UserConstant;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.model.dto.apikey.ApiKeyCreateRequest;
import com.yupi.airouter.model.entity.ApiKey;
import com.yupi.airouter.model.entity.User;
import com.yupi.airouter.model.vo.ApiKeyVO;
import com.yupi.airouter.service.ApiKeyService;
import com.yupi.airouter.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import jakarta.annotation.Resource;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.stream.Collectors;

/**
 * API Key 接口
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@RestController
@RequestMapping("/api/key")
@Slf4j
public class ApiKeyController {

    @Resource
    private ApiKeyService apiKeyService;

    @Resource
    private UserService userService;

    /**
     * 创建 API Key
     */
    @PostMapping("/create")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "创建 API Key")
    public BaseResponse<ApiKeyVO> createApiKey(@RequestBody ApiKeyCreateRequest request,
                                                HttpServletRequest httpRequest) {
        User loginUser = userService.getLoginUser(httpRequest);
        ApiKey apiKey = apiKeyService.createApiKey(request.getKeyName(), loginUser);

        // 转换为 VO（完整显示 Key 值）
        ApiKeyVO apiKeyVO = BeanUtil.copyProperties(apiKey, ApiKeyVO.class);

        return ResultUtils.success(apiKeyVO);
    }

    /**
     * 获取我的 API Key 列表（分页）
     */
    @GetMapping("/list/my")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "获取我的 API Key 列表")
    public BaseResponse<Page<ApiKeyVO>> listMyApiKeys(
            @RequestParam(defaultValue = "1") long pageNum,
            @RequestParam(defaultValue = "10") long pageSize,
            HttpServletRequest request) {
        User loginUser = userService.getLoginUser(request);

        // 分页查询
        QueryWrapper queryWrapper = QueryWrapper.create()
                .eq("userId", loginUser.getId())
                .eq("isDelete", 0)
                .orderBy("createTime", false);

        Page<ApiKey> apiKeyPage = apiKeyService.page(Page.of(pageNum, pageSize), queryWrapper);

        // 转换为 VO（列表中部分隐藏 Key 值）
        Page<ApiKeyVO> apiKeyVOPage = new Page<>(pageNum, pageSize, apiKeyPage.getTotalRow());
        List<ApiKeyVO> apiKeyVOList = apiKeyPage.getRecords().stream()
                .map(apiKey -> {
                    ApiKeyVO vo = BeanUtil.copyProperties(apiKey, ApiKeyVO.class);
                    // 隐藏部分 Key 值（只显示前8位和后4位）
                    if (vo.getKeyValue() != null && vo.getKeyValue().length() > 12) {
                        String key = vo.getKeyValue();
                        vo.setKeyValue(key.substring(0, 8) + "****" + key.substring(key.length() - 4));
                    }
                    return vo;
                })
                .collect(Collectors.toList());

        apiKeyVOPage.setRecords(apiKeyVOList);
        return ResultUtils.success(apiKeyVOPage);
    }

    /**
     * 撤销 API Key
     */
    @PostMapping("/revoke")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "撤销 API Key")
    public BaseResponse<Boolean> revokeApiKey(@RequestBody DeleteRequest deleteRequest,
                                               HttpServletRequest request) {
        if (deleteRequest == null || deleteRequest.getId() == null) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR);
        }

        User loginUser = userService.getLoginUser(request);
        boolean result = apiKeyService.revokeApiKey(deleteRequest.getId(), loginUser.getId());

        return ResultUtils.success(result);
    }
}
