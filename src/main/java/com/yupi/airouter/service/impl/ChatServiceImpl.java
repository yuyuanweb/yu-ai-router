package com.yupi.airouter.service.impl;

import cn.hutool.core.util.IdUtil;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.model.dto.chat.ChatMessage;
import com.yupi.airouter.model.dto.chat.ChatRequest;
import com.yupi.airouter.model.dto.chat.ChatResponse;
import com.yupi.airouter.service.ChatService;
import com.yupi.airouter.service.RequestLogService;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.ObjectUtils;
import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.chat.messages.AssistantMessage;
import org.springframework.ai.chat.messages.Message;
import org.springframework.ai.chat.messages.SystemMessage;
import org.springframework.ai.chat.messages.UserMessage;
import org.springframework.ai.chat.model.ChatModel;
import org.springframework.ai.chat.prompt.Prompt;
import org.springframework.ai.openai.OpenAiChatOptions;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Flux;

import java.util.List;
import java.util.stream.Collectors;

/**
 * 聊天服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Slf4j
@Service
public class ChatServiceImpl implements ChatService {

    @Resource
    private ChatModel chatModel;

    @Resource
    private RequestLogService requestLogService;

    @Override
    public ChatResponse chat(ChatRequest chatRequest, Long userId, Long apiKeyId) {
        long startTime = System.currentTimeMillis();
        String modelName = chatRequest.getModel();

        try {
            // 构建 Prompt
            Prompt prompt = buildPrompt(chatRequest);

            // 调用 AI 模型
            ChatClient chatClient = ChatClient.builder(chatModel).build();
            org.springframework.ai.chat.model.ChatResponse aiResponse = chatClient
                    .prompt(prompt)
                    .call()
                    .chatResponse();

            // 转换响应格式
            ChatResponse response = convertResponse(aiResponse, modelName);

            // 记录请求日志
            long duration = System.currentTimeMillis() - startTime;
            requestLogService.logRequest(userId, apiKeyId, modelName,
                    response.getUsage().getPromptTokens(),
                    response.getUsage().getCompletionTokens(),
                    response.getUsage().getTotalTokens(),
                    (int) duration, "success", null);

            return response;
        } catch (Exception e) {
            log.error("调用模型失败", e);
            long duration = System.currentTimeMillis() - startTime;
            requestLogService.logRequest(userId, apiKeyId, modelName, 0, 0, 0,
                    (int) duration, "failed", e.getMessage());
            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "调用模型失败: " + e.getMessage());
        }
    }

    @Override
    public Flux<String> chatStream(ChatRequest chatRequest, Long userId, Long apiKeyId) {
        String modelName = chatRequest.getModel();
        long startTime = System.currentTimeMillis();

        // Token 计数器和上次内容
        final int[] promptTokens = {0};
        final int[] completionTokens = {0};

        try {
            // 构建 Prompt
            Prompt prompt = buildPrompt(chatRequest);

            // 调用流式 AI 模型
            ChatClient chatClient = ChatClient.builder(chatModel).build();
            Flux<org.springframework.ai.chat.model.ChatResponse> flux = chatClient
                    .prompt(prompt)
                    .stream()
                    .chatResponse();

            return flux.flatMap(response -> {
                // 获取 Token 信息（通常只有最后一个 chunk 才有）
                if (response.getMetadata().getUsage() != null) {
                    Integer promptToken = response.getMetadata().getUsage().getPromptTokens();
                    Integer completion = response.getMetadata().getUsage().getCompletionTokens();

                    if (promptToken != null && promptToken > 0) {
                        promptTokens[0] = promptToken;
                    }
                    if (completion != null && completion > 0) {
                        completionTokens[0] = completion;
                    }
                }

                // 检查 result 是否为 null（最后一个 chunk 只有 token 信息，没有内容）
                if (ObjectUtils.isEmpty(response.getResult()) || ObjectUtils.isEmpty(response.getResult().getOutput())) {
                    // 跳过没有内容的 chunk
                    return Flux.empty();
                }
                // 将文本中的换行符转义为 \\n，避免与 SSE 格式冲突
                String text = response.getResult().getOutput().getText();
                String escapedText = text.replace("\n", "\\n");
                return Flux.just(escapedText + "\n\n");
            }).doOnComplete(() -> {
                // 流结束时记录日志
                long duration = System.currentTimeMillis() - startTime;
                int totalTokens = promptTokens[0] + completionTokens[0];

                requestLogService.logRequest(userId, apiKeyId, modelName,
                        promptTokens[0], completionTokens[0], totalTokens,
                        (int) duration, "success", null);
            }).doOnError(error -> {
                // 流错误时记录日志
                log.error("流式调用模型失败", error);
                long duration = System.currentTimeMillis() - startTime;
                requestLogService.logRequest(userId, apiKeyId, modelName, 0, 0, 0,
                        (int) duration, "failed", error.getMessage());
            });
        } catch (Exception e) {
            log.error("流式调用模型失败", e);
            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "流式调用模型失败: " + e.getMessage());
        }
    }

    /**
     * 构建 Prompt
     */
    private Prompt buildPrompt(ChatRequest chatRequest) {
        List<Message> messages = chatRequest.getMessages().stream()
                .map(msg -> switch (msg.getRole()) {
                    case "system" -> new SystemMessage(msg.getContent());
                    case "assistant" -> new AssistantMessage(msg.getContent());
                    default -> new UserMessage(msg.getContent());
                })
                .collect(Collectors.toList());

        // 构建选项
        OpenAiChatOptions.Builder optionsBuilder = OpenAiChatOptions.builder()
                .model(chatRequest.getModel())
                // 启用流式响应的 token 统计
                .streamUsage(true);

        if (chatRequest.getTemperature() != null) {
            optionsBuilder.temperature(chatRequest.getTemperature());
        }
        if (chatRequest.getMaxTokens() != null) {
            optionsBuilder.maxTokens(chatRequest.getMaxTokens());
        }

        return new Prompt(messages, optionsBuilder.build());
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
}
