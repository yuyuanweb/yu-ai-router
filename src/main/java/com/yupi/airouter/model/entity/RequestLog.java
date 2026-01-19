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
     * 使用的模型名称
     */
    @Column("modelName")
    private String modelName;

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
