/**
 * 插件执行请求
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.model.dto.plugin;

import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.util.Map;

@Data
public class PluginExecuteRequest implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 插件标识
     */
    private String pluginKey;

    /**
     * 输入内容
     */
    private String input;

    /**
     * 文件 URL（可选）
     */
    private String fileUrl;

    /**
     * 文件字节数组（直接传递文件内容）
     */
    private transient byte[] fileBytes;

    /**
     * 文件类型（可选）
     */
    private String fileType;

    /**
     * 额外参数
     */
    private Map<String, Object> params;
}
