package com.yupi.airouter.model.dto.chat;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;
import java.util.List;

/**
 * 聊天响应（兼容 OpenAI 格式）
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ChatResponse implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 响应ID
     */
    private String id;

    /**
     * 对象类型
     */
    private String object = "chat.completion";

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
    private List<Choice> choices;

    /**
     * Token使用情况
     */
    private Usage usage;

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class Choice {
        /**
         * 索引
         */
        private Integer index;

        /**
         * 消息
         */
        private ChatMessage message;

        /**
         * 结束原因
         */
        private String finishReason;
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class Usage {
        /**
         * 输入Token数
         */
        private Integer promptTokens;

        /**
         * 输出Token数
         */
        private Integer completionTokens;

        /**
         * 总Token数
         */
        private Integer totalTokens;
    }
}
