/**
 * 聊天请求
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.ai.sdk.model;

import com.google.gson.annotations.SerializedName;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.ArrayList;
import java.util.List;

/**
 * 聊天请求
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ChatRequest {

    /**
     * 模型名称
     */
    private String model;

    /**
     * 消息列表
     */
    private List<ChatMessage> messages;

    /**
     * 是否流式输出
     */
    @Builder.Default
    private Boolean stream = false;

    /**
     * 温度参数 0-1
     */
    private Double temperature;

    /**
     * 最大 Token 数
     */
    @SerializedName("max_tokens")
    private Integer maxTokens;

    /**
     * 是否启用深度思考
     */
    @SerializedName("enable_reasoning")
    private Boolean enableReasoning;

    /**
     * 路由策略
     */
    @SerializedName("routing_strategy")
    private String routingStrategy;

    /**
     * 创建简单的用户消息请求
     */
    public static ChatRequest simple(String userMessage) {
        List<ChatMessage> messages = new ArrayList<>();
        messages.add(ChatMessage.user(userMessage));
        return ChatRequest.builder()
                .messages(messages)
                .build();
    }

    /**
     * 创建带模型的用户消息请求
     */
    public static ChatRequest withModel(String model, String userMessage) {
        List<ChatMessage> messages = new ArrayList<>();
        messages.add(ChatMessage.user(userMessage));
        return ChatRequest.builder()
                .model(model)
                .messages(messages)
                .build();
    }
}
