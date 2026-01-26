/**
 * 插件执行结果视图对象
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
import java.util.Map;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class PluginExecuteVO implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 是否执行成功
     */
    private boolean success;

    /**
     * 插件标识
     */
    private String pluginKey;

    /**
     * 结果内容
     */
    private String content;

    /**
     * 错误信息
     */
    private String errorMessage;

    /**
     * 执行耗时（毫秒）
     */
    private long duration;

    /**
     * 额外数据
     */
    private Map<String, Object> data;
}
