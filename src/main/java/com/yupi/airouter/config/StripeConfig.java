package com.yupi.airouter.config;

import com.stripe.Stripe;
import jakarta.annotation.PostConstruct;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

/**
 * Stripe 支付配置
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Configuration
@ConfigurationProperties(prefix = "stripe")
@Data
@Slf4j
public class StripeConfig {

    /**
     * Stripe API 密钥
     */
    private String apiKey;

    /**
     * Webhook 签名密钥
     */
    private String webhookSecret;

    /**
     * 支付成功回调URL
     */
    private String successUrl;

    /**
     * 支付取消回调URL
     */
    private String cancelUrl;

    @PostConstruct
    public void init() {
        // 初始化 Stripe API Key（沙箱环境）
        Stripe.apiKey = apiKey;
        if (apiKey != null && apiKey.startsWith("sk_test_")) {
            log.info("Stripe API 初始化完成 - 沙箱环境（测试模式）");
            log.info("支付成功回调URL: {}", successUrl);
            log.info("支付取消回调URL: {}", cancelUrl);
        } else if (apiKey == null || "sk_test_YOUR_TEST_KEY".equals(apiKey)) {
            log.warn("Stripe API 密钥未配置，请在 application.yml 中配置测试密钥");
        } else {
            log.warn("检测到非测试密钥，请确保使用 sk_test_ 开头的测试密钥");
        }
    }
}
