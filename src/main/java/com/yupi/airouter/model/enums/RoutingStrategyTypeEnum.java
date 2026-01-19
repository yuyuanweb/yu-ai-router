package com.yupi.airouter.model.enums;

import lombok.Getter;

/**
 * 路由策略类型枚举
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Getter
public enum RoutingStrategyTypeEnum {

    AUTO("auto", "自动路由"),
    COST_FIRST("cost_first", "成本优先"),
    LATENCY_FIRST("latency_first", "延迟优先"),
    ROUND_ROBIN("round_robin", "轮询"),
    FIXED("fixed", "固定模型");

    private final String value;
    private final String text;

    RoutingStrategyTypeEnum(String value, String text) {
        this.value = value;
        this.text = text;
    }

    /**
     * 根据 value 获取枚举
     */
    public static RoutingStrategyTypeEnum getEnumByValue(String value) {
        if (value == null) {
            return null;
        }
        for (RoutingStrategyTypeEnum anEnum : RoutingStrategyTypeEnum.values()) {
            if (anEnum.value.equals(value)) {
                return anEnum;
            }
        }
        return null;
    }
}
