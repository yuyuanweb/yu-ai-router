package com.yupi.airouter.service;

import com.yupi.airouter.model.entity.Model;

import java.util.List;

/**
 * 路由服务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface RoutingService {

    /**
     * 根据策略选择模型
     *
     * @param strategyType   策略类型
     * @param modelType      模型类型（可选）
     * @param requestedModel 用户请求的模型标识（固定模型策略使用）
     * @return 选中的模型
     */
    Model selectModel(String strategyType, String modelType, String requestedModel);

    /**
     * 获取 Fallback 模型列表
     *
     * @param strategyType   策略类型
     * @param modelType      模型类型（可选）
     * @param requestedModel 用户请求的模型标识（固定模型策略使用）
     * @return Fallback 模型列表
     */
    List<Model> getFallbackModels(String strategyType, String modelType, String requestedModel);
}
