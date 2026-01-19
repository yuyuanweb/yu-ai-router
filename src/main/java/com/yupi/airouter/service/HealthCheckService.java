package com.yupi.airouter.service;

import com.yupi.airouter.model.entity.ModelProvider;

/**
 * 健康检查服务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface HealthCheckService {

    /**
     * 检查提供者健康状态
     *
     * @param provider 模型提供者
     * @return 是否健康
     */
    boolean checkProviderHealth(ModelProvider provider);

    /**
     * 检查所有提供者的健康状态
     */
    void checkAllProviders();

    /**
     * 获取提供者的平均延迟
     *
     * @param provider 模型提供者
     * @return 平均延迟（毫秒）
     */
    int measureLatency(ModelProvider provider);

    /**
     * 从请求日志统计各模型的指标并更新到 model 表和 model_provider 表
     * 包括：平均调用时间、成功率、健康状态、综合得分
     */
    void syncModelMetricsFromRequestLog();
}
