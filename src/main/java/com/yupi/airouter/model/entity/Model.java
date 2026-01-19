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
 * 模型 实体类
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("model")
public class Model implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * id
     */
    @Id(keyType = KeyType.Generator, value = KeyGenerators.snowFlakeId)
    private Long id;

    /**
     * 提供者id
     */
    @Column("providerId")
    private Long providerId;

    /**
     * 模型标识（如：qwen-plus）
     */
    @Column("modelKey")
    private String modelKey;

    /**
     * 模型显示名称
     */
    @Column("modelName")
    private String modelName;

    /**
     * 模型类型：chat/embedding/image/audio
     */
    @Column("modelType")
    private String modelType;

    /**
     * 模型描述
     */
    @Column("description")
    private String description;

    /**
     * 上下文长度限制
     */
    @Column("contextLength")
    private Integer contextLength;

    /**
     * 输入价格（元/千Token）
     */
    @Column("inputPrice")
    private BigDecimal inputPrice;

    /**
     * 输出价格（元/千Token）
     */
    @Column("outputPrice")
    private BigDecimal outputPrice;

    /**
     * 状态：active/inactive/deprecated
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
    private java.math.BigDecimal successRate;

    /**
     * 综合得分（越低越好）
     */
    @Column("score")
    private java.math.BigDecimal score;

    /**
     * 优先级（越大越优先）
     */
    @Column("priority")
    private Integer priority;

    /**
     * 默认超时时间（毫秒）
     */
    @Column("defaultTimeout")
    private Integer defaultTimeout;

    /**
     * 是否支持深度思考：0=不支持，1=支持
     */
    @Column("supportReasoning")
    private Integer supportReasoning;

    /**
     * 能力标签（JSON数组）
     */
    @Column("capabilities")
    private String capabilities;

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
