package com.yupi.airouter.controller;

import com.mybatisflex.core.paginate.Page;
import com.yupi.airouter.annotation.AuthCheck;
import com.yupi.airouter.common.BaseResponse;
import com.yupi.airouter.common.ResultUtils;
import com.yupi.airouter.constant.UserConstant;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.exception.ThrowUtils;
import com.yupi.airouter.model.dto.log.RequestLogQueryRequest;
import com.yupi.airouter.model.entity.RequestLog;
import com.yupi.airouter.model.entity.User;
import com.yupi.airouter.service.BillingService;
import com.yupi.airouter.service.RequestLogService;
import com.yupi.airouter.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import jakarta.annotation.Resource;
import jakarta.servlet.http.HttpServletRequest;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.web.bind.annotation.*;

import java.io.Serializable;
import java.math.BigDecimal;
import java.time.LocalDate;
import java.util.List;
import java.util.Map;

/**
 * 统计接口
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@RestController
@RequestMapping("/stats")
public class StatsController {

    @Resource
    private RequestLogService requestLogService;

    @Resource
    private UserService userService;

    @Resource
    private BillingService billingService;

    @Resource
    private com.yupi.airouter.service.QuotaService quotaService;

    /**
     * 获取我的 Token 消耗统计
     */
    @GetMapping("/my/tokens")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "获取我的 Token 消耗统计")
    public BaseResponse<TokenStatsVO> getMyTokenStats(HttpServletRequest request) {
        User loginUser = userService.getLoginUser(request);
        Long totalTokens = requestLogService.countUserTokens(loginUser.getId());

        TokenStatsVO vo = new TokenStatsVO();
        vo.setTotalTokens(totalTokens);

        return ResultUtils.success(vo);
    }

    /**
     * 获取我的请求日志
     */
    @GetMapping("/my/logs")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "获取我的请求日志")
    public BaseResponse<List<RequestLog>> getMyLogs(@RequestParam(defaultValue = "100") Integer limit,
                                                     HttpServletRequest request) {
        User loginUser = userService.getLoginUser(request);
        List<RequestLog> logs = requestLogService.listUserLogs(loginUser.getId(), limit);

        return ResultUtils.success(logs);
    }

    /**
     * 获取我的费用统计
     */
    @GetMapping("/my/cost")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "获取我的费用统计")
    public BaseResponse<CostStatsVO> getMyCostStats(HttpServletRequest request) {
        User loginUser = userService.getLoginUser(request);
        BigDecimal totalCost = billingService.getUserTotalCost(loginUser.getId());
        BigDecimal todayCost = billingService.getUserTodayCost(loginUser.getId());

        CostStatsVO vo = new CostStatsVO();
        vo.setTotalCost(totalCost);
        vo.setTodayCost(todayCost);

        return ResultUtils.success(vo);
    }

    /**
     * 获取我的综合统计数据（个人中心使用）
     */
    @GetMapping("/my/summary")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "获取我的综合统计数据")
    public BaseResponse<UserSummaryStatsVO> getMySummaryStats(HttpServletRequest request) {
        User loginUser = userService.getLoginUser(request);
        Long userId = loginUser.getId();

        UserSummaryStatsVO vo = new UserSummaryStatsVO();
        
        // Token 统计
        vo.setTotalTokens(requestLogService.countUserTokens(userId));
        
        // 配额信息
        vo.setTokenQuota(loginUser.getTokenQuota());
        vo.setUsedTokens(loginUser.getUsedTokens() != null ? loginUser.getUsedTokens() : 0L);
        vo.setRemainingQuota(quotaService.getRemainingQuota(userId));
        
        // 费用统计
        vo.setTotalCost(billingService.getUserTotalCost(userId));
        vo.setTodayCost(billingService.getUserTodayCost(userId));
        
        // 请求统计
        vo.setTotalRequests(requestLogService.countUserRequests(userId));
        vo.setSuccessRequests(requestLogService.countUserSuccessRequests(userId));

        return ResultUtils.success(vo);
    }

    /**
     * 获取我的每日统计数据（图表使用）
     */
    @GetMapping("/my/daily")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "获取我的每日统计数据")
    public BaseResponse<List<Map<String, Object>>> getMyDailyStats(
            @RequestParam(required = false) String startDate,
            @RequestParam(required = false) String endDate,
            HttpServletRequest request) {
        User loginUser = userService.getLoginUser(request);
        
        // 默认查询最近7天
        LocalDate end = endDate != null ? LocalDate.parse(endDate) : LocalDate.now();
        LocalDate start = startDate != null ? LocalDate.parse(startDate) : end.minusDays(6);
        
        List<Map<String, Object>> dailyStats = requestLogService.getUserDailyStats(
                loginUser.getId(), start, end);

        return ResultUtils.success(dailyStats);
    }

    // region 调用历史查询

    /**
     * 分页查询我的调用历史
     */
    @PostMapping("/history/my/page")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "分页查询我的调用历史")
    public BaseResponse<Page<RequestLog>> pageMyHistory(@RequestBody RequestLogQueryRequest queryRequest,
                                                         HttpServletRequest request) {
        ThrowUtils.throwIf(queryRequest == null, ErrorCode.PARAMS_ERROR);
        User loginUser = userService.getLoginUser(request);
        
        // 强制只查询当前用户的数据
        queryRequest.setUserId(loginUser.getId());
        
        Page<RequestLog> page = requestLogService.pageByQuery(queryRequest);
        return ResultUtils.success(page);
    }

    /**
     * 获取调用历史详情
     */
    @GetMapping("/history/detail")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "获取调用历史详情")
    public BaseResponse<RequestLog> getHistoryDetail(@RequestParam Long id, HttpServletRequest request) {
        ThrowUtils.throwIf(id == null || id <= 0, ErrorCode.PARAMS_ERROR);
        
        User loginUser = userService.getLoginUser(request);
        RequestLog requestLog = requestLogService.getById(id);
        
        ThrowUtils.throwIf(requestLog == null, ErrorCode.NOT_FOUND_ERROR);
        
        // 校验是否是当前用户的数据（非管理员只能查看自己的）
        if (!UserConstant.ADMIN_ROLE.equals(loginUser.getUserRole())) {
            ThrowUtils.throwIf(!loginUser.getId().equals(requestLog.getUserId()), 
                    ErrorCode.NO_AUTH_ERROR, "只能查看自己的调用历史");
        }
        
        return ResultUtils.success(requestLog);
    }

    /**
     * 分页查询所有调用历史（仅管理员）
     */
    @PostMapping("/history/page")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    @Operation(summary = "分页查询所有调用历史（仅管理员）")
    public BaseResponse<Page<RequestLog>> pageHistory(@RequestBody RequestLogQueryRequest queryRequest) {
        ThrowUtils.throwIf(queryRequest == null, ErrorCode.PARAMS_ERROR);
        Page<RequestLog> page = requestLogService.pageByQuery(queryRequest);
        return ResultUtils.success(page);
    }

    // endregion

    /**
     * Token 统计视图对象
     */
    @Data
    @NoArgsConstructor
    @AllArgsConstructor
    public static class TokenStatsVO implements Serializable {
        /**
         * 总Token数
         */
        private Long totalTokens;
    }

    /**
     * 费用统计视图对象
     */
    @Data
    @NoArgsConstructor
    @AllArgsConstructor
    public static class CostStatsVO implements Serializable {
        /**
         * 总消费金额（元）
         */
        private BigDecimal totalCost;

        /**
         * 今日消费金额（元）
         */
        private BigDecimal todayCost;
    }

    /**
     * 用户综合统计视图对象
     */
    @Data
    @NoArgsConstructor
    public static class UserSummaryStatsVO implements Serializable {
        /**
         * 总Token数（来自日志统计）
         */
        private Long totalTokens;
        
        /**
         * Token配额（-1表示无限制）
         */
        private Long tokenQuota;
        
        /**
         * 已使用Token数（来自用户表）
         */
        private Long usedTokens;
        
        /**
         * 剩余配额（-1表示无限制）
         */
        private Long remainingQuota;
        
        /**
         * 总消费金额（元）
         */
        private BigDecimal totalCost;
        
        /**
         * 今日消费金额（元）
         */
        private BigDecimal todayCost;
        
        /**
         * 总请求数
         */
        private Long totalRequests;
        
        /**
         * 成功请求数
         */
        private Long successRequests;
    }
}
