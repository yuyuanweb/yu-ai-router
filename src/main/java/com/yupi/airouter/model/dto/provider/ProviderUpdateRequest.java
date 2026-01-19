package com.yupi.airouter.model.dto.provider;

import lombok.Data;

import java.io.Serial;
import java.io.Serializable;

/**
 * 模型提供者更新请求
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
public class ProviderUpdateRequest implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * id
     */
    private Long id;

    /**
     * 显示名称
     */
    private String displayName;

    /**
     * API基础URL
     */
    private String baseUrl;

    /**
     * API密钥
     */
    private String apiKey;

    /**
     * 状态：active/inactive/maintenance
     */
    private String status;

    /**
     * 优先级（越大越优先）
     */
    private Integer priority;

    /**
     * 额外配置（JSON格式）
     */
    private String config;
}
