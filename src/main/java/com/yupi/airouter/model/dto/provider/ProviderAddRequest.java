package com.yupi.airouter.model.dto.provider;

import lombok.Data;

import java.io.Serial;
import java.io.Serializable;

/**
 * 模型提供者创建请求
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
public class ProviderAddRequest implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 提供者名称（如：qwen/zhipu/deepseek）
     */
    private String providerName;

    /**
     * 显示名称（如：通义千问/智谱AI/DeepSeek）
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
     * 优先级（越大越优先）
     */
    private Integer priority;

    /**
     * 额外配置（JSON格式）
     */
    private String config;
}
