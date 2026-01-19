package com.yupi.airouter.model.dto.model;

import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.math.BigDecimal;

/**
 * 模型更新请求
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
public class ModelUpdateRequest implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * id
     */
    private Long id;

    /**
     * 模型显示名称
     */
    private String modelName;

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
     * 状态：active/inactive/deprecated
     */
    private String status;

    /**
     * 优先级（越大越优先）
     */
    private Integer priority;

    /**
     * 默认超时时间（毫秒）
     */
    private Integer defaultTimeout;

    /**
     * 能力标签（JSON数组）
     */
    private String capabilities;
}
