package com.yupi.airouter.strategy;

import com.yupi.airouter.model.entity.Model;

import java.util.List;

/**
 * 路由策略接口
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface RoutingStrategyInterface {

    /**
     * 选择最优模型（从数据库查询并返回最符合策略的一条数据）
     *
     * @param modelType      模型类型（可选，如 chat/embedding）
     * @param requestedModel 用户请求的模型标识（固定模型策略使用）
     * @return 选中的模型，如果没有可用模型返回 null
     */
    Model selectModel(String modelType, String requestedModel);

    /**
     * 获取 Fallback 模型列表（除了主选模型外的备选模型）
     *
     * @param modelType      模型类型（可选）
     * @param requestedModel 用户请求的模型标识（固定模型策略使用）
     * @return Fallback 模型列表
     */
    List<Model> getFallbackModels(String modelType, String requestedModel);

    /**
     * 获取策略类型
     *
     * @return 策略类型
     */
    String getStrategyType();
}
