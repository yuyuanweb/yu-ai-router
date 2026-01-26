/**
 * 插件类型枚举
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.model.enums;

import lombok.Getter;

@Getter
public enum PluginTypeEnum {

    BUILTIN("builtin", "内置插件"),
    CUSTOM("custom", "自定义插件");

    private final String value;
    private final String text;

    PluginTypeEnum(String value, String text) {
        this.value = value;
        this.text = text;
    }

    /**
     * 根据值获取枚举
     */
    public static PluginTypeEnum getEnumByValue(String value) {
        if (value == null) {
            return null;
        }
        for (PluginTypeEnum typeEnum : PluginTypeEnum.values()) {
            if (typeEnum.value.equals(value)) {
                return typeEnum;
            }
        }
        return null;
    }
}
