package com.yupi.airouter.service;

import com.yupi.airouter.model.dto.chat.ChatRequest;
import com.yupi.airouter.model.dto.chat.ChatResponse;

import java.util.Optional;

/**
 * 响应缓存服务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface CacheService {

    /**
     * 获取缓存的响应
     *
     * @param request 聊天请求
     * @return 缓存的响应，如果不存在则返回空
     */
    Optional<ChatResponse> getCachedResponse(ChatRequest request);

    /**
     * 缓存响应
     *
     * @param request  聊天请求
     * @param response 聊天响应
     */
    void cacheResponse(ChatRequest request, ChatResponse response);

    /**
     * 生成缓存Key
     *
     * @param request 聊天请求
     * @return 缓存Key
     */
    String generateCacheKey(ChatRequest request);

    /**
     * 清除用户的所有缓存
     *
     * @param userId 用户ID
     */
    void clearUserCache(Long userId);

    /**
     * 检查缓存是否启用
     *
     * @return 是否启用缓存
     */
    boolean isCacheEnabled();
}
