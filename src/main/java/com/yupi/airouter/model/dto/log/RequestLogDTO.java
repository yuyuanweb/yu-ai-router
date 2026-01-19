package com.yupi.airouter.model.dto.log;

import lombok.Builder;
import lombok.Data;

import java.io.Serial;
import java.io.Serializable;

/**
 * 请求日志记录对象
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@Builder
public class RequestLogDTO implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 链路追踪ID
     */
    private String traceId;

    /**
     * 用户ID
     */
    private Long userId;

    /**
     * API Key ID
     */
    private Long apiKeyId;

    /**
     * 实际调用的模型ID
     */
    private Long modelId;

    /**
     * 请求的模型标识
     */
    private String requestModel;

    /**
     * 请求类型：chat/embedding/image
     */
    private String requestType;

    /**
     * 调用来源：web/api
     */
    private String source;

    /**
     * 输入Token数
     */
    private Integer promptTokens;

    /**
     * 输出Token数
     */
    private Integer completionTokens;

    /**
     * 总Token数
     */
    private Integer totalTokens;

    /**
     * 请求耗时（毫秒）
     */
    private Integer duration;

    /**
     * 状态：success/failed
     */
    private String status;

    /**
     * 错误信息
     */
    private String errorMessage;

    /**
     * 错误码
     */
    private String errorCode;

    /**
     * 使用的路由策略
     */
    private String routingStrategy;

    /**
     * 是否为Fallback请求
     */
    private Boolean isFallback;

    /**
     * 客户端IP
     */
    private String clientIp;

    /**
     * User-Agent
     */
    private String userAgent;
}
