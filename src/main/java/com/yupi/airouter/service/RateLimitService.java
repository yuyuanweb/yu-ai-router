/**
 * 限流服务接口
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.service;

import java.time.Duration;

/**
 * 限流服务
 * 提供基于 Redisson 的限流能力
 */
public interface RateLimitService {

    /**
     * 尝试获取限流许可
     *
     * @param key      限流键（如 API Key 或 IP）
     * @param limit    限制次数
     * @param duration 时间窗口
     * @return true 表示允许访问，false 表示被限流
     */
    boolean tryAcquire(String key, int limit, Duration duration);

    /**
     * 获取当前可用许可数
     *
     * @param key 限流键
     * @return 当前可用的许可数
     */
    long getAvailablePermits(String key);

    /**
     * 检查 API Key 是否被限流
     *
     * @param apiKey API Key 值
     * @param limit  限制次数（每秒）
     * @return true 表示允许访问，false 表示被限流
     */
    boolean checkApiKeyRateLimit(String apiKey, int limit);

    /**
     * 检查 IP 是否被限流
     *
     * @param ip    客户端 IP
     * @param limit 限制次数（每秒）
     * @return true 表示允许访问，false 表示被限流
     */
    boolean checkIpRateLimit(String ip, int limit);
}
