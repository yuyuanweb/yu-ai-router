package com.yupi.airouter.service.impl;

import com.stripe.exception.SignatureVerificationException;
import com.stripe.exception.StripeException;
import com.stripe.model.Event;
import com.stripe.model.checkout.Session;
import com.stripe.net.Webhook;
import com.stripe.param.checkout.SessionCreateParams;
import com.yupi.airouter.config.StripeConfig;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.model.entity.RechargeRecord;
import com.yupi.airouter.service.RechargeService;
import com.yupi.airouter.service.StripePaymentService;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;

/**
 * Stripe 支付服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Service
@Slf4j
public class StripePaymentServiceImpl implements StripePaymentService {

    @Resource
    private StripeConfig stripeConfig;

    @Resource
    private RechargeService rechargeService;

    @Override
    public Session createCheckoutSession(Long userId, BigDecimal amount, String successUrl, String cancelUrl) {
        try {
            // 创建充值记录
            RechargeRecord record = rechargeService.createRechargeRecord(userId, amount, "stripe");

            // 将金额转换为最小货币单位（分）
            long amountInCents = amount.multiply(new BigDecimal("100")).longValue();

            // 创建 Stripe Checkout Session
            SessionCreateParams params = SessionCreateParams.builder()
                    .setMode(SessionCreateParams.Mode.PAYMENT)
                    .setSuccessUrl(successUrl + "?session_id={CHECKOUT_SESSION_ID}")
                    .setCancelUrl(cancelUrl)
                    .addLineItem(
                            SessionCreateParams.LineItem.builder()
                                    .setPriceData(
                                            SessionCreateParams.LineItem.PriceData.builder()
                                                    // 人民币
                                                    .setCurrency("cny")
                                                    // 金额（分）
                                                    .setUnitAmount(amountInCents)
                                                    .setProductData(
                                                            SessionCreateParams.LineItem.PriceData.ProductData.builder()
                                                                    .setName("Yu AI Router 账户充值")
                                                                    .setDescription("充值金额：¥" + amount)
                                                                    .build()
                                                    )
                                                    .build()
                                    )
                                    .setQuantity(1L)
                                    .build()
                    )
                    .putMetadata("userId", userId.toString())
                    .putMetadata("recordId", record.getId().toString())
                    .putMetadata("amount", amount.toString())
                    .build();

            Session session = Session.create(params);
            log.info("创建 Stripe 支付会话成功：用户 {}, 金额 ¥{}, SessionID {}", userId, amount, session.getId());
            return session;

        } catch (StripeException e) {
            log.error("创建 Stripe 支付会话失败：{}", e.getMessage(), e);
            throw new BusinessException(ErrorCode.OPERATION_ERROR, "创建支付会话失败：" + e.getMessage());
        }
    }

    @Override
    public boolean handlePaymentSuccess(String sessionId) {
        try {
            // 获取会话信息
            Session session = Session.retrieve(sessionId);

            if (!"paid".equals(session.getPaymentStatus())) {
                log.warn("支付会话 {} 状态不是已支付：{}", sessionId, session.getPaymentStatus());
                return false;
            }

            // 从元数据中获取充值记录ID
            String recordIdStr = session.getMetadata().get("recordId");
            if (recordIdStr == null) {
                log.error("支付会话 {} 缺少 recordId 元数据", sessionId);
                return false;
            }

            Long recordId = Long.parseLong(recordIdStr);

            // 完成充值
            rechargeService.completeRecharge(recordId, sessionId);

            log.info("处理支付成功回调完成：SessionID {}", sessionId);
            return true;

        } catch (StripeException e) {
            log.error("获取支付会话失败：{}", e.getMessage(), e);
            throw new BusinessException(ErrorCode.OPERATION_ERROR, "获取支付会话失败");
        }
    }

    @Override
    public Event constructWebhookEvent(String payload, String sigHeader) {
        try {
            return Webhook.constructEvent(payload, sigHeader, stripeConfig.getWebhookSecret());
        } catch (SignatureVerificationException e) {
            log.error("Webhook 签名验证失败：{}", e.getMessage());
            throw new BusinessException(ErrorCode.FORBIDDEN_ERROR, "Webhook 签名验证失败");
        }
    }
}
