package com.yupi.airouter.service;

import com.yupi.airouter.model.dto.chat.ChatRequest;
import com.yupi.airouter.model.dto.chat.ChatResponse;
import com.yupi.airouter.model.dto.chat.StreamResponse;
import reactor.core.publisher.Flux;

/**
 * 聊天服务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface ChatService {

    /**
     * 非流式聊天
     *
     * @param chatRequest 聊天请求
     * @param userId      用户ID
     * @param apiKeyId    API Key ID
     * @param clientIp    客户端IP
     * @param userAgent   User-Agent
     * @return 聊天响应
     */
    ChatResponse chat(ChatRequest chatRequest, Long userId, Long apiKeyId, String clientIp, String userAgent);

    /**
     * 流式聊天（返回 OpenAI SSE 格式）
     *
     * @param chatRequest 聊天请求
     * @param userId      用户ID
     * @param apiKeyId    API Key ID
     * @param clientIp    客户端IP
     * @param userAgent   User-Agent
     * @return 流式聊天响应
     */
    Flux<StreamResponse> chatStream(ChatRequest chatRequest, Long userId, Long apiKeyId, String clientIp, String userAgent);
}
