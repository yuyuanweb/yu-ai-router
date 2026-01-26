package com.yupi.airouter.model.dto.chat;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.util.List;

/**
 * 聊天请求（兼容 OpenAI 格式）
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
public class ChatRequest implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 模型名称（如：qwen-plus）
     */
    private String model;

    /**
     * 消息列表
     */
    private List<ChatMessage> messages;

    /**
     * 是否流式返回
     */
    private Boolean stream = false;

    /**
     * 温度参数（0-1）
     */
    private Double temperature;

    /**
     * 最大生成Token数
     */
    @JsonProperty("max_tokens")
    private Integer maxTokens;

    /**
     * 是否启用深度思考
     * 仅部分模型支持
     */
    @JsonProperty("enable_reasoning")
    private Boolean enableReasoning;

    /**
     * 路由策略类型
     * auto: 自动路由（综合考虑成本、延迟、优先级）
     * cost_first: 成本优先
     * latency_first: 延迟优先
     * fixed: 固定模型（使用指定的 model）
     */
    @JsonProperty("routing_strategy")
    private String routingStrategy;

    // ========== 插件相关参数 ==========

    /**
     * 插件标识
     * web_search: Web搜索插件
     * pdf_parser: PDF解析插件
     * image_recognition: 图片识别插件
     */
    @JsonProperty("plugin_key")
    private String pluginKey;

    /**
     * 文件 URL（用于 PDF 解析、图片识别）
     */
    @JsonProperty("file_url")
    private String fileUrl;

    /**
     * 文件字节数组（直接传递文件内容）
     * 注意：此字段不通过 JSON 传递，仅用于内部文件上传接口
     */
    @JsonIgnore
    private transient byte[] fileBytes;

    /**
     * 文件类型（如 image/png、application/pdf）
     */
    @JsonProperty("file_type")
    private String fileType;
}
