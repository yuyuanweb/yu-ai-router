package com.yupi.airouter.adapter;

import com.yupi.airouter.model.dto.chat.ChatRequest;
import com.yupi.airouter.model.dto.chat.StreamChunk;
import com.yupi.airouter.model.entity.Model;
import com.yupi.airouter.model.entity.ModelProvider;
import org.springframework.ai.chat.model.ChatResponse;
import reactor.core.publisher.Flux;

/**
 * 模型适配器接口
 * 不同的模型提供者实现该接口，封装各自的调用逻辑
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface ModelAdapter {

    /**
     * 同步调用模型
     *
     * @param model       模型信息
     * @param provider    提供者信息
     * @param chatRequest 聊天请求
     * @return 聊天响应
     */
    ChatResponse invoke(Model model, ModelProvider provider, ChatRequest chatRequest);

    /**
     * 流式调用模型（返回原始响应）
     *
     * @param model       模型信息
     * @param provider    提供者信息
     * @param chatRequest 聊天请求
     * @return 流式聊天响应
     */
    Flux<ChatResponse> invokeStream(Model model, ModelProvider provider, ChatRequest chatRequest);

    /**
     * 流式调用模型（返回统一格式的响应块）
     * 每个适配器负责将自己模型的响应转换为统一的 StreamChunk 格式
     *
     * @param model       模型信息
     * @param provider    提供者信息
     * @param chatRequest 聊天请求
     * @return 统一格式的流式响应
     */
    Flux<StreamChunk> invokeStreamChunk(Model model, ModelProvider provider, ChatRequest chatRequest);

    /**
     * 判断是否支持该提供者
     *
     * @param providerName 提供者名称
     * @return 是否支持
     */
    boolean supports(String providerName);
}
