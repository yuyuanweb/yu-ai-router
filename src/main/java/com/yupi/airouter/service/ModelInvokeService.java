package com.yupi.airouter.service;

import com.yupi.airouter.model.dto.chat.ChatRequest;
import com.yupi.airouter.model.dto.chat.StreamChunk;
import com.yupi.airouter.model.entity.Model;
import com.yupi.airouter.model.entity.ModelProvider;
import org.springframework.ai.chat.model.ChatResponse;
import reactor.core.publisher.Flux;

/**
 * 模型调用服务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface ModelInvokeService {

    /**
     * 调用指定的模型（非流式）
     *
     * @param model        模型信息
     * @param provider     提供者信息
     * @param chatRequest  聊天请求
     * @return AI 响应
     */
    ChatResponse invoke(Model model, ModelProvider provider, ChatRequest chatRequest);

    /**
     * 调用指定的模型（流式，返回原始响应）
     *
     * @param model        模型信息
     * @param provider     提供者信息
     * @param chatRequest  聊天请求
     * @return 流式AI响应
     */
    Flux<ChatResponse> invokeStream(Model model, ModelProvider provider, ChatRequest chatRequest);

    /**
     * 调用指定的模型（流式，返回统一格式的响应块）
     * 各适配器负责将模型响应转换为统一的 StreamChunk 格式
     *
     * @param model        模型信息
     * @param provider     提供者信息
     * @param chatRequest  聊天请求
     * @return 统一格式的流式响应
     */
    Flux<StreamChunk> invokeStreamChunk(Model model, ModelProvider provider, ChatRequest chatRequest);
}
