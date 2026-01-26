package com.yupi.airouter.controller;

import cn.hutool.core.util.StrUtil;
import cn.hutool.extra.servlet.JakartaServletUtil;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.yupi.airouter.annotation.AuthCheck;
import com.yupi.airouter.annotation.RateLimit;
import com.yupi.airouter.common.ResultUtils;
import com.yupi.airouter.constant.UserConstant;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.model.dto.chat.ChatMessage;
import com.yupi.airouter.model.dto.chat.ChatRequest;
import com.yupi.airouter.model.dto.chat.ChatResponse;
import com.yupi.airouter.model.entity.User;
import com.yupi.airouter.service.ChatService;
import com.yupi.airouter.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import jakarta.annotation.Resource;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.util.List;

/**
 * 内部 Chat 接口（用于前端页面调用）
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@RestController
@RequestMapping("/internal/chat")
@Slf4j
public class InternalChatController {

    @Resource
    private ChatService chatService;

    @Resource
    private UserService userService;

    @Resource
    private ObjectMapper objectMapper;

    /**
     * 内部聊天接口（网页端对话）
     * 支持流式和非流式响应
     */
    @PostMapping(value = "/completions", produces = {MediaType.APPLICATION_JSON_VALUE, MediaType.TEXT_EVENT_STREAM_VALUE})
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "内部聊天接口")
    @RateLimit(type = RateLimit.LimitType.IP, limit = 30)
    public Object chatCompletions(@RequestBody ChatRequest request,
                                  HttpServletRequest httpRequest) {
        User loginUser = userService.getLoginUser(httpRequest);

        // 1. 参数校验
        if (request.getMessages() == null || request.getMessages().isEmpty()) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "messages 不能为空");
        }

        // 2. 获取客户端IP和User-Agent
        String clientIp = JakartaServletUtil.getClientIP(httpRequest);
        String userAgent = httpRequest.getHeader("User-Agent");

        // 3. 判断是否为流式请求
        Boolean stream = request.getStream();
        if (stream != null && stream) {
            // 流式响应（网页端调用，不传 apiKeyId）- 转换为 SSE 格式
            return chatService.chatStream(request, loginUser.getId(), null, clientIp, userAgent)
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
            // 非流式响应（网页端调用，不传 apiKeyId）
            ChatResponse response = chatService.chat(request, loginUser.getId(), null, clientIp, userAgent);
            return ResultUtils.success(response);
        }
    }

    /**
     * 内部聊天接口 - 支持文件上传（图片识别、PDF解析）
     * 使用 multipart/form-data 格式
     */
    @PostMapping(value = "/completions/upload", consumes = MediaType.MULTIPART_FORM_DATA_VALUE,
            produces = {MediaType.APPLICATION_JSON_VALUE, MediaType.TEXT_EVENT_STREAM_VALUE})
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "内部聊天接口（支持文件上传）")
    @RateLimit(type = RateLimit.LimitType.IP, limit = 30)
    public Object chatCompletionsWithFile(
            @RequestPart("file") MultipartFile file,
            @RequestParam(value = "model", required = false) String model,
            @RequestParam(value = "messages", required = true) String messagesJson,
            @RequestParam(value = "stream", required = false, defaultValue = "false") Boolean stream,
            @RequestParam(value = "routing_strategy", required = false) String routingStrategy,
            @RequestParam(value = "plugin_key", required = false) String pluginKey,
            @RequestParam(value = "enable_reasoning", required = false) Boolean enableReasoning,
            HttpServletRequest httpRequest) {

        User loginUser = userService.getLoginUser(httpRequest);

        try {
            // 1. 解析消息列表
            List<ChatMessage> messages = objectMapper.readValue(messagesJson, new TypeReference<List<ChatMessage>>() {});
            
            if (messages == null || messages.isEmpty()) {
                throw new BusinessException(ErrorCode.PARAMS_ERROR, "messages 不能为空");
            }

            // 2. 构建 ChatRequest
            ChatRequest request = new ChatRequest();
            request.setMessages(messages);
            request.setModel(model);
            request.setStream(stream);
            request.setRoutingStrategy(routingStrategy);
            request.setPluginKey(pluginKey);
            request.setEnableReasoning(enableReasoning);

            // 3. 处理上传的文件
            if (file != null && !file.isEmpty()) {
                // 直接使用字节数组
                request.setFileBytes(file.getBytes());

                // 设置文件类型
                String contentType = file.getContentType();
                request.setFileType(contentType);

                // 根据文件类型自动设置插件（如果未指定）
                if (StrUtil.isBlank(request.getPluginKey())) {
                    if (contentType != null) {
                        if (contentType.startsWith("image/")) {
                            request.setPluginKey("image_recognition");
                        } else if (contentType.equals("application/pdf")) {
                            request.setPluginKey("pdf_parser");
                        }
                    }
                }

                log.info("文件上传成功: name={}, size={}, type={}, plugin={}",
                        file.getOriginalFilename(), file.getSize(), contentType, request.getPluginKey());
            }

            // 4. 设置默认模型（如果未指定）
            if (request.getModel() == null || request.getModel().isEmpty()) {
                request.setModel("qwen-plus");
            }

            // 5. 获取客户端IP和User-Agent
            String clientIp = JakartaServletUtil.getClientIP(httpRequest);
            String userAgent = httpRequest.getHeader("User-Agent");

            // 6. 判断是否为流式请求
            if (stream != null && stream) {
                return chatService.chatStream(request, loginUser.getId(), null, clientIp, userAgent)
                        .map(streamResponse -> {
                            try {
                                String json = objectMapper.writeValueAsString(streamResponse);
                                return json + "\n\n";
                            } catch (Exception e) {
                                log.error("Failed to serialize stream response", e);
                                return "";
                            }
                        });
            } else {
                ChatResponse response = chatService.chat(request, loginUser.getId(), null, clientIp, userAgent);
                return ResultUtils.success(response);
            }
        } catch (IOException e) {
            log.error("处理请求失败", e);
            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "处理请求失败: " + e.getMessage());
        }
    }
}
