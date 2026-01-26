/**
 * PDF 解析插件
 * 解析 PDF 文档内容，提取文本信息
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.plugin.impl;

import cn.hutool.core.util.StrUtil;
import cn.hutool.json.JSONObject;
import cn.hutool.json.JSONUtil;
import com.yupi.airouter.plugin.Plugin;
import com.yupi.airouter.plugin.PluginContext;
import com.yupi.airouter.plugin.PluginResult;
import lombok.extern.slf4j.Slf4j;
import org.apache.pdfbox.Loader;
import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.text.PDFTextStripper;
import org.springframework.stereotype.Component;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;

@Slf4j
@Component
public class PdfParserPlugin implements Plugin {

    private static final String PLUGIN_KEY = "pdf_parser";
    private static final String PLUGIN_NAME = "PDF解析";
    private static final String DESCRIPTION = "解析PDF文档内容，提取文本信息";

    /**
     * 最大解析页数
     */
    private int maxPages = 50;

    /**
     * 最大文本长度（字符）
     */
    private int maxTextLength = 50000;

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
                this.maxPages = configJson.getInt("maxPages", 50);
                this.maxTextLength = configJson.getInt("maxTextLength", 50000);
                log.info("PDF解析插件初始化: maxPages={}, maxTextLength={}", maxPages, maxTextLength);
            } catch (Exception e) {
                log.warn("解析PDF解析插件配置失败，使用默认配置", e);
            }
        }
    }

    @Override
    public boolean supports(PluginContext context) {
        // 需要有文件 URL 或文件字节数组
        return StrUtil.isNotBlank(context.getFileUrl()) 
                || (context.getFileBytes() != null && context.getFileBytes().length > 0);
    }

    @Override
    public PluginResult execute(PluginContext context) {
        log.info("执行PDF解析");

        try {
            byte[] pdfBytes;

            // 获取 PDF 数据
            if (context.getFileBytes() != null && context.getFileBytes().length > 0) {
                // 直接使用上传的文件字节数组
                pdfBytes = context.getFileBytes();
                log.info("使用上传的文件，大小: {} bytes", pdfBytes.length);
            } else if (StrUtil.isNotBlank(context.getFileUrl())) {
                // 从 URL 下载
                pdfBytes = downloadFile(context.getFileUrl());
                log.info("从URL下载PDF，大小: {} bytes", pdfBytes.length);
            } else {
                return PluginResult.fail("请提供PDF文件或文件URL");
            }

            // 解析 PDF
            String text = parsePdf(pdfBytes);

            if (StrUtil.isBlank(text)) {
                return PluginResult.fail("PDF文档为空或无法提取文本");
            }

            // 构建结果
            StringBuilder contentBuilder = new StringBuilder();
            contentBuilder.append("## PDF文档内容\n\n");
            
            // 如果有用户问题，添加说明
            if (StrUtil.isNotBlank(context.getInput())) {
                contentBuilder.append("用户问题: ").append(context.getInput()).append("\n\n");
                contentBuilder.append("以下是PDF文档的内容，请根据内容回答用户问题：\n\n");
            }
            
            contentBuilder.append("```\n");
            contentBuilder.append(text);
            contentBuilder.append("\n```");

            // 构建额外数据
            Map<String, Object> data = new HashMap<>();
            data.put("textLength", text.length());
            data.put("truncated", text.length() >= maxTextLength);

            return PluginResult.success(contentBuilder.toString(), data);
        } catch (Exception e) {
            log.error("PDF解析失败", e);
            return PluginResult.fail("PDF解析失败: " + e.getMessage());
        }
    }

    /**
     * 下载文件
     */
    private byte[] downloadFile(String url) throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .timeout(Duration.ofSeconds(60))
                .GET()
                .build();

        HttpResponse<byte[]> response = httpClient.send(request, HttpResponse.BodyHandlers.ofByteArray());

        if (response.statusCode() != 200) {
            throw new RuntimeException("下载文件失败，状态码: " + response.statusCode());
        }

        return response.body();
    }

    /**
     * 解析 PDF 文档
     */
    private String parsePdf(byte[] pdfBytes) throws Exception {
        try (InputStream inputStream = new ByteArrayInputStream(pdfBytes);
             PDDocument document = Loader.loadPDF(pdfBytes)) {

            int totalPages = document.getNumberOfPages();
            log.info("PDF总页数: {}", totalPages);

            PDFTextStripper stripper = new PDFTextStripper();

            // 限制解析页数
            int endPage = Math.min(totalPages, maxPages);
            stripper.setStartPage(1);
            stripper.setEndPage(endPage);

            String text = stripper.getText(document);

            // 限制文本长度
            if (text.length() > maxTextLength) {
                text = text.substring(0, maxTextLength) + "\n\n[文档内容过长，已截断...]";
                log.info("PDF文本已截断，原长度: {}", text.length());
            }

            // 清理文本（移除多余空白）
            text = text.replaceAll("\\n{3,}", "\n\n").trim();

            return text;
        }
    }
}
