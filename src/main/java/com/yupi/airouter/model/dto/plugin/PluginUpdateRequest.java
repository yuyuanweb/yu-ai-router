/**
 * 插件更新请求
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.model.dto.plugin;

import lombok.Data;

import java.io.Serial;
import java.io.Serializable;

@Data
public class PluginUpdateRequest implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 插件 ID
     */
    private Long id;

    /**
     * 插件名称
     */
    private String pluginName;

    /**
     * 插件描述
     */
    private String description;

    /**
     * 插件配置（JSON）
     */
    private String config;

    /**
     * 状态：active/inactive
     */
    private String status;

    /**
     * 优先级
     */
    private Integer priority;
}
