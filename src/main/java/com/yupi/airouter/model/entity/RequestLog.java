package com.yupi.airouter.model.entity;

import com.mybatisflex.annotation.Column;
import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.annotation.Table;
import com.mybatisflex.core.keygen.KeyGenerators;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;
import java.time.LocalDateTime;

/**
 * 请求日志 实体类
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("request_log")
public class RequestLog implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * id
     */
    @Id(keyType = KeyType.Generator, value = KeyGenerators.snowFlakeId)
    private Long id;

    /**
     * 链路追踪ID
     */
    @Column("traceId")
    private String traceId;

    /**
     * 用户id
     */
    @Column("userId")
    private Long userId;

    /**
     * API Key id
     */
    @Column("apiKeyId")
    private Long apiKeyId;

    /**
     * 实际调用的模型id
     */
    @Column("modelId")
    private Long modelId;

    /**
     * 请求的模型标识
     */
    @Column("requestModel")
    private String requestModel;

    /**
     * 使用的模型名称（兼容字段）
     */
    @Column("modelName")
    private String modelName;

    /**
     * 请求类型：chat/embedding/image
     */
    @Column("requestType")
    private String requestType;

    /**
     * 调用来源：web/api
     */
    @Column("source")
    private String source;

    /**
     * 输入Token数
     */
    @Column("promptTokens")
    private Integer promptTokens;

    /**
     * 输出Token数
     */
    @Column("completionTokens")
    private Integer completionTokens;

    /**
     * 总Token数
     */
    @Column("totalTokens")
    private Integer totalTokens;

    /**
     * 请求耗时（毫秒）
     */
    @Column("duration")
    private Integer duration;

    /**
     * 状态：success/failed
     */
    @Column("status")
    private String status;

    /**
     * 错误信息
     */
    @Column("errorMessage")
    private String errorMessage;

    /**
     * 错误码
     */
    @Column("errorCode")
    private String errorCode;

    /**
     * 使用的路由策略
     */
    @Column("routingStrategy")
    private String routingStrategy;

    /**
     * 是否为Fallback请求
     */
    @Column("isFallback")
    private Integer isFallback;

    /**
     * 客户端IP
     */
    @Column("clientIp")
    private String clientIp;

    /**
     * User-Agent
     */
    @Column("userAgent")
    private String userAgent;

    /**
     * 创建时间
     */
    @Column("createTime")
    private LocalDateTime createTime;

    /**
     * 更新时间
     */
    @Column("updateTime")
    private LocalDateTime updateTime;
}
