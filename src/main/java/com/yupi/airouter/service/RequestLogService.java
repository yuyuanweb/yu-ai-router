package com.yupi.airouter.service;

import com.mybatisflex.core.paginate.Page;
import com.yupi.airouter.model.dto.log.RequestLogDTO;
import com.yupi.airouter.model.dto.log.RequestLogQueryRequest;
import com.yupi.airouter.model.entity.RequestLog;

import java.time.LocalDate;
import java.util.List;
import java.util.Map;

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

    /**
     * 分页获取用户请求日志
     *
     * @param userId   用户ID
     * @param pageNum  页码
     * @param pageSize 每页大小
     * @return 分页结果
     */
    Page<RequestLog> pageUserLogs(Long userId, long pageNum, long pageSize);

    /**
     * 获取用户按日期统计的消耗数据
     *
     * @param userId    用户ID
     * @param startDate 开始日期
     * @param endDate   结束日期
     * @return 日期 -> 统计数据的映射
     */
    List<Map<String, Object>> getUserDailyStats(Long userId, LocalDate startDate, LocalDate endDate);

    /**
     * 获取用户请求总数
     *
     * @param userId 用户ID
     * @return 请求总数
     */
    Long countUserRequests(Long userId);

    /**
     * 获取用户成功请求数
     *
     * @param userId 用户ID
     * @return 成功请求数
     */
    Long countUserSuccessRequests(Long userId);

    /**
     * 按条件分页查询请求日志
     *
     * @param queryRequest 查询请求
     * @return 分页结果
     */
    Page<RequestLog> pageByQuery(RequestLogQueryRequest queryRequest);

    /**
     * 根据ID获取请求日志
     *
     * @param id 日志ID
     * @return 请求日志
     */
    RequestLog getById(Long id);
}
