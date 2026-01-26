/**
 * 插件执行结果
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.plugin;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.Map;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class PluginResult {

    /**
     * 是否执行成功
     */
    private boolean success;

    /**
     * 结果内容（可以直接添加到对话上下文）
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

    /**
     * 创建成功结果
     */
    public static PluginResult success(String content) {
        return PluginResult.builder()
                .success(true)
                .content(content)
                .build();
    }

    /**
     * 创建成功结果（带数据）
     */
    public static PluginResult success(String content, Map<String, Object> data) {
        return PluginResult.builder()
                .success(true)
                .content(content)
                .data(data)
                .build();
    }

    /**
     * 创建失败结果
     */
    public static PluginResult fail(String errorMessage) {
        return PluginResult.builder()
                .success(false)
                .errorMessage(errorMessage)
                .build();
    }
}
