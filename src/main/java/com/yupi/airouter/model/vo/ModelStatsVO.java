/**
 * 模型统计 VO
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.model.vo;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;

/**
 * 模型调用统计视图对象
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ModelStatsVO implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 模型标识
     */
    private String modelKey;

    /**
     * 请求数
     */
    private Long requestCount;

    /**
     * 成功数
     */
    private Long successCount;

    /**
     * 失败数
     */
    private Long failedCount;

    /**
     * 成功率
     */
    private Double successRate;

    /**
     * Token 消耗
     */
    private Long totalTokens;

    /**
     * 平均响应时间
     */
    private Double avgDuration;
}
