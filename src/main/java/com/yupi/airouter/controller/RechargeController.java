package com.yupi.airouter.controller;

import com.mybatisflex.core.paginate.Page;
import com.stripe.model.checkout.Session;
import com.yupi.airouter.annotation.AuthCheck;
import com.yupi.airouter.common.BaseResponse;
import com.yupi.airouter.common.ResultUtils;
import com.yupi.airouter.config.StripeConfig;
import com.yupi.airouter.constant.UserConstant;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.model.entity.RechargeRecord;
import com.yupi.airouter.model.entity.User;
import com.yupi.airouter.service.RechargeService;
import com.yupi.airouter.service.StripePaymentService;
import com.yupi.airouter.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import jakarta.annotation.Resource;
import jakarta.servlet.http.HttpServletRequest;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.*;

import java.io.Serializable;
import java.math.BigDecimal;

/**
 * 充值控制器
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@RestController
@RequestMapping("/recharge")
@Slf4j
public class RechargeController {

    @Resource
    private RechargeService rechargeService;

    @Resource
    private StripePaymentService stripePaymentService;

    @Resource
    private UserService userService;

    @Resource
    private StripeConfig stripeConfig;

    /**
     * 创建充值订单（Stripe）
     */
    @PostMapping("/stripe/create")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "创建Stripe充值订单")
    public BaseResponse<CreateRechargeResponse> createStripeRecharge(
            @RequestBody CreateRechargeRequest request,
            HttpServletRequest httpRequest) {

        if (request.getAmount() == null || request.getAmount().compareTo(BigDecimal.ZERO) <= 0) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "充值金额必须大于0");
        }

        // 设置最小充值金额1元，最大10000元
        if (request.getAmount().compareTo(BigDecimal.ONE) < 0 ||
            request.getAmount().compareTo(new BigDecimal("10000")) > 0) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "充值金额必须在1-10000元之间");
        }

        User loginUser = userService.getLoginUser(httpRequest);

        // 从配置获取回调URL
        String successUrl = stripeConfig.getSuccessUrl();
        String cancelUrl = stripeConfig.getCancelUrl();

        // 创建 Stripe 支付会话
        Session session = stripePaymentService.createCheckoutSession(
                loginUser.getId(),
                request.getAmount(),
                successUrl,
                cancelUrl
        );

        CreateRechargeResponse response = new CreateRechargeResponse();
        response.setCheckoutUrl(session.getUrl());
        response.setSessionId(session.getId());

        return ResultUtils.success(response);
    }

    /**
     * 充值成功回调
     */
    @GetMapping("/stripe/success")
    @Operation(summary = "Stripe充值成功回调")
    public BaseResponse<String> stripeSuccess(@RequestParam("session_id") String sessionId) {
        log.info("收到Stripe支付成功回调，SessionID: {}", sessionId);

        boolean success = stripePaymentService.handlePaymentSuccess(sessionId);
        if (success) {
            return ResultUtils.success("充值成功！");
        } else {
            throw new BusinessException(ErrorCode.OPERATION_ERROR, "充值处理失败");
        }
    }

    /**
     * 充值取消回调
     */
    @GetMapping("/stripe/cancel")
    @Operation(summary = "Stripe充值取消回调")
    public BaseResponse<String> stripeCancel() {
        log.info("用户取消了Stripe支付");
        return ResultUtils.success("您取消了充值");
    }

    /**
     * 获取我的充值记录（分页）
     */
    @GetMapping("/list/my")
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "获取我的充值记录")
    public BaseResponse<Page<RechargeRecord>> getMyRechargeRecords(
            @RequestParam(defaultValue = "1") int pageNum,
            @RequestParam(defaultValue = "10") int pageSize,
            HttpServletRequest request) {
        User loginUser = userService.getLoginUser(request);
        Page<RechargeRecord> page = rechargeService.listUserRechargeRecords(
                loginUser.getId(), pageNum, pageSize);
        return ResultUtils.success(page);
    }

    /**
     * 创建充值请求
     */
    @Data
    public static class CreateRechargeRequest {
        private BigDecimal amount; // 充值金额（元）
    }

    /**
     * 创建充值响应
     */
    @Data
    @NoArgsConstructor
    @AllArgsConstructor
    public static class CreateRechargeResponse implements Serializable {
        private String checkoutUrl; // Stripe支付页面URL
        private String sessionId; // Stripe会话ID
    }
}
