package com.yupi.airouter.service;

import com.yupi.airouter.model.entity.RequestLog;

import java.util.List;

/**
 * 请求日志服务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface RequestLogService {

    /**
     * 记录请求日志
     *
     * @param userId 用户ID
     * @param apiKeyId API Key ID
     * @param modelName 模型名称
     * @param promptTokens 输入Token数
     * @param completionTokens 输出Token数
     * @param totalTokens 总Token数
     * @param duration 请求耗时
     * @param status 状态
     * @param errorMessage 错误信息
     */
    void logRequest(Long userId, Long apiKeyId, String modelName,
                    Integer promptTokens, Integer completionTokens, Integer totalTokens,
                    Integer duration, String status, String errorMessage);

    /**
     * 获取用户的请求日志
     *
     * @param userId 用户ID
     * @param limit 限制条数
     * @return 请求日志列表
     */
    List<RequestLog> listUserLogs(Long userId, Integer limit);

    /**
     * 统计用户的 Token 消耗
     *
     * @param userId 用户ID
     * @return 总Token数
     */
    Long countUserTokens(Long userId);
}
