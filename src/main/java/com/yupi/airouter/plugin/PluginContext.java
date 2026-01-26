/**
 * 插件执行上下文
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
public class PluginContext {

    /**
     * 用户 ID
     */
    private Long userId;

    /**
     * 输入内容（用户的问题或请求）
     */
    private String input;

    /**
     * 文件 URL（用于 PDF 解析、图片识别等）
     */
    private String fileUrl;

    /**
     * 文件字节数组（直接传递文件内容）
     */
    private byte[] fileBytes;

    /**
     * 文件类型（如 application/pdf, image/png）
     */
    private String fileType;

    /**
     * 额外参数
     */
    private Map<String, Object> params;

    /**
     * 获取参数值
     */
    @SuppressWarnings("unchecked")
    public <T> T getParam(String key) {
        if (params == null) {
            return null;
        }
        return (T) params.get(key);
    }

    /**
     * 获取参数值（带默认值）
     */
    @SuppressWarnings("unchecked")
    public <T> T getParam(String key, T defaultValue) {
        if (params == null) {
            return defaultValue;
        }
        T value = (T) params.get(key);
        return value != null ? value : defaultValue;
    }
}
