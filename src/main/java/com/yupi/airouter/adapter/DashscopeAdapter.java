package com.yupi.airouter.adapter;

import com.alibaba.cloud.ai.dashscope.api.DashScopeApi;
import com.alibaba.cloud.ai.dashscope.chat.DashScopeChatModel;
import com.alibaba.cloud.ai.dashscope.chat.DashScopeChatOptions;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.model.dto.chat.ChatMessage;
import com.yupi.airouter.model.dto.chat.ChatRequest;
import com.yupi.airouter.model.dto.chat.StreamChunk;
import com.yupi.airouter.model.entity.Model;
import com.yupi.airouter.model.entity.ModelProvider;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.ObjectUtils;
import org.springframework.ai.chat.messages.AssistantMessage;
import org.springframework.ai.chat.messages.Message;
import org.springframework.ai.chat.messages.SystemMessage;
import org.springframework.ai.chat.messages.UserMessage;
import org.springframework.ai.chat.metadata.Usage;
import org.springframework.ai.chat.model.ChatResponse;
import org.springframework.ai.chat.prompt.Prompt;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Flux;

import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * 阿里云 DashScope（通义千问）模型适配器
 * 使用 spring-ai-alibaba-starter-dashscope 依赖
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Slf4j
@Component
public class DashscopeAdapter implements ModelAdapter {

    /**
     * 支持的提供者名称
     */
    private static final Set<String> SUPPORTED_PROVIDERS = Set.of("qwen", "dashscope", "tongyi", "aliyun");

    @Override
    public ChatResponse invoke(Model model, ModelProvider provider, ChatRequest chatRequest) {
        try {
            DashScopeChatModel chatModel = createChatModel(provider, model, chatRequest);
            Prompt prompt = buildPrompt(chatRequest);
            return chatModel.call(prompt);
        } catch (Exception e) {
            log.error("DashScope 适配器调用模型 {} 失败", model.getModelKey(), e);
            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "调用模型失败: " + e.getMessage());
        }
    }

    @Override
    public Flux<ChatResponse> invokeStream(Model model, ModelProvider provider, ChatRequest chatRequest) {
        try {
            DashScopeChatModel chatModel = createChatModel(provider, model, chatRequest);

            Prompt prompt = buildPrompt(chatRequest);
            return chatModel.stream(prompt);
        } catch (Exception e) {
            log.error("DashScope 适配器流式调用模型 {} 失败", model.getModelKey(), e);
            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "流式调用模型失败: " + e.getMessage());
        }
    }

    @Override
    public Flux<StreamChunk> invokeStreamChunk(Model model, ModelProvider provider, ChatRequest chatRequest) {
        return invokeStream(model, provider, chatRequest)
                .map(this::convertToStreamChunk)
                .filter(chunk -> !chunk.isEmpty());
    }

    /**
     * 将 DashScope（通义千问）响应转换为统一的 StreamChunk 格式
     * 通义千问深度思考内容存储在 metadata 的 "reasoningContent" 字段
     */
    private StreamChunk convertToStreamChunk(ChatResponse response) {
        StreamChunk.StreamChunkBuilder builder = StreamChunk.builder();

        // 提取 Token 信息
        if (ObjectUtils.isNotEmpty(response.getMetadata()) && response.getMetadata().getUsage() != null) {
            Usage dashScopeAiUsage = response.getMetadata().getUsage();
            Integer promptTokens = dashScopeAiUsage.getPromptTokens();
            Integer completionTokens = dashScopeAiUsage.getCompletionTokens();
            if (promptTokens != null && promptTokens > 0) {
                builder.promptTokens(promptTokens);
            }
            if (completionTokens != null && completionTokens > 0) {
                builder.completionTokens(completionTokens);
            }
        }

        // 提取内容
        if (ObjectUtils.isNotEmpty(response.getResult()) && ObjectUtils.isNotEmpty(response.getResult().getOutput())) {
            // 提取深度思考内容 - 通义千问使用 reasoningContent 字段
            AssistantMessage assistantMessage = response.getResult().getOutput();
            Map<String, Object> metadata = assistantMessage.getMetadata();
            if (ObjectUtils.isNotEmpty(metadata)) {
                Object reasoningContent = metadata.get("reasoningContent");
                if (reasoningContent != null && !reasoningContent.toString().isEmpty()) {
                    builder.reasoningContent(reasoningContent.toString());
                }
            }

            // 提取普通文本内容
            String text = assistantMessage.getText();
            if (text != null && !text.isEmpty()) {
                builder.text(text);
            }
        }

        StreamChunk chunk = builder.build();
        // 如果没有任何内容，标记为空
        if (!chunk.hasText() && !chunk.hasReasoningContent()
                && chunk.getPromptTokens() == null && chunk.getCompletionTokens() == null) {
            return StreamChunk.empty();
        }
        return chunk;
    }

    @Override
    public boolean supports(String providerName) {
        return providerName != null && SUPPORTED_PROVIDERS.contains(providerName.toLowerCase());
    }

    /**
     * 创建 DashScope ChatModel
     */
    private DashScopeChatModel createChatModel(ModelProvider provider, Model model, ChatRequest chatRequest) {
        DashScopeApi dashScopeApi = DashScopeApi.builder()
                .apiKey(provider.getApiKey())
                .build();

        DashScopeChatOptions options = buildOptions(model, chatRequest);

        return DashScopeChatModel.builder()
                .dashScopeApi(dashScopeApi)
                .defaultOptions(options)
                .build();
    }

    /**
     * 构建 DashScope 选项
     */
    private DashScopeChatOptions buildOptions(Model model, ChatRequest chatRequest) {
        DashScopeChatOptions options = new DashScopeChatOptions();
        options.setModel(model.getModelKey());

        if (chatRequest.getTemperature() != null) {
            options.setTemperature(chatRequest.getTemperature());
        }
        if (chatRequest.getMaxTokens() != null) {
            options.setMaxTokens(chatRequest.getMaxTokens());
        }

        // 通义千问深度思考支持
        if (chatRequest.getEnableReasoning() != null && chatRequest.getEnableReasoning()
                && model.getSupportReasoning() != null && model.getSupportReasoning() == 1) {
            options.setEnableThinking(true);
        }

        return options;
    }

    /**
     * 构建 Prompt
     */
    private Prompt buildPrompt(ChatRequest chatRequest) {
        List<Message> messages = chatRequest.getMessages().stream()
                .map(this::convertMessage)
                .collect(Collectors.toList());

        return new Prompt(messages);
    }

    /**
     * 转换消息类型
     */
    private Message convertMessage(ChatMessage msg) {
        String role = msg.getRole();
        if ("system".equals(role)) {
            return new SystemMessage(msg.getContent());
        } else if ("assistant".equals(role)) {
            return new AssistantMessage(msg.getContent());
        } else {
            return new UserMessage(msg.getContent());
        }
    }
}
