package com.yupi.airouter.service.impl;

import cn.hutool.core.util.StrUtil;
import cn.hutool.crypto.digest.DigestUtil;
import cn.hutool.json.JSONUtil;
import com.yupi.airouter.model.dto.chat.ChatRequest;
import com.yupi.airouter.model.dto.chat.ChatResponse;
import com.yupi.airouter.service.CacheService;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RBucket;
import org.redisson.api.RedissonClient;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.time.Duration;
import java.util.Optional;

/**
 * 响应缓存服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Slf4j
@Service
public class CacheServiceImpl implements CacheService {

    @Resource
    private RedissonClient redissonClient;

    /**
     * 缓存前缀
     */
    private static final String CACHE_PREFIX = "ai:response:";

    /**
     * 是否启用缓存
     */
    @Value("${ai.cache.enabled:true}")
    private boolean cacheEnabled;

    /**
     * 缓存过期时间（秒）
     */
    @Value("${ai.cache.ttl:3600}")
    private long cacheTtl;

    @Override
    public Optional<ChatResponse> getCachedResponse(ChatRequest request) {
        if (!isCacheEnabled()) {
            return Optional.empty();
        }

        try {
            String cacheKey = generateCacheKey(request);
            RBucket<String> bucket = redissonClient.getBucket(cacheKey);
            String cachedJson = bucket.get();

            if (StrUtil.isNotBlank(cachedJson)) {
                ChatResponse response = JSONUtil.toBean(cachedJson, ChatResponse.class);
                log.debug("命中缓存: {}", cacheKey);
                return Optional.of(response);
            }
        } catch (Exception e) {
            log.warn("获取缓存失败", e);
        }

        return Optional.empty();
    }

    @Override
    public void cacheResponse(ChatRequest request, ChatResponse response) {
        if (!isCacheEnabled()) {
            return;
        }

        try {
            String cacheKey = generateCacheKey(request);
            String responseJson = JSONUtil.toJsonStr(response);

            RBucket<String> bucket = redissonClient.getBucket(cacheKey);
            bucket.set(responseJson, Duration.ofSeconds(cacheTtl));
            log.debug("写入缓存: {}", cacheKey);
        } catch (Exception e) {
            log.warn("写入缓存失败", e);
        }
    }

    @Override
    public String generateCacheKey(ChatRequest request) {
        // 基于模型和消息内容生成缓存Key
        StringBuilder keyBuilder = new StringBuilder();
        keyBuilder.append(request.getModel() != null ? request.getModel() : "auto");
        keyBuilder.append(":");

        // 将消息内容序列化并计算MD5
        if (request.getMessages() != null) {
            String messagesJson = JSONUtil.toJsonStr(request.getMessages());
            String messagesMd5 = DigestUtil.md5Hex(messagesJson);
            keyBuilder.append(messagesMd5);
        }

        // 添加温度参数（如果有）
        if (request.getTemperature() != null) {
            keyBuilder.append(":t").append(request.getTemperature());
        }

        return CACHE_PREFIX + keyBuilder;
    }

    @Override
    public void clearUserCache(Long userId) {
        // 由于缓存是基于请求内容的，不是基于用户的
        // 这里可以实现清除特定用户的缓存逻辑（如果需要的话）
        log.info("清除用户 {} 的缓存（当前实现为空操作）", userId);
    }

    @Override
    public boolean isCacheEnabled() {
        return cacheEnabled;
    }
}
