package com.yupi.airouter.model.vo;

import lombok.Data;

import java.io.Serializable;
import java.math.BigDecimal;
import java.time.LocalDateTime;

/**
 * 模型视图对象
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
public class ModelVO implements Serializable {

    /**
     * id
     */
    private Long id;

    /**
     * 提供者id
     */
    private Long providerId;

    /**
     * 提供者名称
     */
    private String providerName;

    /**
     * 提供者显示名称
     */
    private String providerDisplayName;

    /**
     * 模型标识
     */
    private String modelKey;

    /**
     * 模型显示名称
     */
    private String modelName;

    /**
     * 模型类型
     */
    private String modelType;

    /**
     * 模型描述
     */
    private String description;

    /**
     * 上下文长度限制
     */
    private Integer contextLength;

    /**
     * 输入价格（元/千Token）
     */
    private BigDecimal inputPrice;

    /**
     * 输出价格（元/千Token）
     */
    private BigDecimal outputPrice;

    /**
     * 状态
     */
    private String status;

    /**
     * 健康状态：healthy/unhealthy/degraded/unknown
     */
    private String healthStatus;

    /**
     * 平均延迟（毫秒）
     */
    private Integer avgLatency;

    /**
     * 成功率（百分比）
     */
    private BigDecimal successRate;

    /**
     * 优先级
     */
    private Integer priority;

    /**
     * 默认超时时间（毫秒）
     */
    private Integer defaultTimeout;

    /**
     * 是否支持深度思考：0=不支持，1=支持
     */
    private Integer supportReasoning;

    /**
     * 能力标签
     */
    private String capabilities;

    /**
     * 创建时间
     */
    private LocalDateTime createTime;

    /**
     * 更新时间
     */
    private LocalDateTime updateTime;

    private static final long serialVersionUID = 1L;
}
