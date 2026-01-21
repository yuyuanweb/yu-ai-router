/**
 * 流式响应块
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.ai.sdk.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 流式响应块
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ChatChunk {

    /**
     * 响应内容
     */
    private String content;

    /**
     * 深度思考内容（如果有）
     */
    private String reasoningContent;

    /**
     * 是否完成
     */
    private Boolean done;

    /**
     * 模型名称
     */
    private String model;

    /**
     * Token 统计
     */
    private Integer promptTokens;
    private Integer completionTokens;
}
