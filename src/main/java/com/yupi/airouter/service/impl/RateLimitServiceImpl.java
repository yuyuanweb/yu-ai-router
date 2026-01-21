/**
 * 限流服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.service.impl;

import com.yupi.airouter.service.RateLimitService;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RRateLimiter;
import org.redisson.api.RateType;
import org.redisson.api.RedissonClient;
import org.springframework.stereotype.Service;

import java.time.Duration;

/**
 * 基于 Redisson 的限流服务实现
 * 使用 Redisson 的 RRateLimiter 实现令牌桶限流
 */
@Slf4j
@Service
public class RateLimitServiceImpl implements RateLimitService {

    @Resource
    private RedissonClient redissonClient;

    /**
     * API Key 限流前缀
     */
    private static final String RATE_LIMIT_API_KEY_PREFIX = "rate_limit:api_key:";

    /**
     * IP 限流前缀
     */
    private static final String RATE_LIMIT_IP_PREFIX = "rate_limit:ip:";

    @Override
    public boolean tryAcquire(String key, int limit, Duration duration) {
        // 获取或创建限流器
        RRateLimiter rateLimiter = redissonClient.getRateLimiter(key);

        // 设置限流规则（如果已存在则不会重复设置）
        rateLimiter.trySetRate(RateType.OVERALL, limit, duration);

        // 尝试获取许可
        boolean acquired = rateLimiter.tryAcquire(1);

        if (!acquired) {
            log.debug("Rate limit exceeded for key: {}", key);
        }

        return acquired;
    }

    @Override
    public long getAvailablePermits(String key) {
        RRateLimiter rateLimiter = redissonClient.getRateLimiter(key);
        return rateLimiter.availablePermits();
    }

    @Override
    public boolean checkApiKeyRateLimit(String apiKey, int limit) {
        String key = RATE_LIMIT_API_KEY_PREFIX + apiKey;
        return tryAcquire(key, limit, Duration.ofSeconds(1));
    }

    @Override
    public boolean checkIpRateLimit(String ip, int limit) {
        String key = RATE_LIMIT_IP_PREFIX + ip;
        return tryAcquire(key, limit, Duration.ofSeconds(1));
    }
}
