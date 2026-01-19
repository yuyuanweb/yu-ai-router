package com.yupi.airouter.adapter;

import ai.z.openapi.ZhipuAiClient;
import ai.z.openapi.service.model.*;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.model.dto.chat.ChatRequest;
import com.yupi.airouter.model.dto.chat.StreamChunk;
import com.yupi.airouter.model.entity.Model;
import com.yupi.airouter.model.entity.ModelProvider;
import io.reactivex.rxjava3.core.Flowable;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.ObjectUtils;
import org.springframework.ai.chat.messages.AssistantMessage;
import org.springframework.ai.chat.model.ChatResponse;
import org.springframework.ai.chat.model.Generation;
import org.springframework.ai.chat.metadata.ChatResponseMetadata;
import org.springframework.ai.chat.metadata.DefaultUsage;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Flux;

import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * 智谱 AI 模型适配器
 * 使用官方 zai-sdk 依赖
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 * @see <a href="https://docs.bigmodel.cn/cn/api/introduction#java-sdk">智谱AI官方文档</a>
 */
@Slf4j
@Component
public class ZhipuAIAdapter implements ModelAdapter {

    /**
     * 支持的提供者名称
     */
    private static final Set<String> SUPPORTED_PROVIDERS = Set.of("zhipu", "zhipuai", "glm");

    @Override
    public ChatResponse invoke(Model model, ModelProvider provider, ChatRequest chatRequest) {
        try {
            ZhipuAiClient client = createClient(provider);
            ChatCompletionCreateParams request = buildRequest(chatRequest, model, false);

            ChatCompletionResponse response = client.chat().createChatCompletion(request);

            return convertToChatResponse(response);
        } catch (Exception e) {
            log.error("智谱 AI 适配器调用模型 {} 失败", model.getModelKey(), e);
            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "调用模型失败: " + e.getMessage());
        }
    }

    @Override
    public Flux<ChatResponse> invokeStream(Model model, ModelProvider provider, ChatRequest chatRequest) {
        try {
            ZhipuAiClient client = createClient(provider);
            ChatCompletionCreateParams request = buildRequest(chatRequest, model, true);

            ChatCompletionResponse response = client.chat().createChatCompletion(request);

            Flowable<ModelData> responseFlowable = response.getFlowable();
            if (ObjectUtils.isEmpty(responseFlowable)) {
                return Flux.empty();
            }

            // 将 RxJava Flowable 转换为 Reactor Flux
            return Flux.from(responseFlowable)
                    .map(this::convertModelDataToChatResponse);
        } catch (Exception e) {
            log.error("智谱 AI 适配器流式调用模型 {} 失败", model.getModelKey(), e);
            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "流式调用模型失败: " + e.getMessage());
        }
    }

    @Override
    public Flux<StreamChunk> invokeStreamChunk(Model model, ModelProvider provider, ChatRequest chatRequest) {
        try {
            ZhipuAiClient client = createClient(provider);
            ChatCompletionCreateParams request = buildRequest(chatRequest, model, true);

            ChatCompletionResponse response = client.chat().createChatCompletion(request);

            Flowable<ModelData> responseFlowable = response.getFlowable();
            if (ObjectUtils.isEmpty(responseFlowable)) {
                return Flux.empty();
            }

            // 将 RxJava Flowable 转换为 Reactor Flux，并转换为统一的 StreamChunk 格式
            return Flux.from(responseFlowable)
                    .map(this::convertToStreamChunk)
                    .filter(chunk -> !chunk.isEmpty());
        } catch (Exception e) {
            log.error("智谱 AI 适配器流式调用模型 {} 失败", model.getModelKey(), e);
            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "流式调用模型失败: " + e.getMessage());
        }
    }

    /**
     * 将智谱 AI 流式响应块转换为统一的 StreamChunk 格式
     */
    private StreamChunk convertToStreamChunk(ModelData modelData) {
        StreamChunk.StreamChunkBuilder builder = StreamChunk.builder();

        // 提取 Token 信息
        if (modelData.getUsage() != null) {
            int promptTokens = modelData.getUsage().getPromptTokens();
            int completionTokens = modelData.getUsage().getCompletionTokens();
            if (promptTokens > 0) {
                builder.promptTokens(promptTokens);
            }
            if (completionTokens > 0) {
                builder.completionTokens(completionTokens);
            }
            log.debug("智谱 AI Token 统计 - promptTokens: {}, completionTokens: {}", promptTokens, completionTokens);
        }

        // 提取内容
        if (modelData.getChoices() != null && !modelData.getChoices().isEmpty()) {
            Choice choice = modelData.getChoices().getFirst();

            // 流式响应使用 delta 字段
            if (choice.getDelta() != null) {
                Delta delta = choice.getDelta();

                // 提取深度思考内容
                String reasoningContent = delta.getReasoningContent();
                if (reasoningContent != null && !reasoningContent.isEmpty()) {
                    builder.reasoningContent(reasoningContent);
                }

                // 提取普通文本内容
                String content = delta.getContent();
                if (content != null && !content.isEmpty()) {
                    builder.text(content);
                }
            }
        }

        StreamChunk streamChunk = builder.build();
        // 如果没有任何内容，标记为空
        if (!streamChunk.hasText() && !streamChunk.hasReasoningContent()
                && streamChunk.getPromptTokens() == null && streamChunk.getCompletionTokens() == null) {
            return StreamChunk.empty();
        }
        return streamChunk;
    }

    /**
     * 将智谱 AI 同步响应转换为 Spring AI ChatResponse 格式
     */
    private ChatResponse convertToChatResponse(ChatCompletionResponse response) {
        if (response == null || response.getData() == null) {
            return null;
        }

        return convertModelDataToChatResponse(response.getData());
    }

    /**
     * 将 ModelData 转换为 Spring AI ChatResponse 格式
     */
    private ChatResponse convertModelDataToChatResponse(ModelData modelData) {
        if (modelData == null) {
            return null;
        }

        // 提取内容
        String content = "";
        if (modelData.getChoices() != null && !modelData.getChoices().isEmpty()) {
            Choice choice = modelData.getChoices().getFirst();

            // 同步响应使用 message 字段，流式响应使用 delta 字段
            if (choice.getMessage() != null && choice.getMessage().getContent() != null) {
                // content 是 Object 类型，需要转换为 String
                content = choice.getMessage().getContent().toString();
            } else if (choice.getDelta() != null && choice.getDelta().getContent() != null) {
                content = choice.getDelta().getContent();
            }
        }

        // 创建 Generation
        AssistantMessage assistantMessage = new AssistantMessage(content);
        Generation generation = new Generation(assistantMessage);

        // 创建 Usage
        DefaultUsage usage = null;
        if (modelData.getUsage() != null) {
            usage = new DefaultUsage(
                    modelData.getUsage().getPromptTokens(),
                    modelData.getUsage().getCompletionTokens()
            );
        }

        // 构建 ChatResponseMetadata
        ChatResponseMetadata metadata = ChatResponseMetadata.builder()
                .usage(usage)
                .build();

        return new ChatResponse(List.of(generation), metadata);
    }

    @Override
    public boolean supports(String providerName) {
        return providerName != null && SUPPORTED_PROVIDERS.contains(providerName.toLowerCase());
    }

    /**
     * 创建智谱 AI 客户端
     */
    private ZhipuAiClient createClient(ModelProvider provider) {
        return ZhipuAiClient.builder()
                .ofZHIPU()
                .apiKey(provider.getApiKey())
                .build();
    }

    /**
     * 构建聊天请求参数
     */
    private ChatCompletionCreateParams buildRequest(ChatRequest chatRequest, Model model, boolean stream) {
        // 转换消息列表
        List<ChatMessage> messages = chatRequest.getMessages().stream()
                .map(this::convertMessage)
                .collect(Collectors.toList());

        // 使用 var 或者直接链式调用来避免泛型问题
        ChatCompletionCreateParams.ChatCompletionCreateParamsBuilder<?, ?> builder = ChatCompletionCreateParams.builder()
                .model(model.getModelKey())
                .messages(messages)
                .stream(stream);

        // 设置温度（SDK 要求 Float 类型）
        if (chatRequest.getTemperature() != null) {
            builder.temperature(chatRequest.getTemperature().floatValue());
        }

        // 设置最大 token
        if (chatRequest.getMaxTokens() != null) {
            builder.maxTokens(chatRequest.getMaxTokens());
        }

        // 启用深度思考
        if (chatRequest.getEnableReasoning() != null && chatRequest.getEnableReasoning()) {
            ChatThinking thinking = ChatThinking.builder()
                    .type(ChatThinkingType.ENABLED.value())
                    .build();
            builder.thinking(thinking);
        }

        return builder.build();
    }

    /**
     * 转换消息类型
     */
    private ChatMessage convertMessage(com.yupi.airouter.model.dto.chat.ChatMessage msg) {
        String role = msg.getRole();
        ChatMessageRole messageRole;

        if ("system".equals(role)) {
            messageRole = ChatMessageRole.SYSTEM;
        } else if ("assistant".equals(role)) {
            messageRole = ChatMessageRole.ASSISTANT;
        } else {
            messageRole = ChatMessageRole.USER;
        }

        return ChatMessage.builder()
                .role(messageRole.value())
                .content(msg.getContent())
                .build();
    }
}
