package com.yupi.airouter.service;

import com.stripe.model.Event;
import com.stripe.model.checkout.Session;

import java.math.BigDecimal;

/**
 * Stripe 支付服务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface StripePaymentService {

    /**
     * 创建支付会话
     *
     * @param userId 用户ID
     * @param amount 充值金额（元）
     * @param successUrl 支付成功回调URL
     * @param cancelUrl 支付取消回调URL
     * @return Stripe支付会话
     */
    Session createCheckoutSession(Long userId, BigDecimal amount, String successUrl, String cancelUrl);

    /**
     * 处理支付成功回调
     *
     * @param sessionId Stripe会话ID
     * @return 是否处理成功
     */
    boolean handlePaymentSuccess(String sessionId);

    /**
     * 验证 Webhook 签名
     *
     * @param payload Webhook请求体
     * @param sigHeader Webhook签名
     * @return Stripe事件
     */
    Event constructWebhookEvent(String payload, String sigHeader);
}
