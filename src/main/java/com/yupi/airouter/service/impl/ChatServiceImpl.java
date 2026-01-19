package com.yupi.airouter.service.impl;

import cn.hutool.core.util.IdUtil;
import cn.hutool.core.util.StrUtil;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.model.dto.chat.ChatMessage;
import com.yupi.airouter.model.dto.chat.ChatRequest;
import com.yupi.airouter.model.dto.chat.ChatResponse;
import com.yupi.airouter.model.dto.chat.StreamChunk;
import com.yupi.airouter.model.dto.log.RequestLogDTO;
import com.yupi.airouter.model.entity.Model;
import com.yupi.airouter.model.entity.ModelProvider;
import com.yupi.airouter.service.*;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Flux;

import java.util.List;

/**
 * 聊天服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Slf4j
@Service
public class ChatServiceImpl implements ChatService {

    @Resource
    private RoutingService routingService;

    @Resource
    private ModelInvokeService modelInvokeService;

    @Resource
    private ModelProviderService modelProviderService;

    @Resource
    private RequestLogService requestLogService;

    /**
     * 最大 Fallback 重试次数
     */
    private static final int MAX_FALLBACK_RETRIES = 3;

    @Override
    public ChatResponse chat(ChatRequest chatRequest, Long userId, Long apiKeyId, String clientIp, String userAgent) {
        long startTime = System.currentTimeMillis();
        String requestedModel = chatRequest.getModel();
        String traceId = IdUtil.simpleUUID();

        // 确定路由策略：优先使用请求中指定的策略，否则根据是否指定模型决定
        String strategyType = determineStrategyType(chatRequest.getRoutingStrategy(), requestedModel);

        // 选择模型
        Model selectedModel = routingService.selectModel(strategyType, "chat", requestedModel);

        if (selectedModel == null) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "没有可用的模型");
        }

        // 获取 Fallback 模型列表
        List<Model> fallbackModels = routingService.getFallbackModels(strategyType, "chat", requestedModel);

        // 尝试调用主模型
        boolean isFallback = false;

        try {
            return invokeModelWithFallback(selectedModel, fallbackModels, chatRequest,
                    userId, apiKeyId, traceId, startTime, strategyType, isFallback, clientIp, userAgent);
        } catch (Exception e) {
            log.error("所有模型调用失败", e);
            long duration = System.currentTimeMillis() - startTime;
            // 记录失败日志
            requestLogService.logRequest(RequestLogDTO.builder()
                    .traceId(traceId)
                    .userId(userId)
                    .apiKeyId(apiKeyId)
                    .requestModel(requestedModel)
                    .requestType("chat")
                    .source(apiKeyId != null ? "api" : "web")
                    .duration((int) duration)
                    .status("failed")
                    .errorMessage(e.getMessage())
                    .errorCode("SYSTEM_ERROR")
                    .routingStrategy(strategyType)
                    .isFallback(false)
                    .clientIp(clientIp)
                    .userAgent(userAgent)
                    .build());
            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "调用模型失败: " + e.getMessage());
        }
    }

    /**
     * 带 Fallback 的模型调用
     */
    private ChatResponse invokeModelWithFallback(Model primaryModel, List<Model> fallbackModels,
                                                 ChatRequest chatRequest, Long userId, Long apiKeyId,
                                                 String traceId, long startTime, String strategyType, boolean isFallback,
                                                 String clientIp, String userAgent) {
        // 尝试调用主模型
        try {
            return callModel(primaryModel, chatRequest, userId, apiKeyId, traceId, startTime, strategyType, isFallback, clientIp, userAgent);
        } catch (Exception e) {
            log.warn("模型 {} 调用失败，尝试 Fallback", primaryModel.getModelKey(), e);

            // 如果主模型失败且有 Fallback 模型，尝试 Fallback
            if (fallbackModels != null && !fallbackModels.isEmpty()) {
                int retries = Math.min(fallbackModels.size(), MAX_FALLBACK_RETRIES);
                for (int i = 0; i < retries; i++) {
                    Model fallbackModel = fallbackModels.get(i);
                    try {
                        log.info("尝试 Fallback 模型: {}", fallbackModel.getModelKey());
                        return callModel(fallbackModel, chatRequest, userId, apiKeyId, traceId, startTime, strategyType, true, clientIp, userAgent);
                    } catch (Exception fallbackException) {
                        log.warn("Fallback 模型 {} 调用失败", fallbackModel.getModelKey(), fallbackException);
                        if (i == retries - 1) {
                            throw fallbackException;
                        }
                    }
                }
            }
            throw e;
        }
    }

    /**
     * 调用单个模型
     */
    private ChatResponse callModel(Model model, ChatRequest chatRequest, Long userId, Long apiKeyId,
                                   String traceId, long startTime, String strategyType, boolean isFallback,
                                   String clientIp, String userAgent) {
        // 获取提供者信息
        ModelProvider provider = modelProviderService.getById(model.getProviderId());
        if (provider == null) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "模型提供者不存在");
        }

        try {
            // 调用模型
            org.springframework.ai.chat.model.ChatResponse aiResponse =
                    modelInvokeService.invoke(model, provider, chatRequest);

            // 转换响应格式
            ChatResponse response = convertResponse(aiResponse, model.getModelKey());

            // 记录请求日志
            long duration = System.currentTimeMillis() - startTime;
            requestLogService.logRequest(RequestLogDTO.builder()
                    .traceId(traceId)
                    .userId(userId)
                    .apiKeyId(apiKeyId)
                    .modelId(model.getId())
                    .requestModel(model.getModelKey())
                    .requestType("chat")
                    .source(apiKeyId != null ? "api" : "web")
                    .promptTokens(response.getUsage().getPromptTokens())
                    .completionTokens(response.getUsage().getCompletionTokens())
                    .totalTokens(response.getUsage().getTotalTokens())
                    .duration((int) duration)
                    .status("success")
                    .routingStrategy(strategyType)
                    .isFallback(isFallback)
                    .clientIp(clientIp)
                    .userAgent(userAgent)
                    .build());

            return response;
        } catch (Exception e) {
            long duration = System.currentTimeMillis() - startTime;
            // 记录失败日志
            requestLogService.logRequest(RequestLogDTO.builder()
                    .traceId(traceId)
                    .userId(userId)
                    .apiKeyId(apiKeyId)
                    .modelId(model.getId())
                    .requestModel(model.getModelKey())
                    .requestType("chat")
                    .source(apiKeyId != null ? "api" : "web")
                    .duration((int) duration)
                    .status("failed")
                    .errorMessage(e.getMessage())
                    .errorCode("MODEL_ERROR")
                    .routingStrategy(strategyType)
                    .isFallback(isFallback)
                    .clientIp(clientIp)
                    .userAgent(userAgent)
                    .build());
            throw e;
        }
    }

    @Override
    public Flux<String> chatStream(ChatRequest chatRequest, Long userId, Long apiKeyId, String clientIp, String userAgent) {
        String requestedModel = chatRequest.getModel();
        long startTime = System.currentTimeMillis();
        String traceId = IdUtil.simpleUUID();

        // Token 计数器
        final int[] promptTokens = {0};
        final int[] completionTokens = {0};

        // 深度思考状态追踪
        final boolean[] thinkingStarted = {false};
        final boolean[] thinkingEnded = {false};
        final boolean isReasoningEnabled = chatRequest.getEnableReasoning() != null && chatRequest.getEnableReasoning();

        try {
            // 确定路由策略：优先使用请求中指定的策略，否则根据是否指定模型决定
            String strategyType = determineStrategyType(chatRequest.getRoutingStrategy(), requestedModel);

            // 选择模型（策略内部查询数据库）
            Model selectedModel = routingService.selectModel(strategyType, "chat", requestedModel);

            if (selectedModel == null) {
                return Flux.error(new BusinessException(ErrorCode.PARAMS_ERROR, "没有可用的模型"));
            }

            // 获取提供者信息
            ModelProvider provider = modelProviderService.getById(selectedModel.getProviderId());
            if (provider == null) {
                return Flux.error(new BusinessException(ErrorCode.PARAMS_ERROR, "模型提供者不存在"));
            }

            // 调用流式模型，获取统一格式的响应块
            Flux<StreamChunk> chunkFlux = modelInvokeService.invokeStreamChunk(selectedModel, provider, chatRequest);

            // 将统一格式的响应块转换为输出字符串
            return chunkFlux.flatMap(chunk -> {
                // 更新 Token 统计
                if (chunk.getPromptTokens() != null && chunk.getPromptTokens() > 0) {
                    promptTokens[0] = chunk.getPromptTokens();
                }
                if (chunk.getCompletionTokens() != null && chunk.getCompletionTokens() > 0) {
                    completionTokens[0] = chunk.getCompletionTokens();
                }

                // 处理深度思考内容
                if (isReasoningEnabled && chunk.hasReasoningContent()) {
                    String reasoningContent = chunk.getReasoningContent();

                    // 思考内容存在且未开始
                    if (!thinkingStarted[0]) {
                        thinkingStarted[0] = true;
                        String thinkingText = "[THINKING]" + escapeNewlines(reasoningContent);
                        return Flux.just(thinkingText + "\n\n");
                    }

                    // 思考内容正在返回
                    if (!thinkingEnded[0]) {
                        return Flux.just(escapeNewlines(reasoningContent) + "\n\n");
                    }
                }

                // 思考内容结束，开始返回答案
                if (isReasoningEnabled && thinkingStarted[0] && !thinkingEnded[0] && !chunk.hasReasoningContent()) {
                    thinkingEnded[0] = true;
                    String endTag = "[/THINKING]\n\n";
                    if (chunk.hasText()) {
                        return Flux.just(endTag + escapeNewlines(chunk.getText()) + "\n\n");
                    }
                    return Flux.just(endTag);
                }

                // 返回普通文本内容
                if ((!isReasoningEnabled || thinkingEnded[0]) && chunk.hasText()) {
                    return Flux.just(escapeNewlines(chunk.getText()) + "\n\n");
                }

                return Flux.empty();
            }).doOnComplete(() -> {
                // 流结束时记录日志
                long duration = System.currentTimeMillis() - startTime;
                int totalTokens = promptTokens[0] + completionTokens[0];
                requestLogService.logRequest(RequestLogDTO.builder()
                        .traceId(traceId)
                        .userId(userId)
                        .apiKeyId(apiKeyId)
                        .modelId(selectedModel.getId())
                        .requestModel(selectedModel.getModelKey())
                        .requestType("chat")
                        .source(apiKeyId != null ? "api" : "web")
                        .promptTokens(promptTokens[0])
                        .completionTokens(completionTokens[0])
                        .totalTokens(totalTokens)
                        .duration((int) duration)
                        .status("success")
                        .routingStrategy(strategyType)
                        .isFallback(false)
                        .clientIp(clientIp)
                        .userAgent(userAgent)
                        .build());
            }).doOnError(error -> {
                // 流错误时记录日志
                log.error("流式调用模型失败", error);
                long duration = System.currentTimeMillis() - startTime;
                requestLogService.logRequest(RequestLogDTO.builder()
                        .traceId(traceId)
                        .userId(userId)
                        .apiKeyId(apiKeyId)
                        .modelId(selectedModel.getId())
                        .requestModel(selectedModel.getModelKey())
                        .requestType("chat")
                        .source(apiKeyId != null ? "api" : "web")
                        .duration((int) duration)
                        .status("failed")
                        .errorMessage(error.getMessage())
                        .errorCode("STREAM_ERROR")
                        .routingStrategy(strategyType)
                        .isFallback(false)
                        .clientIp(clientIp)
                        .userAgent(userAgent)
                        .build());
            });
        } catch (Exception e) {
            log.error("流式调用模型失败", e);
            return Flux.error(new BusinessException(ErrorCode.SYSTEM_ERROR, "流式调用模型失败: " + e.getMessage()));
        }
    }

    /**
     * 转义换行符
     */
    private String escapeNewlines(String text) {
        if (text == null) {
            return null;
        }
        return text.replace("\n", "\\n");
    }

    /**
     * 转换响应格式
     */
    private ChatResponse convertResponse(org.springframework.ai.chat.model.ChatResponse aiResponse, String modelName) {
        String content = aiResponse.getResult().getOutput().getText();

        ChatResponse.Usage usage = ChatResponse.Usage.builder()
                .promptTokens(aiResponse.getMetadata().getUsage().getPromptTokens() != null ?
                        aiResponse.getMetadata().getUsage().getPromptTokens() : 0)
                .completionTokens(aiResponse.getMetadata().getUsage().getCompletionTokens() != null ?
                        aiResponse.getMetadata().getUsage().getCompletionTokens() : 0)
                .totalTokens(aiResponse.getMetadata().getUsage().getTotalTokens() != null ?
                        aiResponse.getMetadata().getUsage().getTotalTokens() : 0)
                .build();

        ChatResponse.Choice choice = ChatResponse.Choice.builder()
                .index(0)
                .message(new ChatMessage("assistant", content))
                .finishReason(aiResponse.getResult().getMetadata().getFinishReason())
                .build();

        return ChatResponse.builder()
                .id(IdUtil.simpleUUID())
                .object("chat.completion")
                .created(System.currentTimeMillis() / 1000)
                .model(modelName)
                .choices(List.of(choice))
                .usage(usage)
                .build();
    }

    /**
     * 确定路由策略类型
     *
     * @param requestedStrategy 请求中指定的策略
     * @param requestedModel    请求中指定的模型
     * @return 最终使用的策略类型
     */
    private String determineStrategyType(String requestedStrategy, String requestedModel) {
        // 如果请求中指定了策略，优先使用
        if (StrUtil.isNotBlank(requestedStrategy)) {
            return requestedStrategy;
        }
        // 如果指定了模型但没有指定策略，使用固定模型策略
        if (StrUtil.isNotBlank(requestedModel)) {
            return "fixed";
        }
        // 默认使用自动路由策略
        return "auto";
    }
}
