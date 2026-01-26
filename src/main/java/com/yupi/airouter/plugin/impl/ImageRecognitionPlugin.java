/**
 * 图片识别插件
 * 识别图片内容，返回图片描述
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.plugin.impl;

import cn.hutool.core.util.StrUtil;
import cn.hutool.json.JSONObject;
import cn.hutool.json.JSONUtil;
import com.mybatisflex.core.query.QueryWrapper;
import com.yupi.airouter.model.entity.ModelProvider;
import com.yupi.airouter.plugin.Plugin;
import com.yupi.airouter.plugin.PluginContext;
import com.yupi.airouter.plugin.PluginResult;
import com.yupi.airouter.service.ModelProviderService;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.util.Base64;
import java.util.HashMap;
import java.util.Map;

@Slf4j
@Component
public class ImageRecognitionPlugin implements Plugin {

    private static final String PLUGIN_KEY = "image_recognition";
    private static final String PLUGIN_NAME = "图片识别";
    private static final String DESCRIPTION = "识别图片内容，返回图片描述";

    /**
     * 使用的视觉模型
     */
    private String model = "qwen-vl-plus";

    /**
     * 最大图片大小（字节）
     */
    private int maxImageSize = 4 * 1024 * 1024; // 4MB

    /**
     * 默认提示词
     */
    private String defaultPrompt = "请描述这张图片的内容，包括主要对象、场景、颜色、文字等信息。";

    @Resource
    private ModelProviderService modelProviderService;

    /**
     * HTTP 客户端（支持重定向）
     */
    private final HttpClient httpClient = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(30))
            .followRedirects(HttpClient.Redirect.NORMAL)
            .build();

    @Override
    public String getPluginKey() {
        return PLUGIN_KEY;
    }

    @Override
    public String getPluginName() {
        return PLUGIN_NAME;
    }

    @Override
    public String getDescription() {
        return DESCRIPTION;
    }

    @Override
    public void init(String config) {
        if (StrUtil.isNotBlank(config)) {
            try {
                JSONObject configJson = JSONUtil.parseObj(config);
                this.model = configJson.getStr("model", "qwen-vl-plus");
                this.maxImageSize = configJson.getInt("maxImageSize", 4 * 1024 * 1024);
                this.defaultPrompt = configJson.getStr("defaultPrompt", defaultPrompt);
                log.info("图片识别插件初始化: model={}, maxImageSize={}", model, maxImageSize);
            } catch (Exception e) {
                log.warn("解析图片识别插件配置失败，使用默认配置", e);
            }
        }
    }

    @Override
    public boolean supports(PluginContext context) {
        // 需要有图片 URL 或文件字节数组
        return StrUtil.isNotBlank(context.getFileUrl()) 
                || (context.getFileBytes() != null && context.getFileBytes().length > 0);
    }

    @Override
    public PluginResult execute(PluginContext context) {
        log.info("执行图片识别");

        try {
            // 获取图片数据
            byte[] imageBytes;
            String mimeType = "image/jpeg";

            // 获取图片数据
            if (context.getFileBytes() != null && context.getFileBytes().length > 0) {
                // 直接使用上传的文件字节数组
                imageBytes = context.getFileBytes();
                if (StrUtil.isNotBlank(context.getFileType())) {
                    mimeType = context.getFileType();
                }
                log.info("使用上传的图片，大小: {} bytes", imageBytes.length);
            } else if (StrUtil.isNotBlank(context.getFileUrl())) {
                // 从 URL 下载
                imageBytes = downloadImage(context.getFileUrl());
                mimeType = guessMimeType(context.getFileUrl());
                log.info("从URL下载图片，大小: {} bytes", imageBytes.length);
            } else {
                return PluginResult.fail("请提供图片文件或图片URL");
            }

            // 检查图片大小
            if (imageBytes.length > maxImageSize) {
                return PluginResult.fail("图片大小超过限制（最大 " + (maxImageSize / 1024 / 1024) + "MB）");
            }

            // 获取用户问题或使用默认提示词
            String prompt = StrUtil.isNotBlank(context.getInput()) ? context.getInput() : defaultPrompt;

            // 调用视觉模型
            String result = recognizeImage(imageBytes, mimeType, prompt);

            // 构建结果
            StringBuilder contentBuilder = new StringBuilder();
            contentBuilder.append("## 图片识别结果\n\n");
            if (StrUtil.isNotBlank(context.getInput())) {
                contentBuilder.append("用户问题: ").append(context.getInput()).append("\n\n");
            }
            contentBuilder.append(result);

            // 构建额外数据
            Map<String, Object> data = new HashMap<>();
            data.put("model", model);
            data.put("imageSize", imageBytes.length);

            return PluginResult.success(contentBuilder.toString(), data);
        } catch (Exception e) {
            log.error("图片识别失败", e);
            return PluginResult.fail("图片识别失败: " + e.getMessage());
        }
    }

    /**
     * 下载图片
     */
    private byte[] downloadImage(String url) throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .timeout(Duration.ofSeconds(30))
                .GET()
                .build();

        HttpResponse<byte[]> response = httpClient.send(request, HttpResponse.BodyHandlers.ofByteArray());

        if (response.statusCode() != 200) {
            throw new RuntimeException("下载图片失败，状态码: " + response.statusCode());
        }

        return response.body();
    }

    /**
     * 猜测图片 MIME 类型
     */
    private String guessMimeType(String url) {
        url = url.toLowerCase();
        if (url.contains(".png")) {
            return "image/png";
        } else if (url.contains(".gif")) {
            return "image/gif";
        } else if (url.contains(".webp")) {
            return "image/webp";
        }
        return "image/jpeg";
    }

    /**
     * 调用视觉模型识别图片
     * 使用 HTTP API 直接调用通义千问视觉模型
     */
    private String recognizeImage(byte[] imageBytes, String mimeType, String prompt) {
        // 获取通义千问提供者
        ModelProvider provider = modelProviderService.getOne(
                QueryWrapper.create().eq("providerName", "qwen"));
        if (provider == null) {
            throw new RuntimeException("未找到通义千问提供者配置");
        }

        try {
            // 将图片转换为 Base64
            String base64Image = Base64.getEncoder().encodeToString(imageBytes);
            String imageUrl = "data:" + mimeType + ";base64," + base64Image;

            // 构建请求体（通义千问视觉模型 API 格式）
            JSONObject requestBody = new JSONObject();
            requestBody.set("model", model);
            
            // 构建消息
            JSONObject userMessage = new JSONObject();
            userMessage.set("role", "user");
            
            // 构建多模态内容
            cn.hutool.json.JSONArray contentArray = new cn.hutool.json.JSONArray();
            
            // 添加图片
            JSONObject imageContent = new JSONObject();
            imageContent.set("type", "image_url");
            JSONObject imageUrlObj = new JSONObject();
            imageUrlObj.set("url", imageUrl);
            imageContent.set("image_url", imageUrlObj);
            contentArray.add(imageContent);
            
            // 添加文本
            JSONObject textContent = new JSONObject();
            textContent.set("type", "text");
            textContent.set("text", prompt);
            contentArray.add(textContent);
            
            userMessage.set("content", contentArray);
            
            cn.hutool.json.JSONArray messagesArray = new cn.hutool.json.JSONArray();
            messagesArray.add(userMessage);
            requestBody.set("messages", messagesArray);

            // 发送请求
            HttpRequest request = HttpRequest.newBuilder()
                    .uri(URI.create(provider.getBaseUrl() + "/v1/chat/completions"))
                    .header("Content-Type", "application/json")
                    .header("Authorization", "Bearer " + provider.getApiKey())
                    .timeout(Duration.ofSeconds(60))
                    .POST(HttpRequest.BodyPublishers.ofString(requestBody.toString()))
                    .build();

            HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());

            if (response.statusCode() != 200) {
                log.error("视觉模型调用失败，状态码: {}, 响应: {}", response.statusCode(), response.body());
                throw new RuntimeException("视觉模型调用失败，状态码: " + response.statusCode());
            }

            // 解析响应
            JSONObject responseJson = JSONUtil.parseObj(response.body());
            cn.hutool.json.JSONArray choices = responseJson.getJSONArray("choices");
            if (choices != null && !choices.isEmpty()) {
                JSONObject choice = choices.getJSONObject(0);
                JSONObject message = choice.getJSONObject("message");
                if (message != null) {
                    return message.getStr("content", "无法识别图片内容");
                }
            }

            throw new RuntimeException("模型返回结果格式错误");
        } catch (Exception e) {
            log.error("调用视觉模型失败", e);
            throw new RuntimeException("调用视觉模型失败: " + e.getMessage());
        }
    }
}
