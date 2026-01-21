/**
 * 客户端配置
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.ai.sdk.config;

import lombok.Builder;
import lombok.Data;

/**
 * Yu AI 客户端配置
 */
@Data
@Builder
public class ClientConfig {

    /**
     * API Key
     */
    private String apiKey;

    /**
     * 基础 URL
     */
    @Builder.Default
    private String baseUrl = "http://localhost:8123/api";

    /**
     * 连接超时（毫秒）
     */
    @Builder.Default
    private Integer connectTimeout = 10000;

    /**
     * 读取超时（毫秒）
     */
    @Builder.Default
    private Integer readTimeout = 30000;

    /**
     * 写入超时（毫秒）
     */
    @Builder.Default
    private Integer writeTimeout = 30000;

    /**
     * 最大重试次数
     */
    @Builder.Default
    private Integer maxRetries = 3;

    /**
     * 重试延迟（毫秒）
     */
    @Builder.Default
    private Integer retryDelay = 1000;

    /**
     * 验证配置
     */
    public void validate() {
        if (apiKey == null || apiKey.trim().isEmpty()) {
            throw new IllegalArgumentException("API Key cannot be null or empty");
        }
        if (baseUrl == null || baseUrl.trim().isEmpty()) {
            throw new IllegalArgumentException("Base URL cannot be null or empty");
        }
    }
}
