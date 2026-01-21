/**
 * 限流注解
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.annotation;

import java.lang.annotation.*;
import java.util.concurrent.TimeUnit;

/**
 * 限流注解
 * 支持 API Key 级别和 IP 级别的限流
 */
@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
@Documented
public @interface RateLimit {

    /**
     * 限流类型
     */
    LimitType type() default LimitType.API_KEY;

    /**
     * 限流数量（默认每秒 10 次）
     */
    int limit() default 10;

    /**
     * 时间窗口，默认 1
     */
    int window() default 1;

    /**
     * 时间单位，默认秒
     */
    TimeUnit timeUnit() default TimeUnit.SECONDS;

    /**
     * 限流类型枚举
     */
    enum LimitType {
        /**
         * 基于 API Key 限流
         */
        API_KEY,
        /**
         * 基于 IP 限流
         */
        IP
    }
}
