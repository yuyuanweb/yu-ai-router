package com.yupi.airouter.model.enums;

import lombok.Getter;

/**
 * 模型状态枚举
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Getter
public enum ModelStatusEnum {

    ACTIVE("active", "启用"),
    INACTIVE("inactive", "禁用"),
    DEPRECATED("deprecated", "已废弃");

    private final String value;
    private final String text;

    ModelStatusEnum(String value, String text) {
        this.value = value;
        this.text = text;
    }

    /**
     * 根据 value 获取枚举
     */
    public static ModelStatusEnum getEnumByValue(String value) {
        if (value == null) {
            return null;
        }
        for (ModelStatusEnum anEnum : ModelStatusEnum.values()) {
            if (anEnum.value.equals(value)) {
                return anEnum;
            }
        }
        return null;
    }
}
