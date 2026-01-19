package com.yupi.airouter.model.dto.chat;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.util.List;

/**
 * 聊天请求（兼容 OpenAI 格式）
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
public class ChatRequest implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 模型名称（如：qwen-plus）
     */
    private String model;

    /**
     * 消息列表
     */
    private List<ChatMessage> messages;

    /**
     * 是否流式返回
     */
    private Boolean stream = false;

    /**
     * 温度参数（0-1）
     */
    private Double temperature;

    /**
     * 最大生成Token数
     */
    @JsonProperty("max_tokens")
    private Integer maxTokens;
}
