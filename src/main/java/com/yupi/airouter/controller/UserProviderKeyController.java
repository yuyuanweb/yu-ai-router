/**
 * 用户提供者密钥控制器（BYOK）
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.controller;

import com.yupi.airouter.annotation.AuthCheck;
import com.yupi.airouter.common.BaseResponse;
import com.yupi.airouter.common.DeleteRequest;
import com.yupi.airouter.common.ResultUtils;
import com.yupi.airouter.constant.UserConstant;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.model.dto.byok.UserProviderKeyAddRequest;
import com.yupi.airouter.model.dto.byok.UserProviderKeyUpdateRequest;
import com.yupi.airouter.model.entity.User;
import com.yupi.airouter.model.vo.UserProviderKeyVO;
import com.yupi.airouter.service.UserProviderKeyService;
import com.yupi.airouter.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import jakarta.annotation.Resource;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Slf4j
@RestController
@RequestMapping("/byok")
public class UserProviderKeyController {

    @Resource
    private UserProviderKeyService userProviderKeyService;

    @Resource
    private UserService userService;

    /**
     * 添加用户提供者密钥
     */
    @PostMapping("/add")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "添加用户提供者密钥")
    public BaseResponse<Boolean> addUserProviderKey(@RequestBody UserProviderKeyAddRequest request,
                                                     HttpServletRequest httpRequest) {
        User loginUser = userService.getLoginUser(httpRequest);
        boolean result = userProviderKeyService.addUserProviderKey(request, loginUser.getId());
        return ResultUtils.success(result);
    }

    /**
     * 更新用户提供者密钥
     */
    @PostMapping("/update")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "更新用户提供者密钥")
    public BaseResponse<Boolean> updateUserProviderKey(@RequestBody UserProviderKeyUpdateRequest request,
                                                        HttpServletRequest httpRequest) {
        User loginUser = userService.getLoginUser(httpRequest);
        boolean result = userProviderKeyService.updateUserProviderKey(request, loginUser.getId());
        return ResultUtils.success(result);
    }

    /**
     * 删除用户提供者密钥
     */
    @PostMapping("/delete")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "删除用户提供者密钥")
    public BaseResponse<Boolean> deleteUserProviderKey(@RequestBody DeleteRequest request,
                                                        HttpServletRequest httpRequest) {
        if (request == null || request.getId() == null) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR);
        }
        User loginUser = userService.getLoginUser(httpRequest);
        boolean result = userProviderKeyService.deleteUserProviderKey(request.getId(), loginUser.getId());
        return ResultUtils.success(result);
    }

    /**
     * 获取我的提供者密钥列表
     */
    @GetMapping("/my/list")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "获取我的提供者密钥列表")
    public BaseResponse<List<UserProviderKeyVO>> listMyProviderKeys(HttpServletRequest httpRequest) {
        User loginUser = userService.getLoginUser(httpRequest);
        List<UserProviderKeyVO> list = userProviderKeyService.listUserProviderKeys(loginUser.getId());
        return ResultUtils.success(list);
    }
}
