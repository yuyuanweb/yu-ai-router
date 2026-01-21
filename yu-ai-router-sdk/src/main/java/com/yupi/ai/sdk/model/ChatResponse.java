/**
 * 聊天响应
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
 * 聊天响应
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ChatResponse {

    /**
     * 响应 ID
     */
    private String id;

    /**
     * 对象类型
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
    private List<Choice> choices;

    /**
     * Token 使用统计
     */
    private Usage usage;

    /**
     * 获取第一个选项的消息内容
     */
    public String getContent() {
        if (choices != null && !choices.isEmpty()) {
            Choice firstChoice = choices.get(0);
            if (firstChoice.getMessage() != null) {
                return firstChoice.getMessage().getContent();
            }
        }
        return null;
    }

    /**
     * 选项
     */
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class Choice {
        private Integer index;
        private ChatMessage message;
        private String finishReason;
    }

    /**
     * Token 使用统计
     */
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class Usage {
        private Integer promptTokens;
        private Integer completionTokens;
        private Integer totalTokens;
    }
}
