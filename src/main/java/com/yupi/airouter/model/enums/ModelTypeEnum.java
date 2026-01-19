package com.yupi.airouter.model.enums;

import lombok.Getter;

/**
 * 模型类型枚举
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Getter
public enum ModelTypeEnum {

    CHAT("chat", "对话模型"),
    EMBEDDING("embedding", "向量模型"),
    IMAGE("image", "图像模型"),
    AUDIO("audio", "音频模型");

    private final String value;
    private final String text;

    ModelTypeEnum(String value, String text) {
        this.value = value;
        this.text = text;
    }

    /**
     * 根据 value 获取枚举
     */
    public static ModelTypeEnum getEnumByValue(String value) {
        if (value == null) {
            return null;
        }
        for (ModelTypeEnum anEnum : ModelTypeEnum.values()) {
            if (anEnum.value.equals(value)) {
                return anEnum;
            }
        }
        return null;
    }
}
