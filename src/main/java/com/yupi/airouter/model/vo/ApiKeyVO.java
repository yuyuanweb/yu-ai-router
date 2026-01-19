package com.yupi.airouter.model.vo;

import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.time.LocalDateTime;

/**
 * API Key 视图对象
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
public class ApiKeyVO implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * id
     */
    private Long id;

    /**
     * API Key值（创建时完整显示，列表中部分隐藏）
     */
    private String keyValue;

    /**
     * Key名称/备注
     */
    private String keyName;

    /**
     * 状态：active/inactive/revoked
     */
    private String status;

    /**
     * 已使用Token总数
     */
    private Long totalTokens;

    /**
     * 最后使用时间
     */
    private LocalDateTime lastUsedTime;

    /**
     * 创建时间
     */
    private LocalDateTime createTime;
}
