package com.yupi.airouter.model.vo;

import lombok.Data;

import java.io.Serializable;
import java.math.BigDecimal;
import java.time.LocalDateTime;

/**
 * 模型提供者视图对象
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
public class ProviderVO implements Serializable {

    /**
     * id
     */
    private Long id;

    /**
     * 提供者名称
     */
    private String providerName;

    /**
     * 显示名称
     */
    private String displayName;

    /**
     * API基础URL
     */
    private String baseUrl;

    /**
     * 状态
     */
    private String status;

    /**
     * 健康状态
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
     * 额外配置
     */
    private String config;

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
