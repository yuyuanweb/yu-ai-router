package com.yupi.airouter.controller;

import cn.hutool.core.bean.BeanUtil;
import com.mybatisflex.core.paginate.Page;
import com.yupi.airouter.annotation.AuthCheck;
import com.yupi.airouter.common.BaseResponse;
import com.yupi.airouter.common.DeleteRequest;
import com.yupi.airouter.common.ResultUtils;
import com.yupi.airouter.constant.UserConstant;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.exception.ThrowUtils;
import com.yupi.airouter.model.dto.user.*;
import com.yupi.airouter.model.entity.User;
import com.yupi.airouter.model.vo.LoginUserVO;
import com.yupi.airouter.model.vo.UserVO;
import com.yupi.airouter.service.BillingService;
import com.yupi.airouter.service.QuotaService;
import com.yupi.airouter.service.RequestLogService;
import com.yupi.airouter.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import jakarta.annotation.Resource;
import jakarta.servlet.http.HttpServletRequest;
import org.springframework.web.bind.annotation.*;

import java.math.BigDecimal;

import java.util.List;

/**
 * 用户控制层
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@RestController
@RequestMapping("/user")
public class UserController {

    @Resource
    private UserService userService;

    @Resource
    private QuotaService quotaService;

    @Resource
    private RequestLogService requestLogService;

    @Resource
    private BillingService billingService;

    /**
     * 用户注册
     *
     * @param userRegisterRequest 用户注册请求
     * @return 注册结果
     */
    @PostMapping("/register")
    public BaseResponse<Long> userRegister(@RequestBody UserRegisterRequest userRegisterRequest) {
        ThrowUtils.throwIf(userRegisterRequest == null, ErrorCode.PARAMS_ERROR);
        String userAccount = userRegisterRequest.getUserAccount();
        String userPassword = userRegisterRequest.getUserPassword();
        String checkPassword = userRegisterRequest.getCheckPassword();
        long result = userService.userRegister(userAccount, userPassword, checkPassword);
        return ResultUtils.success(result);
    }

    /**
     * 用户登录
     *
     * @param userLoginRequest 用户登录请求
     * @param request          请求对象
     * @return 脱敏后的用户登录信息
     */
    @PostMapping("/login")
    public BaseResponse<LoginUserVO> userLogin(@RequestBody UserLoginRequest userLoginRequest, HttpServletRequest request) {
        ThrowUtils.throwIf(userLoginRequest == null, ErrorCode.PARAMS_ERROR);
        String userAccount = userLoginRequest.getUserAccount();
        String userPassword = userLoginRequest.getUserPassword();
        LoginUserVO loginUserVO = userService.userLogin(userAccount, userPassword, request);
        return ResultUtils.success(loginUserVO);
    }

    @GetMapping("/get/login")
    public BaseResponse<LoginUserVO> getLoginUser(HttpServletRequest request) {
        User loginUser = userService.getLoginUser(request);
        return ResultUtils.success(userService.getLoginUserVO(loginUser));
    }

    /**
     * 用户注销
     *
     * @param request 请求对象
     * @return
     */
    @PostMapping("/logout")
    public BaseResponse<Boolean> userLogout(HttpServletRequest request) {
        ThrowUtils.throwIf(request == null, ErrorCode.PARAMS_ERROR);
        boolean result = userService.userLogout(request);
        return ResultUtils.success(result);
    }

    /**
     * 创建用户
     */
    @PostMapping("/add")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Long> addUser(@RequestBody UserAddRequest userAddRequest) {
        ThrowUtils.throwIf(userAddRequest == null, ErrorCode.PARAMS_ERROR);
        User user = new User();
        BeanUtil.copyProperties(userAddRequest, user);
        // 默认密码 12345678
        final String DEFAULT_PASSWORD = "12345678";
        String encryptPassword = userService.getEncryptPassword(DEFAULT_PASSWORD);
        user.setUserPassword(encryptPassword);
        boolean result = userService.save(user);
        ThrowUtils.throwIf(!result, ErrorCode.OPERATION_ERROR);
        return ResultUtils.success(user.getId());
    }

    /**
     * 根据 id 获取用户（仅管理员）
     */
    @GetMapping("/get")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<User> getUserById(long id) {
        ThrowUtils.throwIf(id <= 0, ErrorCode.PARAMS_ERROR);
        User user = userService.getById(id);
        ThrowUtils.throwIf(user == null, ErrorCode.NOT_FOUND_ERROR);
        return ResultUtils.success(user);
    }

    /**
     * 根据 id 获取包装类
     */
    @GetMapping("/get/vo")
    public BaseResponse<UserVO> getUserVOById(long id) {
        BaseResponse<User> response = getUserById(id);
        User user = response.getData();
        return ResultUtils.success(userService.getUserVO(user));
    }

    /**
     * 删除用户
     */
    @PostMapping("/delete")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Boolean> deleteUser(@RequestBody DeleteRequest deleteRequest) {
        if (deleteRequest == null || deleteRequest.getId() <= 0) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR);
        }
        boolean b = userService.removeById(deleteRequest.getId());
        return ResultUtils.success(b);
    }

    /**
     * 更新用户
     */
    @PostMapping("/update")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Boolean> updateUser(@RequestBody UserUpdateRequest userUpdateRequest) {
        if (userUpdateRequest == null || userUpdateRequest.getId() == null) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR);
        }
        User user = new User();
        BeanUtil.copyProperties(userUpdateRequest, user);
        boolean result = userService.updateById(user);
        ThrowUtils.throwIf(!result, ErrorCode.OPERATION_ERROR);
        return ResultUtils.success(true);
    }

    /**
     * 分页获取用户封装列表（仅管理员）
     *
     * @param userQueryRequest 查询请求参数
     */
    @PostMapping("/list/page/vo")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Page<UserVO>> listUserVOByPage(@RequestBody UserQueryRequest userQueryRequest) {
        ThrowUtils.throwIf(userQueryRequest == null, ErrorCode.PARAMS_ERROR);
        long pageNum = userQueryRequest.getPageNum();
        long pageSize = userQueryRequest.getPageSize();
        Page<User> userPage = userService.page(Page.of(pageNum, pageSize),
                userService.getQueryWrapper(userQueryRequest));
        // 数据脱敏
        Page<UserVO> userVOPage = new Page<>(pageNum, pageSize, userPage.getTotalRow());
        List<UserVO> userVOList = userService.getUserVOList(userPage.getRecords());
        userVOPage.setRecords(userVOList);
        return ResultUtils.success(userVOPage);
    }

    // region 配额管理

    /**
     * 获取我的配额信息
     */
    @GetMapping("/quota/my")
    @Operation(summary = "获取我的配额信息")
    public BaseResponse<QuotaVO> getMyQuota(HttpServletRequest request) {
        User loginUser = userService.getLoginUser(request);
        long remainingQuota = quotaService.getRemainingQuota(loginUser.getId());
        
        QuotaVO quotaVO = new QuotaVO();
        quotaVO.setTokenQuota(loginUser.getTokenQuota());
        quotaVO.setUsedTokens(loginUser.getUsedTokens() != null ? loginUser.getUsedTokens() : 0L);
        quotaVO.setRemainingQuota(remainingQuota);
        
        return ResultUtils.success(quotaVO);
    }

    /**
     * 设置用户配额（仅管理员）
     */
    @PostMapping("/quota/set")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    @Operation(summary = "设置用户配额")
    public BaseResponse<Boolean> setUserQuota(@RequestBody QuotaUpdateRequest quotaUpdateRequest) {
        ThrowUtils.throwIf(quotaUpdateRequest == null, ErrorCode.PARAMS_ERROR);
        ThrowUtils.throwIf(quotaUpdateRequest.getUserId() == null, ErrorCode.PARAMS_ERROR, "用户ID不能为空");
        ThrowUtils.throwIf(quotaUpdateRequest.getTokenQuota() == null, ErrorCode.PARAMS_ERROR, "配额不能为空");
        
        User user = new User();
        user.setId(quotaUpdateRequest.getUserId());
        user.setTokenQuota(quotaUpdateRequest.getTokenQuota());
        
        boolean result = userService.updateById(user);
        ThrowUtils.throwIf(!result, ErrorCode.OPERATION_ERROR);
        return ResultUtils.success(true);
    }

    /**
     * 重置用户已使用配额（仅管理员）
     */
    @PostMapping("/quota/reset")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    @Operation(summary = "重置用户已使用配额")
    public BaseResponse<Boolean> resetUserQuota(@RequestParam Long userId) {
        ThrowUtils.throwIf(userId == null || userId <= 0, ErrorCode.PARAMS_ERROR, "用户ID不合法");
        
        User user = new User();
        user.setId(userId);
        user.setUsedTokens(0L);
        
        boolean result = userService.updateById(user);
        ThrowUtils.throwIf(!result, ErrorCode.OPERATION_ERROR);
        return ResultUtils.success(true);
    }

    // endregion

    // region 用户状态管理

    /**
     * 禁用用户（仅管理员）
     */
    @PostMapping("/disable")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    @Operation(summary = "禁用用户")
    public BaseResponse<Boolean> disableUser(@RequestParam Long userId) {
        ThrowUtils.throwIf(userId == null || userId <= 0, ErrorCode.PARAMS_ERROR, "用户ID不合法");
        
        // 不能禁用自己
        // 可以在这里添加更多校验
        
        boolean result = userService.disableUser(userId);
        ThrowUtils.throwIf(!result, ErrorCode.OPERATION_ERROR);
        return ResultUtils.success(true);
    }

    /**
     * 启用用户（仅管理员）
     */
    @PostMapping("/enable")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    @Operation(summary = "启用用户")
    public BaseResponse<Boolean> enableUser(@RequestParam Long userId) {
        ThrowUtils.throwIf(userId == null || userId <= 0, ErrorCode.PARAMS_ERROR, "用户ID不合法");
        
        boolean result = userService.enableUser(userId);
        ThrowUtils.throwIf(!result, ErrorCode.OPERATION_ERROR);
        return ResultUtils.success(true);
    }

    // endregion

    // region 用户使用分析

    /**
     * 获取用户使用分析数据（仅管理员）
     */
    @GetMapping("/analysis")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    @Operation(summary = "获取用户使用分析数据")
    public BaseResponse<UserAnalysisVO> getUserAnalysis(@RequestParam Long userId) {
        ThrowUtils.throwIf(userId == null || userId <= 0, ErrorCode.PARAMS_ERROR, "用户ID不合法");
        
        User user = userService.getById(userId);
        ThrowUtils.throwIf(user == null, ErrorCode.NOT_FOUND_ERROR, "用户不存在");
        
        UserAnalysisVO analysisVO = new UserAnalysisVO();
        analysisVO.setUserId(userId);
        analysisVO.setUserAccount(user.getUserAccount());
        analysisVO.setUserName(user.getUserName());
        analysisVO.setUserStatus(user.getUserStatus());
        analysisVO.setUserRole(user.getUserRole());
        
        // 配额信息
        analysisVO.setTokenQuota(user.getTokenQuota());
        analysisVO.setUsedTokens(user.getUsedTokens() != null ? user.getUsedTokens() : 0L);
        analysisVO.setRemainingQuota(quotaService.getRemainingQuota(userId));
        
        // 请求统计
        analysisVO.setTotalRequests(requestLogService.countUserRequests(userId));
        analysisVO.setSuccessRequests(requestLogService.countUserSuccessRequests(userId));
        analysisVO.setTotalTokens(requestLogService.countUserTokens(userId));
        
        // 费用统计
        analysisVO.setTotalCost(billingService.getUserTotalCost(userId));
        analysisVO.setTodayCost(billingService.getUserTodayCost(userId));
        
        return ResultUtils.success(analysisVO);
    }

    // endregion

    /**
     * 用户使用分析视图对象
     */
    @lombok.Data
    public static class UserAnalysisVO implements java.io.Serializable {
        private Long userId;
        private String userAccount;
        private String userName;
        private String userStatus;
        private String userRole;
        private Long tokenQuota;
        private Long usedTokens;
        private Long remainingQuota;
        private Long totalRequests;
        private Long successRequests;
        private Long totalTokens;
        private BigDecimal totalCost;
        private BigDecimal todayCost;
    }

    /**
     * 配额信息视图对象
     */
    @lombok.Data
    public static class QuotaVO implements java.io.Serializable {
        /**
         * Token配额（-1表示无限制）
         */
        private Long tokenQuota;
        
        /**
         * 已使用Token数
         */
        private Long usedTokens;
        
        /**
         * 剩余配额（-1表示无限制）
         */
        private Long remainingQuota;
    }

    /**
     * 配额更新请求
     */
    @lombok.Data
    public static class QuotaUpdateRequest implements java.io.Serializable {
        /**
         * 用户ID
         */
        private Long userId;
        
        /**
         * Token配额（-1表示无限制）
         */
        private Long tokenQuota;
    }
}
