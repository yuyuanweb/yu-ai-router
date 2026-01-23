package com.yupi.airouter.service.impl;

import cn.hutool.core.util.IdUtil;
import cn.hutool.core.util.StrUtil;
import com.mybatisflex.core.paginate.Page;
import com.mybatisflex.core.query.QueryWrapper;
import com.yupi.airouter.mapper.RequestLogMapper;
import com.yupi.airouter.model.dto.log.RequestLogDTO;
import com.yupi.airouter.model.dto.log.RequestLogQueryRequest;
import com.yupi.airouter.model.entity.RequestLog;
import com.yupi.airouter.service.ApiKeyService;
import com.yupi.airouter.service.BillingService;
import com.yupi.airouter.service.RequestLogService;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

import jakarta.annotation.Resource;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * 请求日志服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Service
public class RequestLogServiceImpl implements RequestLogService {

    @Resource
    private RequestLogMapper requestLogMapper;

    @Resource
    private ApiKeyService apiKeyService;

    @Resource
    private BillingService billingService;

    @Override
    @Async
    public void logRequest(RequestLogDTO logDTO) {
        if (logDTO == null) {
            return;
        }

        // 如果没有 traceId，生成一个
        if (logDTO.getTraceId() == null) {
            logDTO.setTraceId(IdUtil.simpleUUID());
        }

        // 计算费用（如果没有提供费用，则自动计算）
        BigDecimal cost = logDTO.getCost();
        if (cost == null && logDTO.getModelId() != null && "success".equals(logDTO.getStatus())) {
            int promptTokens = logDTO.getPromptTokens() != null ? logDTO.getPromptTokens() : 0;
            int completionTokens = logDTO.getCompletionTokens() != null ? logDTO.getCompletionTokens() : 0;
            cost = billingService.calculateCost(logDTO.getModelId(), promptTokens, completionTokens);
        }

        // 创建请求日志
        RequestLog log = RequestLog.builder()
                .traceId(logDTO.getTraceId())
                .userId(logDTO.getUserId())
                .apiKeyId(logDTO.getApiKeyId())
                .modelId(logDTO.getModelId())
                .requestModel(logDTO.getRequestModel())
                .modelName(logDTO.getRequestModel())  // 兼容字段
                .requestType(logDTO.getRequestType() != null ? logDTO.getRequestType() : "chat")
                .source(logDTO.getSource() != null ? logDTO.getSource() : "web")
                .promptTokens(logDTO.getPromptTokens() != null ? logDTO.getPromptTokens() : 0)
                .completionTokens(logDTO.getCompletionTokens() != null ? logDTO.getCompletionTokens() : 0)
                .totalTokens(logDTO.getTotalTokens() != null ? logDTO.getTotalTokens() : 0)
                .cost(cost != null ? cost : BigDecimal.ZERO)
                .duration(logDTO.getDuration() != null ? logDTO.getDuration() : 0)
                .status(logDTO.getStatus())
                .errorMessage(logDTO.getErrorMessage())
                .errorCode(logDTO.getErrorCode())
                .routingStrategy(logDTO.getRoutingStrategy())
                .isFallback(logDTO.getIsFallback() != null && logDTO.getIsFallback() ? 1 : 0)
                .clientIp(logDTO.getClientIp())
                .userAgent(logDTO.getUserAgent())
                .createTime(LocalDateTime.now())
                .updateTime(LocalDateTime.now())
                .build();

        // 保存日志
        requestLogMapper.insert(log);

        // 更新 API Key 的使用统计（仅成功的请求且有 apiKeyId）
        if ("success".equals(logDTO.getStatus()) && 
            logDTO.getApiKeyId() != null && 
            logDTO.getTotalTokens() != null && 
            logDTO.getTotalTokens() > 0) {
            apiKeyService.updateUsageStats(logDTO.getApiKeyId(), logDTO.getTotalTokens());
        }
    }

    @Override
    public List<RequestLog> listUserLogs(Long userId, Integer limit) {
        QueryWrapper queryWrapper = QueryWrapper.create()
                .eq("userId", userId)
                .orderBy("createTime", false)
                .limit(limit != null ? limit : 100);

        return requestLogMapper.selectListByQuery(queryWrapper);
    }

    @Override
    public Long countUserTokens(Long userId) {
        QueryWrapper queryWrapper = QueryWrapper.create()
                .select("SUM(totalTokens)")
                .eq("userId", userId)
                .eq("status", "success");

        Long total = requestLogMapper.selectObjectByQueryAs(queryWrapper, Long.class);
        return total != null ? total : 0L;
    }

    @Override
    public Page<RequestLog> pageUserLogs(Long userId, long pageNum, long pageSize) {
        QueryWrapper queryWrapper = QueryWrapper.create()
                .eq("userId", userId)
                .orderBy("createTime", false);

        return requestLogMapper.paginate(Page.of(pageNum, pageSize), queryWrapper);
    }

    @Override
    public List<Map<String, Object>> getUserDailyStats(Long userId, LocalDate startDate, LocalDate endDate) {
        List<Map<String, Object>> result = new ArrayList<>();
        
        LocalDate currentDate = startDate;
        while (!currentDate.isAfter(endDate)) {
            LocalDateTime dayStart = LocalDateTime.of(currentDate, LocalTime.MIN);
            LocalDateTime dayEnd = LocalDateTime.of(currentDate, LocalTime.MAX);
            
            List<RequestLog> logs = requestLogMapper.selectListByQuery(
                    QueryWrapper.create()
                            .where("userId = " + userId)
                            .and("createTime >= '" + dayStart + "'")
                            .and("createTime <= '" + dayEnd + "'")
            );
            
            long totalTokens = 0;
            long requestCount = 0;
            long successCount = 0;
            BigDecimal totalCost = BigDecimal.ZERO;
            
            for (RequestLog log : logs) {
                requestCount++;
                if ("success".equals(log.getStatus())) {
                    successCount++;
                    if (log.getTotalTokens() != null) {
                        totalTokens += log.getTotalTokens();
                    }
                    if (log.getCost() != null) {
                        totalCost = totalCost.add(log.getCost());
                    }
                }
            }
            
            Map<String, Object> dayStats = new HashMap<>();
            dayStats.put("date", currentDate.toString());
            dayStats.put("totalTokens", totalTokens);
            dayStats.put("requestCount", requestCount);
            dayStats.put("successCount", successCount);
            dayStats.put("totalCost", totalCost);
            result.add(dayStats);
            
            currentDate = currentDate.plusDays(1);
        }
        
        return result;
    }

    @Override
    public Long countUserRequests(Long userId) {
        QueryWrapper queryWrapper = QueryWrapper.create()
                .eq("userId", userId);
        return requestLogMapper.selectCountByQuery(queryWrapper);
    }

    @Override
    public Long countUserSuccessRequests(Long userId) {
        QueryWrapper queryWrapper = QueryWrapper.create()
                .eq("userId", userId)
                .eq("status", "success");
        return requestLogMapper.selectCountByQuery(queryWrapper);
    }

    @Override
    public Page<RequestLog> pageByQuery(RequestLogQueryRequest queryRequest) {
        QueryWrapper queryWrapper = QueryWrapper.create();
        
        // 用户ID条件
        if (queryRequest.getUserId() != null) {
            queryWrapper.eq("userId", queryRequest.getUserId());
        }
        
        // 模型标识
        if (StrUtil.isNotBlank(queryRequest.getRequestModel())) {
            queryWrapper.like("requestModel", queryRequest.getRequestModel());
        }
        
        // 请求类型
        if (StrUtil.isNotBlank(queryRequest.getRequestType())) {
            queryWrapper.eq("requestType", queryRequest.getRequestType());
        }
        
        // 调用来源
        if (StrUtil.isNotBlank(queryRequest.getSource())) {
            queryWrapper.eq("source", queryRequest.getSource());
        }
        
        // 状态
        if (StrUtil.isNotBlank(queryRequest.getStatus())) {
            queryWrapper.eq("status", queryRequest.getStatus());
        }
        
        // 日期范围
        if (StrUtil.isNotBlank(queryRequest.getStartDate())) {
            LocalDateTime startTime = LocalDateTime.of(LocalDate.parse(queryRequest.getStartDate()), LocalTime.MIN);
            queryWrapper.ge("createTime", startTime);
        }
        if (StrUtil.isNotBlank(queryRequest.getEndDate())) {
            LocalDateTime endTime = LocalDateTime.of(LocalDate.parse(queryRequest.getEndDate()), LocalTime.MAX);
            queryWrapper.le("createTime", endTime);
        }
        
        // 按创建时间倒序
        queryWrapper.orderBy("createTime", false);
        
        long pageNum = queryRequest.getPageNum();
        long pageSize = queryRequest.getPageSize();
        return requestLogMapper.paginate(Page.of(pageNum, pageSize), queryWrapper);
    }

    @Override
    public RequestLog getById(Long id) {
        if (id == null) {
            return null;
        }
        return requestLogMapper.selectOneById(id);
    }
}
