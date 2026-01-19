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
import java.math.BigDecimal;
import java.time.LocalDateTime;

/**
 * 模型提供者 实体类
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("model_provider")
public class ModelProvider implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * id
     */
    @Id(keyType = KeyType.Generator, value = KeyGenerators.snowFlakeId)
    private Long id;

    /**
     * 提供者名称（如：qwen/zhipu/deepseek）
     */
    @Column("providerName")
    private String providerName;

    /**
     * 显示名称（如：通义千问/智谱AI/DeepSeek）
     */
    @Column("displayName")
    private String displayName;

    /**
     * API基础URL
     */
    @Column("baseUrl")
    private String baseUrl;

    /**
     * API密钥
     */
    @Column("apiKey")
    private String apiKey;

    /**
     * 状态：active/inactive/maintenance
     */
    @Column("status")
    private String status;

    /**
     * 健康状态：healthy/unhealthy/degraded/unknown
     */
    @Column("healthStatus")
    private String healthStatus;

    /**
     * 平均延迟（毫秒）
     */
    @Column("avgLatency")
    private Integer avgLatency;

    /**
     * 成功率（百分比）
     */
    @Column("successRate")
    private BigDecimal successRate;

    /**
     * 优先级（越大越优先）
     */
    @Column("priority")
    private Integer priority;

    /**
     * 额外配置（JSON格式）
     */
    @Column("config")
    private String config;

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

    /**
     * 是否删除
     */
    @Column(value = "isDelete", isLogicDelete = true)
    private Integer isDelete;
}
