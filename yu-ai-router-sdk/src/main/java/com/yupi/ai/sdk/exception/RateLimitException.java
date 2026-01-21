/**
 * 限流异常
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.ai.sdk.exception;

/**
 * 限流异常 - 请求过于频繁
 */
public class RateLimitException extends YuAIException {

    public RateLimitException(String message) {
        super(429, message);
    }

    public RateLimitException(String message, Throwable cause) {
        super(429, message, cause);
    }
}
