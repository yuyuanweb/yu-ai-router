package com.yupi.airouter.model.dto.image;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;

/**
 * 图片生成请求
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ImageGenerationRequest implements Serializable {

    /**
     * 提示词
     */
    private String prompt;

    /**
     * 模型名称（默认 qwen-image-plus）
     */
    private String model;

    /**
     * 生成图片数量（默认1）
     */
    @JsonProperty("n")
    private Integer n;

    /**
     * 图片尺寸（如 1024x1024）
     */
    private String size;

    /**
     * 图片质量（standard/hd）
     */
    private String quality;

    /**
     * 响应格式（url/b64_json）
     */
    @JsonProperty("response_format")
    private String responseFormat;

    /**
     * 用户标识
     */
    private String user;

    private static final long serialVersionUID = 1L;
}
