package com.yupi.airouter.model.enums;

import lombok.Getter;

/**
 * 提供者状态枚举
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Getter
public enum ProviderStatusEnum {

    ACTIVE("active", "启用"),
    INACTIVE("inactive", "禁用"),
    MAINTENANCE("maintenance", "维护中");

    private final String value;
    private final String text;

    ProviderStatusEnum(String value, String text) {
        this.value = value;
        this.text = text;
    }

    /**
     * 根据 value 获取枚举
     */
    public static ProviderStatusEnum getEnumByValue(String value) {
        if (value == null) {
            return null;
        }
        for (ProviderStatusEnum anEnum : ProviderStatusEnum.values()) {
            if (anEnum.value.equals(value)) {
                return anEnum;
            }
        }
        return null;
    }
}
