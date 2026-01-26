/**
 * 插件配置视图对象
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.model.vo;

import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.util.Date;

@Data
public class PluginConfigVO implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * id
     */
    private Long id;

    /**
     * 插件标识
     */
    private String pluginKey;

    /**
     * 插件名称
     */
    private String pluginName;

    /**
     * 插件类型
     */
    private String pluginType;

    /**
     * 插件描述
     */
    private String description;

    /**
     * 插件配置（JSON）
     */
    private String config;

    /**
     * 状态
     */
    private String status;

    /**
     * 优先级
     */
    private Integer priority;

    /**
     * 创建时间
     */
    private Date createTime;

    /**
     * 更新时间
     */
    private Date updateTime;
}
