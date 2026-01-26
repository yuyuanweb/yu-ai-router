/**
 * 插件状态枚举
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.model.enums;

import lombok.Getter;

@Getter
public enum PluginStatusEnum {

    ACTIVE("active", "启用"),
    INACTIVE("inactive", "禁用");

    private final String value;
    private final String text;

    PluginStatusEnum(String value, String text) {
        this.value = value;
        this.text = text;
    }

    /**
     * 根据值获取枚举
     */
    public static PluginStatusEnum getEnumByValue(String value) {
        if (value == null) {
            return null;
        }
        for (PluginStatusEnum statusEnum : PluginStatusEnum.values()) {
            if (statusEnum.value.equals(value)) {
                return statusEnum;
            }
        }
        return null;
    }
}
