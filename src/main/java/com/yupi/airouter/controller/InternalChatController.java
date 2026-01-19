package com.yupi.airouter.controller;

import com.yupi.airouter.annotation.AuthCheck;
import com.yupi.airouter.common.BaseResponse;
import com.yupi.airouter.common.ResultUtils;
import com.yupi.airouter.constant.UserConstant;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.model.dto.chat.ChatRequest;
import com.yupi.airouter.model.dto.chat.ChatResponse;
import com.yupi.airouter.model.entity.ApiKey;
import com.yupi.airouter.model.entity.User;
import com.yupi.airouter.service.ApiKeyService;
import com.yupi.airouter.service.ChatService;
import com.yupi.airouter.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Flux;

import jakarta.annotation.Resource;
import jakarta.servlet.http.HttpServletRequest;

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
    private ApiKeyService apiKeyService;

    @Resource
    private UserService userService;

    /**
     * 内部聊天接口（使用 API Key ID）
     * 支持流式和非流式响应
     */
    @PostMapping(value = "/completions", produces = {MediaType.APPLICATION_JSON_VALUE, MediaType.TEXT_EVENT_STREAM_VALUE})
    @AuthCheck(mustRole = UserConstant.DEFAULT_ROLE)
    @Operation(summary = "内部聊天接口")
    public Object chatCompletions(@RequestBody ChatRequest request,
                                   @RequestParam Long apiKeyId,
                                   HttpServletRequest httpRequest) {
        User loginUser = userService.getLoginUser(httpRequest);

        // 1. 验证 API Key 归属
        ApiKey apiKey = apiKeyService.getById(apiKeyId);
        if (apiKey == null) {
            throw new BusinessException(ErrorCode.NOT_FOUND_ERROR, "API Key 不存在");
        }

        if (!apiKey.getUserId().equals(loginUser.getId())) {
            throw new BusinessException(ErrorCode.NO_AUTH_ERROR, "无权使用该 API Key");
        }

        if (!"active".equals(apiKey.getStatus())) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "API Key 已失效");
        }

        // 2. 参数校验
        if (request.getMessages() == null || request.getMessages().isEmpty()) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "messages 不能为空");
        }

        // 3. 设置默认模型（如果未指定）
        if (request.getModel() == null || request.getModel().isEmpty()) {
            request.setModel("qwen-plus");
        }

        // 4. 判断是否为流式请求
        Boolean stream = request.getStream();
        if (stream != null && stream) {
            // 流式响应
            return chatService.chatStream(request, loginUser.getId(), apiKey.getId());
        } else {
            // 非流式响应
            ChatResponse response = chatService.chat(request, loginUser.getId(), apiKey.getId());
            return ResultUtils.success(response);
        }
    }
}
