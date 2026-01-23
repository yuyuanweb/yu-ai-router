package com.yupi.airouter.model.dto.image;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;
import java.util.List;

/**
 * 图片生成响应
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ImageGenerationResponse implements Serializable {

    /**
     * 创建时间戳
     */
    private Long created;

    /**
     * 生成的图片列表
     */
    private List<ImageData> data;

    /**
     * 图片数据
     */
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class ImageData implements Serializable {
        /**
         * 图片URL
         */
        private String url;

        /**
         * Base64编码的图片数据
         */
        private String b64Json;

        /**
         * 修订提示词（如果模型返回）
         */
        private String revisedPrompt;
    }

    private static final long serialVersionUID = 1L;
}
