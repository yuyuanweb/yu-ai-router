package com.yupi.airouter.controller;

import cn.hutool.extra.servlet.JakartaServletUtil;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.yupi.airouter.annotation.RateLimit;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.model.dto.chat.ChatRequest;
import com.yupi.airouter.model.entity.ApiKey;
import com.yupi.airouter.service.ApiKeyService;
import com.yupi.airouter.service.ChatService;
import io.swagger.v3.oas.annotations.Operation;
import jakarta.annotation.Resource;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;

/**
 * Chat Completions 接口（兼容 OpenAI 格式）
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@RestController
@RequestMapping("/v1/chat")
@Slf4j
public class ChatController {

    @Resource
    private ChatService chatService;

    @Resource
    private ApiKeyService apiKeyService;

    @Resource
    private ObjectMapper objectMapper;

    /**
     * Chat Completions 接口
     * 支持流式和非流式响应
     */
    @PostMapping(value = "/completions", produces = {MediaType.APPLICATION_JSON_VALUE, MediaType.TEXT_EVENT_STREAM_VALUE})
    @Operation(summary = "Chat Completions")
    @RateLimit(type = RateLimit.LimitType.API_KEY, limit = 60)
    public Object chatCompletions(@RequestBody ChatRequest request,
                                  @RequestHeader(value = "Authorization", required = false) String authorization,
                                  HttpServletRequest httpRequest) {
        // 1. 验证 API Key
        if (authorization == null || !authorization.startsWith("Bearer ")) {
            throw new BusinessException(ErrorCode.NO_AUTH_ERROR, "缺少或无效的 Authorization Header");
        }

        // 去掉 "Bearer " 前缀
        String apiKeyValue = authorization.substring(7);
        ApiKey apiKey = apiKeyService.getByKeyValue(apiKeyValue);

        if (apiKey == null) {
            throw new BusinessException(ErrorCode.NO_AUTH_ERROR, "API Key 无效或已失效");
        }

        // 2. 参数校验
        if (request.getMessages() == null || request.getMessages().isEmpty()) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "messages 不能为空");
        }

        // 3. 获取客户端IP和User-Agent
        String clientIp = JakartaServletUtil.getClientIP(httpRequest);
        String userAgent = httpRequest.getHeader("User-Agent");

        // 4. 判断是否为流式请求
        Boolean stream = request.getStream();
        if (stream != null && stream) {
            // 流式响应 - 转换为 SSE 格式
            return chatService.chatStream(request, apiKey.getUserId(), apiKey.getId(), clientIp, userAgent)
                    .map(streamResponse -> {
                        try {
                            // 将 StreamResponse 转换为 JSON
                            String json = objectMapper.writeValueAsString(streamResponse);
                            return json + "\n\n";
                        } catch (Exception e) {
                            log.error("Failed to serialize stream response", e);
                            return "";
                        }
                    });
        } else {
            // 非流式响应
            return chatService.chat(request, apiKey.getUserId(), apiKey.getId(), clientIp, userAgent);
        }
    }
}
