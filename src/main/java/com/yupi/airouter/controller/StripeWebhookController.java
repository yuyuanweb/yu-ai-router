package com.yupi.airouter.controller;

import com.stripe.model.Event;
import com.stripe.model.checkout.Session;
import com.yupi.airouter.service.StripePaymentService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.*;

/**
 * Stripe Webhook 控制器
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@RestController
@RequestMapping("/webhook/stripe")
@Tag(name = "Stripe Webhook")
@Slf4j
public class StripeWebhookController {

    @Resource
    private StripePaymentService stripePaymentService;

    /**
     * Stripe Webhook 回调
     * Stripe会在支付状态变更时调用这个接口
     */
    @PostMapping
    @Operation(summary = "Stripe Webhook回调")
    public String handleStripeWebhook(
            @RequestBody String payload,
            @RequestHeader("Stripe-Signature") String sigHeader) {

        log.info("收到Stripe Webhook回调");

        try {
            // 验证 Webhook 签名
            Event event = stripePaymentService.constructWebhookEvent(payload, sigHeader);

            // 处理不同的事件类型
            switch (event.getType()) {
                case "checkout.session.completed":
                    // 支付成功
                    Session session = (Session) event.getDataObjectDeserializer()
                            .getObject().orElse(null);
                    if (session != null) {
                        log.info("支付成功：SessionID {}", session.getId());
                        stripePaymentService.handlePaymentSuccess(session.getId());
                    }
                    break;

                case "checkout.session.expired":
                    // 支付会话过期
                    log.info("支付会话过期");
                    break;

                case "charge.refunded":
                    // 退款
                    log.info("订单退款");
                    break;

                default:
                    log.info("未处理的事件类型：{}", event.getType());
            }

            return "success";
        } catch (Exception e) {
            log.error("处理Stripe Webhook失败", e);
            return "error";
        }
    }
}
