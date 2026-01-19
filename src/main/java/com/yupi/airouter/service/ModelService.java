package com.yupi.airouter.service;

import com.mybatisflex.core.query.QueryWrapper;
import com.mybatisflex.core.service.IService;
import com.yupi.airouter.model.dto.model.ModelQueryRequest;
import com.yupi.airouter.model.entity.Model;
import com.yupi.airouter.model.vo.ModelVO;

import java.math.BigDecimal;
import java.util.List;

/**
 * 模型 服务层
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface ModelService extends IService<Model> {

    /**
     * 根据查询条件构造查询参数
     *
     * @param modelQueryRequest 查询请求
     * @return 查询参数
     */
    QueryWrapper getQueryWrapper(ModelQueryRequest modelQueryRequest);

    /**
     * 获取模型信息（含提供者信息）
     *
     * @param model 模型信息
     * @return 模型VO
     */
    ModelVO getModelVO(Model model);

    /**
     * 获取模型信息列表（含提供者信息）
     *
     * @param modelList 模型列表
     * @return 模型VO列表
     */
    List<ModelVO> getModelVOList(List<Model> modelList);

    /**
     * 根据模型Key获取模型
     *
     * @param modelKey 模型标识
     * @return 模型信息
     */
    Model getByModelKey(String modelKey);

    /**
     * 获取所有启用的模型
     *
     * @return 启用的模型列表
     */
    List<Model> getActiveModels();

    /**
     * 根据提供者ID获取所有启用的模型
     *
     * @param providerId 提供者ID
     * @return 启用的模型列表
     */
    List<Model> getActiveModelsByProviderId(Long providerId);

    /**
     * 根据模型类型获取所有启用的模型
     *
     * @param modelType 模型类型
     * @return 启用的模型列表
     */
    List<Model> getActiveModelsByType(String modelType);

    /**
     * 更新模型的健康指标
     *
     * @param modelId      模型ID
     * @param healthStatus 健康状态
     * @param avgLatency   平均延迟（毫秒）
     * @param successRate  成功率（百分比）
     * @param score        综合得分
     */
    void updateModelMetrics(Long modelId, String healthStatus, Integer avgLatency, 
                           BigDecimal successRate, BigDecimal score);
}
