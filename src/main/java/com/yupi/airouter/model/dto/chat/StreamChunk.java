package com.yupi.airouter.model.dto.chat;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 统一的流式响应块
 * 不同模型适配器将各自的响应转换为此统一格式
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class StreamChunk {

    /**
     * 普通文本内容
     */
    private String text;

    /**
     * 深度思考内容（如果有）
     */
    private String reasoningContent;

    /**
     * 输入 Token 数量
     */
    private Integer promptTokens;

    /**
     * 输出 Token 数量
     */
    private Integer completionTokens;

    /**
     * 是否为空内容（跳过处理）
     */
    @Builder.Default
    private boolean empty = false;

    /**
     * 创建空的响应块
     */
    public static StreamChunk empty() {
        return StreamChunk.builder().empty(true).build();
    }

    /**
     * 创建包含普通文本的响应块
     */
    public static StreamChunk ofText(String text) {
        return StreamChunk.builder().text(text).build();
    }

    /**
     * 创建包含思考内容的响应块
     */
    public static StreamChunk ofReasoning(String reasoningContent) {
        return StreamChunk.builder().reasoningContent(reasoningContent).build();
    }

    /**
     * 判断是否有思考内容
     */
    public boolean hasReasoningContent() {
        return reasoningContent != null && !reasoningContent.isEmpty();
    }

    /**
     * 判断是否有普通文本
     */
    public boolean hasText() {
        return text != null && !text.isEmpty();
    }
}
