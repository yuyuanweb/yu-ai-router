/**
 * 限流切面
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.aop;

import cn.hutool.extra.servlet.JakartaServletUtil;
import com.yupi.airouter.annotation.RateLimit;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.service.RateLimitService;
import jakarta.annotation.Resource;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.springframework.stereotype.Component;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;

import java.time.Duration;
import java.util.concurrent.TimeUnit;

/**
 * 限流切面
 * 用于拦截带有 @RateLimit 注解的方法，实现限流控制
 */
@Slf4j
@Aspect
@Component
public class RateLimitAspect {

    @Resource
    private RateLimitService rateLimitService;

    @Around("@annotation(rateLimit)")
    public Object doRateLimit(ProceedingJoinPoint joinPoint, RateLimit rateLimit) throws Throwable {
        // 获取请求对象
        ServletRequestAttributes attributes = (ServletRequestAttributes) RequestContextHolder.getRequestAttributes();
        if (attributes == null) {
            return joinPoint.proceed();
        }
        HttpServletRequest request = attributes.getRequest();

        // 根据限流类型获取限流 key
        String limitKey = getLimitKey(rateLimit, request);
        if (limitKey == null) {
            return joinPoint.proceed();
        }

        // 将 TimeUnit 转换为 Duration
        Duration duration = toDuration(rateLimit.window(), rateLimit.timeUnit());

        // 执行限流检查
        boolean allowed = rateLimitService.tryAcquire(limitKey, rateLimit.limit(), duration);

        if (!allowed) {
            log.warn("Rate limit exceeded, key: {}, limit: {}/{} {}",
                    limitKey, rateLimit.limit(), rateLimit.window(), rateLimit.timeUnit());
            throw new BusinessException(ErrorCode.TOO_MANY_REQUEST, "请求过于频繁，请稍后再试");
        }

        return joinPoint.proceed();
    }

    /**
     * 根据限流类型获取限流 key
     */
    private String getLimitKey(RateLimit rateLimit, HttpServletRequest request) {
        switch (rateLimit.type()) {
            case API_KEY:
                String authorization = request.getHeader("Authorization");
                if (authorization != null && authorization.startsWith("Bearer ")) {
                    return "rate_limit:api_key:" + authorization.substring(7);
                }
                return null;
            case IP:
                return "rate_limit:ip:" + JakartaServletUtil.getClientIP(request);
            default:
                return null;
        }
    }

    /**
     * 将 TimeUnit 转换为 Duration
     */
    private Duration toDuration(long amount, TimeUnit timeUnit) {
        return switch (timeUnit) {
            case NANOSECONDS -> Duration.ofNanos(amount);
            case MICROSECONDS -> Duration.ofNanos(amount * 1000);
            case MILLISECONDS -> Duration.ofMillis(amount);
            case SECONDS -> Duration.ofSeconds(amount);
            case MINUTES -> Duration.ofMinutes(amount);
            case HOURS -> Duration.ofHours(amount);
            case DAYS -> Duration.ofDays(amount);
        };
    }
}
