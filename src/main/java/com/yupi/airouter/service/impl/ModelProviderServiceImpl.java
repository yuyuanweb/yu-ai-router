package com.yupi.airouter.service.impl;

import cn.hutool.core.bean.BeanUtil;
import cn.hutool.core.collection.CollUtil;
import cn.hutool.core.util.StrUtil;
import com.mybatisflex.core.query.QueryWrapper;
import com.mybatisflex.spring.service.impl.ServiceImpl;
import com.yupi.airouter.mapper.ModelProviderMapper;
import com.yupi.airouter.model.dto.provider.ProviderQueryRequest;
import com.yupi.airouter.model.entity.ModelProvider;
import com.yupi.airouter.model.enums.HealthStatusEnum;
import com.yupi.airouter.model.enums.ProviderStatusEnum;
import com.yupi.airouter.model.vo.ProviderVO;
import com.yupi.airouter.service.ModelProviderService;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;

/**
 * 模型提供者 服务层实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Service
public class ModelProviderServiceImpl extends ServiceImpl<ModelProviderMapper, ModelProvider> implements ModelProviderService {

    @Override
    public QueryWrapper getQueryWrapper(ProviderQueryRequest providerQueryRequest) {
        if (providerQueryRequest == null) {
            return QueryWrapper.create();
        }

        String providerName = providerQueryRequest.getProviderName();
        String displayName = providerQueryRequest.getDisplayName();
        String status = providerQueryRequest.getStatus();
        String healthStatus = providerQueryRequest.getHealthStatus();

        // 构造查询条件
        return QueryWrapper.create()
                .like("providerName", providerName, StrUtil.isNotBlank(providerName))
                .like("displayName", displayName, StrUtil.isNotBlank(displayName))
                .eq("status", status, StrUtil.isNotBlank(status))
                .eq("healthStatus", healthStatus, StrUtil.isNotBlank(healthStatus))
                .orderBy("priority", false)
                .orderBy("createTime", false);
    }

    @Override
    public ProviderVO getProviderVO(ModelProvider modelProvider) {
        if (modelProvider == null) {
            return null;
        }
        ProviderVO providerVO = new ProviderVO();
        BeanUtil.copyProperties(modelProvider, providerVO);
        // 脱敏：不返回完整的 API Key
        // providerVO 中没有 apiKey 字段，所以已经自动脱敏
        return providerVO;
    }

    @Override
    public List<ProviderVO> getProviderVOList(List<ModelProvider> providerList) {
        if (CollUtil.isEmpty(providerList)) {
            return new ArrayList<>();
        }
        return providerList.stream()
                .map(this::getProviderVO)
                .collect(Collectors.toList());
    }

    @Override
    public void updateHealthStatus(Long providerId, String healthStatus, Integer avgLatency, BigDecimal successRate) {
        ModelProvider updateProvider = new ModelProvider();
        updateProvider.setId(providerId);
        updateProvider.setHealthStatus(healthStatus);
        updateProvider.setAvgLatency(avgLatency);
        updateProvider.setSuccessRate(successRate);
        this.updateById(updateProvider);
    }

    @Override
    public List<ModelProvider> getHealthyProviders() {
        return this.list(QueryWrapper.create()
                .eq("status", ProviderStatusEnum.ACTIVE.getValue())
                .in("healthStatus", 
                        HealthStatusEnum.HEALTHY.getValue(), 
                        HealthStatusEnum.DEGRADED.getValue())
                .orderBy("priority", false));
    }
}
