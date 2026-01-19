package com.yupi.airouter.service;

import com.mybatisflex.core.query.QueryWrapper;
import com.mybatisflex.core.service.IService;
import com.yupi.airouter.model.dto.provider.ProviderQueryRequest;
import com.yupi.airouter.model.entity.ModelProvider;
import com.yupi.airouter.model.vo.ProviderVO;

import java.math.BigDecimal;
import java.util.List;

/**
 * 模型提供者 服务层
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface ModelProviderService extends IService<ModelProvider> {

    /**
     * 根据查询条件构造查询参数
     *
     * @param providerQueryRequest 查询请求
     * @return 查询参数
     */
    QueryWrapper getQueryWrapper(ProviderQueryRequest providerQueryRequest);

    /**
     * 获取脱敏后的提供者信息
     *
     * @param modelProvider 提供者信息
     * @return 脱敏后的提供者信息
     */
    ProviderVO getProviderVO(ModelProvider modelProvider);

    /**
     * 获取脱敏后的提供者信息列表
     *
     * @param providerList 提供者列表
     * @return 脱敏后的提供者信息列表
     */
    List<ProviderVO> getProviderVOList(List<ModelProvider> providerList);

    /**
     * 更新健康状态
     *
     * @param providerId    提供者ID
     * @param healthStatus  健康状态
     * @param avgLatency    平均延迟
     * @param successRate   成功率
     */
    void updateHealthStatus(Long providerId, String healthStatus, Integer avgLatency, BigDecimal successRate);

    /**
     * 获取所有健康的提供者
     *
     * @return 健康的提供者列表
     */
    List<ModelProvider> getHealthyProviders();
}
