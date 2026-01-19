package com.yupi.airouter.service;

import com.yupi.airouter.model.dto.log.RequestLogDTO;
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
     * @param logDTO 日志记录对象
     */
    void logRequest(RequestLogDTO logDTO);

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
