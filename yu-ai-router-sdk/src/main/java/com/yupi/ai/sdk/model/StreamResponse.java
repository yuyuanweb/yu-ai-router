/**
 * 流式响应（OpenAI SSE 格式）
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.ai.sdk.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

/**
 * 流式响应（OpenAI chat.completion.chunk 格式）
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class StreamResponse {

    /**
     * 响应 ID
     */
    private String id;

    /**
     * 对象类型: chat.completion.chunk
     */
    private String object;

    /**
     * 创建时间戳
     */
    private Long created;

    /**
     * 使用的模型
     */
    private String model;

    /**
     * 选项列表
     */
    private List<StreamChoice> choices;

    /**
     * 流式选项
     */
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class StreamChoice {
        /**
         * 索引
         */
        private Integer index;

        /**
         * 增量内容
         */
        private Delta delta;

        /**
         * 结束原因
         */
        private String finishReason;
    }

    /**
     * 增量内容
     */
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class Delta {
        /**
         * 角色（仅第一个块）
         */
        private String role;

        /**
         * 内容增量
         */
        private String content;

        /**
         * 深度思考内容（仅支持深度思考的模型）
         */
        private String reasoningContent;
    }
}
