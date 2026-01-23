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
 * 图片生成记录 实体类
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("image_generation_record")
public class ImageGenerationRecord implements Serializable {

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
     * 使用的模型id
     */
    @Column("modelId")
    private Long modelId;

    /**
     * 模型标识
     */
    @Column("modelKey")
    private String modelKey;

    /**
     * 生成提示词
     */
    @Column("prompt")
    private String prompt;

    /**
     * 修订后的提示词
     */
    @Column("revisedPrompt")
    private String revisedPrompt;

    /**
     * 图片URL
     */
    @Column("imageUrl")
    private String imageUrl;

    /**
     * Base64图片数据
     */
    @Column("imageData")
    private String imageData;

    /**
     * 图片尺寸
     */
    @Column("size")
    private String size;

    /**
     * 图片质量
     */
    @Column("quality")
    private String quality;

    /**
     * 状态：success/failed
     */
    @Column("status")
    private String status;

    /**
     * 生成费用（元）
     */
    @Column("cost")
    private BigDecimal cost;

    /**
     * 耗时（毫秒）
     */
    @Column("duration")
    private Integer duration;

    /**
     * 错误信息
     */
    @Column("errorMessage")
    private String errorMessage;

    /**
     * 客户端IP
     */
    @Column("clientIp")
    private String clientIp;

    /**
     * 创建时间
     */
    @Column("createTime")
    private LocalDateTime createTime;

}
