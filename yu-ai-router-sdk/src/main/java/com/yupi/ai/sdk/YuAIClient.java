/**
 * Yu AI 客户端主入口
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.ai.sdk;

import com.yupi.ai.sdk.callback.StreamCallback;
import com.yupi.ai.sdk.config.ClientConfig;
import com.yupi.ai.sdk.http.HttpClient;
import com.yupi.ai.sdk.model.ChatMessage;
import com.yupi.ai.sdk.model.ChatRequest;
import com.yupi.ai.sdk.model.ChatResponse;
import lombok.extern.slf4j.Slf4j;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * Yu AI 客户端 - SDK 主入口类
 * <p>
 * 使用示例：
 * <pre>
 * YuAIClient client = YuAIClient.builder()
 *     .apiKey("sk-xxx")
 *     .baseUrl("http://localhost:8123/api")
 *     .build();
 * 
 * // 同步调用
 * ChatResponse response = client.chat("你好");
 * System.out.println(response.getContent());
 * 
 * // 流式调用
 * client.chatStream("讲个故事", new StreamCallback() {
 *     public void onMessage(ChatChunk chunk) {
 *         System.out.print(chunk.getContent());
 *     }
 *     public void onComplete() {
 *         System.out.println("\n完成");
 *     }
 *     public void onError(Throwable e) {
 *         e.printStackTrace();
 *     }
 * });
 * </pre>
 */
@Slf4j
public class YuAIClient {

    private final ClientConfig config;
    private final HttpClient httpClient;

    private YuAIClient(ClientConfig config) {
        this.config = config;
        this.config.validate();
        this.httpClient = new HttpClient(config);
        log.info("YuAIClient initialized with baseUrl: {}", config.getBaseUrl());
    }

    /**
     * 创建客户端构建器
     */
    public static ClientConfigBuilder builder() {
        return new ClientConfigBuilder();
    }

    /**
     * 同步聊天 - 简单文本
     *
     * @param message 用户消息
     * @return 聊天响应
     */
    public ChatResponse chat(String message) {
        return chat(ChatRequest.simple(message));
    }

    /**
     * 同步聊天 - 指定模型
     *
     * @param model   模型名称
     * @param message 用户消息
     * @return 聊天响应
     */
    public ChatResponse chat(String model, String message) {
        return chat(ChatRequest.withModel(model, message));
    }

    /**
     * 同步聊天 - 完整请求
     *
     * @param request 聊天请求
     * @return 聊天响应
     */
    public ChatResponse chat(ChatRequest request) {
        log.debug("Sending chat request: {}", request);
        ChatResponse response = httpClient.chat(request);
        log.debug("Received chat response: {}", response);
        return response;
    }

    /**
     * 流式聊天 - 简单文本
     *
     * @param message  用户消息
     * @param callback 流式回调
     */
    public void chatStream(String message, StreamCallback callback) {
        chatStream(ChatRequest.simple(message), callback);
    }

    /**
     * 流式聊天 - 指定模型
     *
     * @param model    模型名称
     * @param message  用户消息
     * @param callback 流式回调
     */
    public void chatStream(String model, String message, StreamCallback callback) {
        chatStream(ChatRequest.withModel(model, message), callback);
    }

    /**
     * 流式聊天 - 完整请求
     *
     * @param request  聊天请求
     * @param callback 流式回调
     */
    public void chatStream(ChatRequest request, StreamCallback callback) {
        log.debug("Sending stream chat request: {}", request);
        httpClient.chatStream(request, callback);
    }

    /**
     * 关闭客户端，释放资源
     */
    public void close() {
        httpClient.close();
        log.info("YuAIClient closed");
    }

    /**
     * 客户端配置构建器
     */
    public static class ClientConfigBuilder {
        private String apiKey;
        private String baseUrl = "http://localhost:8123/api";
        private Integer connectTimeout = 10000;
        private Integer readTimeout = 30000;
        private Integer writeTimeout = 30000;
        private Integer maxRetries = 3;
        private Integer retryDelay = 1000;

        public ClientConfigBuilder apiKey(String apiKey) {
            this.apiKey = apiKey;
            return this;
        }

        public ClientConfigBuilder baseUrl(String baseUrl) {
            this.baseUrl = baseUrl;
            return this;
        }

        public ClientConfigBuilder connectTimeout(Integer connectTimeout) {
            this.connectTimeout = connectTimeout;
            return this;
        }

        public ClientConfigBuilder readTimeout(Integer readTimeout) {
            this.readTimeout = readTimeout;
            return this;
        }

        public ClientConfigBuilder writeTimeout(Integer writeTimeout) {
            this.writeTimeout = writeTimeout;
            return this;
        }

        public ClientConfigBuilder maxRetries(Integer maxRetries) {
            this.maxRetries = maxRetries;
            return this;
        }

        public ClientConfigBuilder retryDelay(Integer retryDelay) {
            this.retryDelay = retryDelay;
            return this;
        }

        public YuAIClient build() {
            ClientConfig config = ClientConfig.builder()
                    .apiKey(apiKey)
                    .baseUrl(baseUrl)
                    .connectTimeout(connectTimeout)
                    .readTimeout(readTimeout)
                    .writeTimeout(writeTimeout)
                    .maxRetries(maxRetries)
                    .retryDelay(retryDelay)
                    .build();
            return new YuAIClient(config);
        }
    }
}
