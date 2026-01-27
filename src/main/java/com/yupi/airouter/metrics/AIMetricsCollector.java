/**
 * AI 网关业务指标收集器
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.metrics;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

import java.util.concurrent.ConcurrentHashMap;

@Slf4j
@Component
public class AIMetricsCollector {

    private final MeterRegistry meterRegistry;

    private Counter totalRequestsCounter;
    private Counter totalTokensCounter;
    private Counter totalErrorsCounter;

    // 按模型和用户维度的计数器缓存
    private final ConcurrentHashMap<String, Counter> modelRequestCounters = new ConcurrentHashMap<>();
    private final ConcurrentHashMap<String, Counter> modelTokenCounters = new ConcurrentHashMap<>();
    private final ConcurrentHashMap<String, Counter> userRequestCounters = new ConcurrentHashMap<>();
    private final ConcurrentHashMap<String, Timer> modelTimers = new ConcurrentHashMap<>();

    public AIMetricsCollector(MeterRegistry meterRegistry) {
        this.meterRegistry = meterRegistry;
    }

    @PostConstruct
    public void init() {
        // 初始化全局计数器
        totalRequestsCounter = Counter.builder("ai.requests.total")
                .description("AI 请求总数")
                .register(meterRegistry);

        totalTokensCounter = Counter.builder("ai.tokens.total")
                .description("Token 消耗总量")
                .register(meterRegistry);

        totalErrorsCounter = Counter.builder("ai.errors.total")
                .description("错误总数")
                .register(meterRegistry);

        log.info("AI 指标收集器初始化完成");
    }

    /**
     * 记录请求
     */
    public void recordRequest(String modelKey, Long userId, String apiKeyId) {
        // 全局请求计数
        totalRequestsCounter.increment();

        // 按模型计数
        if (modelKey != null) {
            getOrCreateModelRequestCounter(modelKey).increment();
        }

        // 按用户计数
        if (userId != null) {
            getOrCreateUserRequestCounter(userId.toString()).increment();
        }
    }

    /**
     * 记录 Token 消耗
     */
    public void recordTokens(String modelKey, int tokens) {
        // 全局 Token 计数
        totalTokensCounter.increment(tokens);

        // 按模型计数
        if (modelKey != null) {
            getOrCreateModelTokenCounter(modelKey).increment(tokens);
        }
    }

    /**
     * 记录错误
     */
    public void recordError(String modelKey, String errorType) {
        totalErrorsCounter.increment();

        Counter.builder("ai.errors")
                .tag("model", modelKey != null ? modelKey : "unknown")
                .tag("type", errorType != null ? errorType : "unknown")
                .register(meterRegistry)
                .increment();
    }

    /**
     * 记录响应时间
     */
    public void recordResponseTime(String modelKey, long durationMillis) {
        if (modelKey != null) {
            Timer timer = modelTimers.computeIfAbsent(modelKey, k ->
                    Timer.builder("ai.response.time")
                            .description("AI 响应时间")
                            .tag("model", k)
                            .register(meterRegistry)
            );
            timer.record(durationMillis, java.util.concurrent.TimeUnit.MILLISECONDS);
        }
    }

    /**
     * 获取或创建模型请求计数器
     */
    private Counter getOrCreateModelRequestCounter(String modelKey) {
        return modelRequestCounters.computeIfAbsent(modelKey, k ->
                Counter.builder("ai.requests.by_model")
                        .description("按模型统计的请求数")
                        .tag("model", k)
                        .register(meterRegistry)
        );
    }

    /**
     * 获取或创建模型 Token 计数器
     */
    private Counter getOrCreateModelTokenCounter(String modelKey) {
        return modelTokenCounters.computeIfAbsent(modelKey, k ->
                Counter.builder("ai.tokens.by_model")
                        .description("按模型统计的 Token 消耗")
                        .tag("model", k)
                        .register(meterRegistry)
        );
    }

    /**
     * 获取或创建用户请求计数器
     */
    private Counter getOrCreateUserRequestCounter(String userId) {
        return userRequestCounters.computeIfAbsent(userId, k ->
                Counter.builder("ai.requests.by_user")
                        .description("按用户统计的请求数")
                        .tag("user_id", k)
                        .register(meterRegistry)
        );
    }
}
