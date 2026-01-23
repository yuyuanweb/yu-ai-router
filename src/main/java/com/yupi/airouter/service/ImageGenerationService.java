package com.yupi.airouter.service;

import com.mybatisflex.core.paginate.Page;
import com.mybatisflex.core.service.IService;
import com.yupi.airouter.model.dto.image.ImageGenerationRequest;
import com.yupi.airouter.model.dto.image.ImageGenerationResponse;
import com.yupi.airouter.model.entity.ImageGenerationRecord;

/**
 * 图片生成服务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface ImageGenerationService extends IService<ImageGenerationRecord> {

    /**
     * 生成图片
     *
     * @param request 图片生成请求
     * @param userId 用户ID
     * @param apiKeyId API Key ID
     * @param clientIp 客户端IP
     * @return 图片生成响应
     */
    ImageGenerationResponse generateImage(ImageGenerationRequest request, Long userId, Long apiKeyId, String clientIp);

    /**
     * 获取用户的图片生成记录（分页）
     *
     * @param userId 用户ID
     * @param pageNum 页码
     * @param pageSize 每页大小
     * @return 分页结果
     */
    Page<ImageGenerationRecord> listUserRecords(Long userId, int pageNum, int pageSize);
}
