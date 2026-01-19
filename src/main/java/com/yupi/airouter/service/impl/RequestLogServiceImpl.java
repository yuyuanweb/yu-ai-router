package com.yupi.airouter.service.impl;

import com.mybatisflex.core.query.QueryWrapper;
import com.yupi.airouter.mapper.RequestLogMapper;
import com.yupi.airouter.model.entity.RequestLog;
import com.yupi.airouter.service.ApiKeyService;
import com.yupi.airouter.service.RequestLogService;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

import jakarta.annotation.Resource;
import java.time.LocalDateTime;
import java.util.List;

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

    @Override
    @Async
    public void logRequest(Long userId, Long apiKeyId, String modelName,
                           Integer promptTokens, Integer completionTokens, Integer totalTokens,
                           Integer duration, String status, String errorMessage) {
        // 创建请求日志
        RequestLog log = RequestLog.builder()
                .userId(userId)
                .apiKeyId(apiKeyId)
                .modelName(modelName)
                .promptTokens(promptTokens)
                .completionTokens(completionTokens)
                .totalTokens(totalTokens)
                .duration(duration)
                .status(status)
                .errorMessage(errorMessage)
                .createTime(LocalDateTime.now())
                .updateTime(LocalDateTime.now())
                .build();

        // 保存日志
        requestLogMapper.insert(log);

        // 更新 API Key 的使用统计（仅成功的请求）
        if ("success".equals(status) && apiKeyId != null && totalTokens != null && totalTokens > 0) {
            apiKeyService.updateUsageStats(apiKeyId, totalTokens);
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
}
