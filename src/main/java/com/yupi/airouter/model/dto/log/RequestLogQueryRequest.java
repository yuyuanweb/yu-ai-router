package com.yupi.airouter.model.dto.log;

import com.yupi.airouter.common.PageRequest;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.io.Serial;

/**
 * 请求日志查询请求
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@EqualsAndHashCode(callSuper = true)
public class RequestLogQueryRequest extends PageRequest {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 用户ID
     */
    private Long userId;

    /**
     * 模型标识
     */
    private String requestModel;

    /**
     * 请求类型：chat/embedding/image
     */
    private String requestType;

    /**
     * 调用来源：web/api
     */
    private String source;

    /**
     * 状态：success/failed
     */
    private String status;

    /**
     * 开始日期（格式：yyyy-MM-dd）
     */
    private String startDate;

    /**
     * 结束日期（格式：yyyy-MM-dd）
     */
    private String endDate;
}
