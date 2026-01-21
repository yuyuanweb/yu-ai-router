package com.yupi.airouter.model.dto.chat;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;
import java.util.List;

/**
 * 流式响应（兼容 OpenAI SSE 格式）
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class StreamResponse implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 响应ID
     */
    private String id;

    /**
     * 对象类型: chat.completion.chunk
     */
    private String object = "chat.completion.chunk";

    /**
     * 创建时间戳
     */
    private Long created;

    /**
     * 模型名称
     */
    private String model;

    /**
     * 选择列表
     */
    private List<StreamChoice> choices;

    /**
     * 流式选项（使用 delta 而非 message）
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
         * 增量内容（delta）
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
         * 角色（仅第一个块包含）
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
