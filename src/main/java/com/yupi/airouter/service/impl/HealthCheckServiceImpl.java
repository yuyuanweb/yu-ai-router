package com.yupi.airouter.service.impl;

import cn.hutool.core.collection.CollUtil;
import cn.hutool.http.HttpRequest;
import cn.hutool.http.HttpResponse;
import com.mybatisflex.core.query.QueryWrapper;
import com.yupi.airouter.model.entity.Model;
import com.yupi.airouter.model.entity.ModelProvider;
import com.yupi.airouter.model.enums.HealthStatusEnum;
import com.yupi.airouter.model.enums.ProviderStatusEnum;
import com.yupi.airouter.mapper.RequestLogMapper;
import com.yupi.airouter.service.HealthCheckService;
import com.yupi.airouter.service.ModelProviderService;
import com.yupi.airouter.service.ModelService;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.stream.Collectors;

/**
 * 健康检查服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Slf4j
@Service
public class HealthCheckServiceImpl implements HealthCheckService {

    @Resource
    private ModelProviderService modelProviderService;

    @Resource
    private ModelService modelService;

    @Resource
    private RequestLogMapper requestLogMapper;

    /**
     * 存储每个提供者最近的健康检查结果（用于计算成功率）
     * Key: providerId, Value: 最近100次检查结果列表（true=成功，false=失败）
     */
    private final ConcurrentHashMap<Long, List<Boolean>> healthHistoryMap = new ConcurrentHashMap<>();

    /**
     * 健康检查历史记录最大数量
     */
    private static final int MAX_HISTORY_SIZE = 100;

    /**
     * 健康检查超时时间（毫秒）
     */
    private static final int HEALTH_CHECK_TIMEOUT = 10000;

    /**
     * 统计请求日志的时间范围（小时）
     */
    private static final int STATS_HOURS = 24;

    /**
     * 成本权重（用于计算综合得分）
     */
    private static final double COST_WEIGHT = 0.3;

    /**
     * 延迟权重（用于计算综合得分）
     */
    private static final double LATENCY_WEIGHT = 0.3;

    /**
     * 成功率权重（用于计算综合得分）
     */
    private static final double SUCCESS_RATE_WEIGHT = 0.2;

    /**
     * 优先级权重（用于计算综合得分）
     */
    private static final double PRIORITY_WEIGHT = 0.2;

    @Override
    public boolean checkProviderHealth(ModelProvider provider) {
        if (provider == null) {
            return false;
        }

        try {
            // 测量延迟
            long startTime = System.currentTimeMillis();
            boolean isHealthy = sendHealthCheckRequest(provider);
            long endTime = System.currentTimeMillis();
            int latency = (int) (endTime - startTime);

            // 记录健康检查结果
            recordHealthCheckResult(provider.getId(), isHealthy);

            // 计算成功率
            BigDecimal successRate = calculateSuccessRate(provider.getId());

            // 确定健康状态
            String healthStatus = determineHealthStatus(successRate);

            // 更新提供者的健康信息
            modelProviderService.updateHealthStatus(
                    provider.getId(),
                    healthStatus,
                    latency,
                    successRate
            );

            log.debug("提供者 {} 健康检查完成: 健康={}, 延迟={}ms, 成功率={}%",
                    provider.getDisplayName(), isHealthy, latency, successRate);

            return isHealthy;
        } catch (Exception e) {
            log.error("提供者 {} 健康检查失败", provider.getDisplayName(), e);
            
            // 记录失败
            recordHealthCheckResult(provider.getId(), false);
            BigDecimal successRate = calculateSuccessRate(provider.getId());
            
            // 更新为不健康状态
            modelProviderService.updateHealthStatus(
                    provider.getId(),
                    HealthStatusEnum.UNHEALTHY.getValue(),
                    0,
                    successRate
            );
            
            return false;
        }
    }

    @Override
    public void checkAllProviders() {
        // 获取所有启用的提供者
        List<ModelProvider> providers = modelProviderService.list();
        if (CollUtil.isEmpty(providers)) {
            log.warn("没有找到任何模型提供者");
            return;
        }

        log.info("开始健康检查，共 {} 个提供者", providers.size());

        // 并行检查所有提供者
        providers.parallelStream()
                .filter(provider -> ProviderStatusEnum.ACTIVE.getValue().equals(provider.getStatus()))
                .forEach(this::checkProviderHealth);

        // 从请求日志同步模型指标
        syncModelMetricsFromRequestLog();

        log.info("健康检查完成");
    }

    @Override
    public int measureLatency(ModelProvider provider) {
        try {
            long startTime = System.currentTimeMillis();
            sendHealthCheckRequest(provider);
            long endTime = System.currentTimeMillis();
            return (int) (endTime - startTime);
        } catch (Exception e) {
            log.error("测量提供者 {} 延迟失败", provider.getDisplayName(), e);
            return Integer.MAX_VALUE;
        }
    }

    @Override
    public void syncModelMetricsFromRequestLog() {
        log.info("开始从请求日志同步模型指标");

        // 获取所有启用的模型
        List<Model> activeModels = modelService.getActiveModels();
        if (CollUtil.isEmpty(activeModels)) {
            log.warn("没有找到任何启用的模型");
            return;
        }

        // 计算统计时间范围
        LocalDateTime startTime = LocalDateTime.now().minusHours(STATS_HOURS);

        // 使用数据库聚合函数批量查询所有模型的统计数据
        List<ModelStats> modelStatsList = queryModelStatsFromDb(startTime);

        // 转换为 Map 便于查找
        Map<Long, ModelStats> modelStatsMap = modelStatsList.stream()
                .collect(Collectors.toMap(s -> s.modelId, s -> s, (a, b) -> a));

        // 计算归一化参数（用于综合得分计算）
        NormalizationParams normParams = calculateNormalizationParams(activeModels, modelStatsMap);

        // 更新每个模型的指标
        for (Model model : activeModels) {
            ModelStats stats = modelStatsMap.getOrDefault(model.getId(), new ModelStats());
            
            // 计算综合得分
            BigDecimal score = calculateScore(model, stats, normParams);
            
            // 根据成功率确定健康状态
            String healthStatus = determineHealthStatusBySuccessRate(stats.successRate);

            // 更新模型指标
            modelService.updateModelMetrics(
                    model.getId(),
                    healthStatus,
                    stats.avgLatency,
                    stats.successRate,
                    score
            );

            log.debug("模型 {} 指标更新: 延迟={}ms, 成功率={}%, 得分={}",
                    model.getModelKey(), stats.avgLatency, stats.successRate, score);
        }

        // 同步更新提供者的汇总指标（基于模型统计数据）
        syncProviderMetricsFromModelStats(activeModels, modelStatsMap);

        log.info("模型指标同步完成");
    }

    /**
     * 使用数据库聚合函数批量查询所有模型的统计数据
     */
    private List<ModelStats> queryModelStatsFromDb(LocalDateTime startTime) {
        // 使用 SQL 聚合函数一次性查询所有模型的统计数据
        // SELECT modelId, 
        //        AVG(CASE WHEN status = 'success' THEN duration ELSE NULL END) as avgLatency,
        //        COUNT(*) as totalRequests,
        //        SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) * 100.0 / COUNT(*) as successRate
        // FROM request_log 
        // WHERE createTime >= ? AND modelId IS NOT NULL
        // GROUP BY modelId
        
        QueryWrapper queryWrapper = QueryWrapper.create()
                .select("modelId",
                        "AVG(CASE WHEN status = 'success' THEN duration ELSE NULL END) as avgLatency",
                        "COUNT(*) as totalRequests",
                        "IFNULL(SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) * 100.0 / NULLIF(COUNT(*), 0), 100) as successRate")
                .from("request_log")
                .ge("createTime", startTime)
                .isNotNull("modelId")
                .groupBy("modelId");

        @SuppressWarnings("unchecked")
        List<Map<String, Object>> results = (List<Map<String, Object>>) (List<?>) requestLogMapper.selectListByQueryAs(queryWrapper, Map.class);

        List<ModelStats> statsList = new ArrayList<>();
        for (Map<String, Object> row : results) {
            ModelStats stats = new ModelStats();
            stats.modelId = row.get("modelId") != null ? ((Number) row.get("modelId")).longValue() : null;
            stats.avgLatency = row.get("avgLatency") != null ? ((Number) row.get("avgLatency")).intValue() : 0;
            stats.totalRequests = row.get("totalRequests") != null ? ((Number) row.get("totalRequests")).intValue() : 0;
            stats.successRate = row.get("successRate") != null 
                    ? BigDecimal.valueOf(((Number) row.get("successRate")).doubleValue()).setScale(2, RoundingMode.HALF_UP)
                    : BigDecimal.valueOf(100);
            statsList.add(stats);
        }

        return statsList;
    }

    /**
     * 计算归一化参数
     */
    private NormalizationParams calculateNormalizationParams(List<Model> models, Map<Long, ModelStats> statsMap) {
        NormalizationParams params = new NormalizationParams();

        // 成本范围
        params.minCost = models.stream()
                .map(m -> m.getInputPrice().add(m.getOutputPrice()))
                .min(BigDecimal::compareTo)
                .orElse(BigDecimal.ZERO);
        params.maxCost = models.stream()
                .map(m -> m.getInputPrice().add(m.getOutputPrice()))
                .max(BigDecimal::compareTo)
                .orElse(BigDecimal.ONE);

        // 延迟范围
        params.minLatency = statsMap.values().stream()
                .mapToInt(s -> s.avgLatency)
                .filter(l -> l > 0)
                .min()
                .orElse(0);
        params.maxLatency = statsMap.values().stream()
                .mapToInt(s -> s.avgLatency)
                .max()
                .orElse(10000);

        // 优先级范围
        params.minPriority = models.stream()
                .mapToInt(Model::getPriority)
                .min()
                .orElse(0);
        params.maxPriority = models.stream()
                .mapToInt(Model::getPriority)
                .max()
                .orElse(100);

        return params;
    }

    /**
     * 计算模型的综合得分（越低越好）
     */
    private BigDecimal calculateScore(Model model, ModelStats stats, NormalizationParams params) {
        // 1. 成本得分（归一化到 0-1）
        BigDecimal modelCost = model.getInputPrice().add(model.getOutputPrice());
        double costScore = normalize(modelCost, params.minCost, params.maxCost);

        // 2. 延迟得分（归一化到 0-1）
        int latency = stats.avgLatency > 0 ? stats.avgLatency : 5000; // 默认5秒
        double latencyScore = normalize(latency, params.minLatency, params.maxLatency);

        // 3. 成功率得分（归一化到 0-1，成功率越高得分越低）
        double successRateScore = 1.0 - (stats.successRate.doubleValue() / 100.0);

        // 4. 优先级得分（归一化到 0-1，优先级越高得分越低）
        double priorityScore = 1.0 - normalize(model.getPriority(), params.minPriority, params.maxPriority);

        // 综合得分
        double score = costScore * COST_WEIGHT +
                      latencyScore * LATENCY_WEIGHT +
                      successRateScore * SUCCESS_RATE_WEIGHT +
                      priorityScore * PRIORITY_WEIGHT;

        return BigDecimal.valueOf(score).setScale(4, RoundingMode.HALF_UP);
    }

    /**
     * 归一化数值到 0-1 范围
     */
    private double normalize(BigDecimal value, BigDecimal min, BigDecimal max) {
        BigDecimal range = max.subtract(min);
        if (range.compareTo(BigDecimal.ZERO) == 0) {
            return 0.0;
        }
        return value.subtract(min).divide(range, 4, RoundingMode.HALF_UP).doubleValue();
    }

    /**
     * 归一化数值到 0-1 范围
     */
    private double normalize(int value, int min, int max) {
        int range = max - min;
        if (range == 0) {
            return 0.0;
        }
        return (double) (value - min) / range;
    }

    /**
     * 根据成功率确定健康状态
     */
    private String determineHealthStatusBySuccessRate(BigDecimal successRate) {
        if (successRate == null) {
            return HealthStatusEnum.UNKNOWN.getValue();
        }
        
        double rate = successRate.doubleValue();
        if (rate >= 80.0) {
            return HealthStatusEnum.HEALTHY.getValue();
        } else if (rate >= 50.0) {
            return HealthStatusEnum.DEGRADED.getValue();
        } else {
            return HealthStatusEnum.UNHEALTHY.getValue();
        }
    }

    /**
     * 基于模型统计数据同步更新提供者的汇总指标
     */
    private void syncProviderMetricsFromModelStats(List<Model> models, Map<Long, ModelStats> statsMap) {
        // 按提供者分组
        Map<Long, List<Model>> providerModels = models.stream()
                .collect(Collectors.groupingBy(Model::getProviderId));

        for (Map.Entry<Long, List<Model>> entry : providerModels.entrySet()) {
            Long providerId = entry.getKey();
            List<Model> providerModelList = entry.getValue();

            // 计算该提供者所有模型的平均延迟和成功率
            int totalLatency = 0;
            int latencyCount = 0;
            double totalSuccessRate = 0;
            int successRateCount = 0;

            for (Model model : providerModelList) {
                ModelStats stats = statsMap.get(model.getId());
                if (stats != null && stats.avgLatency > 0) {
                    totalLatency += stats.avgLatency;
                    latencyCount++;
                }
                if (stats != null && stats.totalRequests > 0) {
                    totalSuccessRate += stats.successRate.doubleValue();
                    successRateCount++;
                }
            }

            int avgLatency = latencyCount > 0 ? totalLatency / latencyCount : 0;
            BigDecimal avgSuccessRate = successRateCount > 0 
                    ? BigDecimal.valueOf(totalSuccessRate / successRateCount).setScale(2, RoundingMode.HALF_UP)
                    : BigDecimal.valueOf(100);

            // 根据成功率确定健康状态
            String healthStatus = determineHealthStatusBySuccessRate(avgSuccessRate);

            modelProviderService.updateHealthStatus(providerId, healthStatus, avgLatency, avgSuccessRate);
            
            log.debug("提供者 {} 指标更新: 延迟={}ms, 成功率={}%, 健康状态={}",
                    providerId, avgLatency, avgSuccessRate, healthStatus);
        }
    }

    /**
     * 发送健康检查请求
     */
    private boolean sendHealthCheckRequest(ModelProvider provider) {
        try {
            String baseUrl = provider.getBaseUrl();
            String apiKey = provider.getApiKey();

            // 构造一个简单的测试请求（检查模型列表接口）
            String url = baseUrl + "/models";
            
            HttpResponse response = HttpRequest.get(url)
                    .header("Authorization", "Bearer " + apiKey)
                    .timeout(HEALTH_CHECK_TIMEOUT)
                    .execute();

            // 检查响应状态码
            return response.getStatus() == 200 || response.getStatus() == 401;
        } catch (Exception e) {
            log.warn("提供者 {} 健康检查请求失败: {}", provider.getDisplayName(), e.getMessage());
            return false;
        }
    }

    /**
     * 记录健康检查结果
     */
    private void recordHealthCheckResult(Long providerId, boolean isHealthy) {
        healthHistoryMap.computeIfAbsent(providerId, k -> new ArrayList<>());
        List<Boolean> history = healthHistoryMap.get(providerId);

        history.add(isHealthy);

        if (history.size() > MAX_HISTORY_SIZE) {
            history.remove(0);
        }
    }

    /**
     * 计算成功率
     */
    private BigDecimal calculateSuccessRate(Long providerId) {
        List<Boolean> history = healthHistoryMap.get(providerId);
        if (CollUtil.isEmpty(history)) {
            return BigDecimal.valueOf(100.00);
        }

        long successCount = history.stream().filter(result -> result).count();
        double rate = (double) successCount / history.size() * 100;

        return BigDecimal.valueOf(rate).setScale(2, RoundingMode.HALF_UP);
    }

    /**
     * 确定健康状态
     */
    private String determineHealthStatus(BigDecimal successRate) {
        double rate = successRate.doubleValue();
        
        if (rate >= 80.0) {
            return HealthStatusEnum.HEALTHY.getValue();
        } else if (rate >= 50.0) {
            return HealthStatusEnum.DEGRADED.getValue();
        } else {
            return HealthStatusEnum.UNHEALTHY.getValue();
        }
    }

    /**
     * 模型统计数据内部类
     */
    private static class ModelStats {
        Long modelId;
        int avgLatency = 0;
        BigDecimal successRate = BigDecimal.valueOf(100);
        int totalRequests = 0;
    }

    /**
     * 归一化参数内部类
     */
    private static class NormalizationParams {
        BigDecimal minCost = BigDecimal.ZERO;
        BigDecimal maxCost = BigDecimal.ONE;
        int minLatency = 0;
        int maxLatency = 10000;
        int minPriority = 0;
        int maxPriority = 100;
    }
}
