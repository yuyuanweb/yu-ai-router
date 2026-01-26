/**
 * 插件配置实体
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.model.entity;

import com.mybatisflex.annotation.Column;
import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.annotation.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;
import java.util.Date;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("plugin_config")
public class PluginConfig implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * id
     */
    @Id(keyType = KeyType.Auto)
    private Long id;

    /**
     * 插件标识
     */
    @Column("pluginKey")
    private String pluginKey;

    /**
     * 插件名称
     */
    @Column("pluginName")
    private String pluginName;

    /**
     * 插件类型：builtin/custom
     */
    @Column("pluginType")
    private String pluginType;

    /**
     * 插件描述
     */
    @Column("description")
    private String description;

    /**
     * 插件配置（JSON）
     */
    @Column("config")
    private String config;

    /**
     * 状态：active/inactive
     */
    @Column("status")
    private String status;

    /**
     * 优先级
     */
    @Column("priority")
    private Integer priority;

    /**
     * 创建时间
     */
    @Column("createTime")
    private Date createTime;

    /**
     * 更新时间
     */
    @Column("updateTime")
    private Date updateTime;

    /**
     * 是否删除
     */
    @Column(value = "isDelete", isLogicDelete = true)
    private Integer isDelete;
}
