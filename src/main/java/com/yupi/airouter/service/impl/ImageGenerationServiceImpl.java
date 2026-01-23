package com.yupi.airouter.service.impl;

import cn.hutool.core.util.StrUtil;
import cn.hutool.http.HttpRequest;
import cn.hutool.http.HttpResponse;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.mybatisflex.core.paginate.Page;
import com.mybatisflex.core.query.QueryWrapper;
import com.mybatisflex.spring.service.impl.ServiceImpl;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.mapper.ImageGenerationRecordMapper;
import com.yupi.airouter.model.dto.image.ImageGenerationRequest;
import com.yupi.airouter.model.dto.image.ImageGenerationResponse;
import com.yupi.airouter.model.entity.ImageGenerationRecord;
import com.yupi.airouter.model.entity.Model;
import com.yupi.airouter.model.entity.ModelProvider;
import com.yupi.airouter.service.*;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * 图片生成服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Service
@Slf4j
public class ImageGenerationServiceImpl extends ServiceImpl<ImageGenerationRecordMapper, ImageGenerationRecord> implements ImageGenerationService {

    @Resource
    private ModelService modelService;

    @Resource
    private ModelProviderService modelProviderService;

    @Resource
    private QuotaService quotaService;

    @Resource
    private BalanceService balanceService;

    @Resource
    private BillingService billingService;

    @Resource
    private UserService userService;

    @Resource
    private ObjectMapper objectMapper;

    private static final String DEFAULT_MODEL = "qwen-image-plus";
    private static final String DEFAULT_SIZE = "1024*1024";
    private static final int DEFAULT_N = 1;

    @Override
    @Transactional(rollbackFor = Exception.class)
    public ImageGenerationResponse generateImage(ImageGenerationRequest request, Long userId, Long apiKeyId, String clientIp) {
        long startTime = System.currentTimeMillis();

        // 参数校验
        if (StrUtil.isBlank(request.getPrompt())) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "提示词不能为空");
        }

        // 检查用户状态
        if (userId != null && userService.isUserDisabled(userId)) {
            throw new BusinessException(ErrorCode.FORBIDDEN_ERROR, "账号已被禁用，无法使用服务");
        }

        // 检查用户配额
        if (userId != null && !quotaService.checkQuota(userId)) {
            throw new BusinessException(ErrorCode.OPERATION_ERROR, "Token配额已用尽，请联系管理员增加配额");
        }

        // 设置默认值
        String modelKey = StrUtil.isNotBlank(request.getModel()) ? request.getModel() : DEFAULT_MODEL;
        String size = StrUtil.isNotBlank(request.getSize()) ? request.getSize() : DEFAULT_SIZE;
        int n = 1;

        // 查询模型信息
        Model model = modelService.getByModelKey(modelKey);
        if (model == null || !"image".equals(model.getModelType())) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "模型不存在或不是绘图模型");
        }

        // 查询提供者信息
        ModelProvider provider = modelProviderService.getById(model.getProviderId());
        if (provider == null) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "模型提供者不存在");
        }

        // 预估费用（绘图模型按张计费）
        BigDecimal estimatedCost = model.getInputPrice() != null
                ? model.getInputPrice().multiply(new BigDecimal(n))
                : BigDecimal.ZERO;

        // 检查余额
        if (userId != null && !balanceService.checkBalance(userId, estimatedCost)) {
            throw new BusinessException(ErrorCode.OPERATION_ERROR,
                    "账户余额不足，生成" + n + "张图片预计需要¥" + estimatedCost + "，请先充值");
        }

        try {
            // 调用模型生成图片
            ImageGenerationResponse response = callImageGenerationModel(model, provider, request, size, n);

            long duration = System.currentTimeMillis() - startTime;

            // 计算实际费用（按实际生成的图片数量）
            int actualImageCount = response.getData() != null ? response.getData().size() : n;
            BigDecimal actualCost = model.getInputPrice() != null
                    ? model.getInputPrice().multiply(new BigDecimal(actualImageCount))
                    : BigDecimal.ZERO;

            // 记录生成成功
            for (ImageGenerationResponse.ImageData imageData : response.getData()) {
                ImageGenerationRecord record = ImageGenerationRecord.builder()
                        .userId(userId)
                        .apiKeyId(apiKeyId)
                        .modelId(model.getId())
                        .modelKey(modelKey)
                        .prompt(request.getPrompt())
                        .revisedPrompt(imageData.getRevisedPrompt())
                        .imageUrl(imageData.getUrl())
                        .imageData(imageData.getB64Json())
                        .size(size)
                        .quality(request.getQuality())
                        .status("success")
                        .cost(model.getInputPrice() != null ? model.getInputPrice() : BigDecimal.ZERO)
                        .duration((int) duration)
                        .clientIp(clientIp)
                        .createTime(LocalDateTime.now())
                        .build();
                save(record);
            }

            // 扣减配额和余额
            if (userId != null) {
                // 绘图消耗固定Token（假设1000个Token/张）
                int tokensPerImage = 1000;
                int totalTokens = actualImageCount * tokensPerImage;
                quotaService.deductTokens(userId, totalTokens);

                // 扣减余额
                if (actualCost.compareTo(BigDecimal.ZERO) > 0) {
                    String description = apiKeyId != null
                            ? "API图片生成 - " + modelKey + " x" + actualImageCount
                            : "网页图片生成 - " + modelKey + " x" + actualImageCount;
                    balanceService.deductBalance(userId, actualCost, null, description);
                }
            }

            log.info("图片生成成功：用户 {}, 模型 {}, 数量 {}, 耗时 {}ms", userId, modelKey, actualImageCount, duration);
            return response;

        } catch (Exception e) {
            long duration = System.currentTimeMillis() - startTime;

            // 记录生成失败
            ImageGenerationRecord record = ImageGenerationRecord.builder()
                    .userId(userId)
                    .apiKeyId(apiKeyId)
                    .modelId(model.getId())
                    .modelKey(modelKey)
                    .prompt(request.getPrompt())
                    .size(size)
                    .quality(request.getQuality())
                    .status("failed")
                    .cost(BigDecimal.ZERO)
                    .duration((int) duration)
                    .errorMessage(e.getMessage())
                    .clientIp(clientIp)
                    .createTime(LocalDateTime.now())
                    .build();
            save(record);

            log.error("图片生成失败：用户 {}, 模型 {}, 错误：{}", userId, modelKey, e.getMessage(), e);
            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "图片生成失败: " + e.getMessage());
        }
    }

    /**
     * 调用通义万相生成图片（异步模式 + 轮询）
     */
    private ImageGenerationResponse callImageGenerationModel(Model model, ModelProvider provider,
                                                             ImageGenerationRequest request, String size, int n) {
        try {
            // 构建请求JSON
            Map<String, Object> body = new HashMap<>();
            body.put("model", model.getModelKey());
            body.put("input", Map.of("prompt", request.getPrompt()));
            body.put("parameters", Map.of("size", size, "n", n));
            String requestBody = objectMapper.writeValueAsString(body);

            // 调用通义万相API（异步模式）
            String apiUrl = provider.getBaseUrl().replace("/compatible-mode", "") + "/api/v1/services/aigc/text2image/image-synthesis";

            HttpResponse httpResponse = HttpRequest.post(apiUrl)
                    .header("Authorization", "Bearer " + provider.getApiKey())
                    .header("Content-Type", "application/json")
                    .header("X-DashScope-Async", "enable")  // 异步模式
                    .body(requestBody)
                    .timeout(30000)
                    .execute();

            String responseBody = httpResponse.body();
            log.info("通义万相创建任务响应：{}", responseBody);

            // 解析响应，获取任务ID
            JsonNode rootNode = objectMapper.readTree(responseBody);

            if (!rootNode.has("output") || !rootNode.get("output").has("task_id")) {
                String errorMsg = rootNode.has("message") ? rootNode.get("message").asText() : "未返回任务ID";
                throw new BusinessException(ErrorCode.SYSTEM_ERROR, "创建图片生成任务失败：" + errorMsg);
            }

            String taskId = rootNode.get("output").get("task_id").asText();
            log.info("图片生成任务ID：{}", taskId);

            // 轮询任务状态
            String taskApiUrl = provider.getBaseUrl().replace("/compatible-mode", "") + "/api/v1/tasks/" + taskId;

            int maxRetries = 60;  // 最多轮询60次
            int retryCount = 0;

            while (retryCount < maxRetries) {
                Thread.sleep(2000);  // 等待2秒

                HttpResponse taskResponse = HttpRequest.get(taskApiUrl)
                        .header("Authorization", "Bearer " + provider.getApiKey())
                        .timeout(10000)
                        .execute();

                String taskResponseBody = taskResponse.body();
                JsonNode taskNode = objectMapper.readTree(taskResponseBody);

                if (!taskNode.has("output") || !taskNode.get("output").has("task_status")) {
                    retryCount++;
                    continue;
                }

                String taskStatus = taskNode.get("output").get("task_status").asText();
                log.info("任务状态：{} (轮询次数: {})", taskStatus, retryCount + 1);

                if ("SUCCEEDED".equals(taskStatus)) {
                    // 任务成功，提取图片URL
                    List<ImageGenerationResponse.ImageData> imageDataList = new ArrayList<>();

                    if (taskNode.has("output") && taskNode.get("output").has("results")) {
                        JsonNode results = taskNode.get("output").get("results");
                        for (JsonNode result : results) {
                            if (result.has("url")) {
                                ImageGenerationResponse.ImageData imageData = ImageGenerationResponse.ImageData.builder()
                                        .url(result.get("url").asText())
                                        .revisedPrompt(request.getPrompt())
                                        .build();
                                imageDataList.add(imageData);
                            }
                        }
                    }

                    if (imageDataList.isEmpty()) {
                        throw new BusinessException(ErrorCode.SYSTEM_ERROR, "任务成功但未返回图片");
                    }

                    return ImageGenerationResponse.builder()
                            .created(System.currentTimeMillis() / 1000)
                            .data(imageDataList)
                            .build();

                } else if ("FAILED".equals(taskStatus)) {
                    String errorMsg = taskNode.has("output") && taskNode.get("output").has("message")
                            ? taskNode.get("output").get("message").asText()
                            : "任务失败";
                    throw new BusinessException(ErrorCode.SYSTEM_ERROR, "图片生成失败：" + errorMsg);
                }

                // PENDING, RUNNING 状态继续轮询
                retryCount++;
            }

            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "图片生成超时");

        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "图片生成被中断");
        } catch (Exception e) {
            log.error("调用通义万相API失败", e);
            throw new BusinessException(ErrorCode.SYSTEM_ERROR, "调用图片生成API失败: " + e.getMessage());
        }
    }

    @Override
    public Page<ImageGenerationRecord> listUserRecords(Long userId, int pageNum, int pageSize) {
        QueryWrapper queryWrapper = QueryWrapper.create()
                .where("userId = ?", userId)
                .orderBy("createTime", false);

        return page(Page.of(pageNum, pageSize), queryWrapper);
    }
}
