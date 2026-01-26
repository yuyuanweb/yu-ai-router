/**
 * Web 搜索插件
 * 实时联网搜索，获取最新信息
 * 集成 SerpApi（支持 Google 搜索）
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.plugin.impl;

import cn.hutool.core.util.StrUtil;
import cn.hutool.json.JSONArray;
import cn.hutool.json.JSONObject;
import cn.hutool.json.JSONUtil;
import com.yupi.airouter.plugin.Plugin;
import com.yupi.airouter.plugin.PluginContext;
import com.yupi.airouter.plugin.PluginResult;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Slf4j
@Component
public class WebSearchPlugin implements Plugin {

    private static final String PLUGIN_KEY = "web_search";
    private static final String PLUGIN_NAME = "Web搜索";
    private static final String DESCRIPTION = "实时联网搜索，获取最新信息（SerpApi）";

    /**
     * SerpApi API 端点
     */
    private static final String SERPAPI_ENDPOINT = "https://serpapi.com/search.json";

    /**
     * SerpApi API Key（从配置文件读取）
     */
    @Value("${plugin.serpapi.api-key:}")
    private String apiKey;

    /**
     * 最大搜索结果数量
     */
    private int maxResults = 5;

    /**
     * 搜索引擎类型：google/bing
     */
    private String searchEngine = "google";

    /**
     * 请求超时时间（毫秒）
     */
    private int timeout = 15000;

    /**
     * HTTP 客户端（支持重定向）
     */
    private final HttpClient httpClient = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(15))
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
                this.maxResults = configJson.getInt("maxResults", 5);
                this.searchEngine = configJson.getStr("searchEngine", "google");
                this.timeout = configJson.getInt("timeout", 15000);
            } catch (Exception e) {
                log.warn("解析Web搜索插件配置失败，使用默认配置", e);
            }
        }
        log.info("Web搜索插件初始化: maxResults={}, searchEngine={}, apiKey已配置={}",
                maxResults, searchEngine, StrUtil.isNotBlank(apiKey));
    }

    @Override
    public boolean supports(PluginContext context) {
        // Web 搜索需要有输入内容
        return StrUtil.isNotBlank(context.getInput());
    }

    @Override
    public PluginResult execute(PluginContext context) {
        String query = context.getInput();
        if (StrUtil.isBlank(query)) {
            return PluginResult.fail("搜索关键词不能为空");
        }

        // 检查 API Key
        if (StrUtil.isBlank(apiKey)) {
            return PluginResult.fail("SerpApi API Key 未配置，请在插件配置中设置 apiKey。注册地址：https://serpapi.com/");
        }

        log.info("执行Web搜索: {}", query);

        try {
            // 使用 SerpApi 进行搜索
            List<SearchResult> results = searchWithSerpApi(query);

            if (results.isEmpty()) {
                return PluginResult.fail("未找到相关搜索结果");
            }

            // 构建搜索结果文本
            StringBuilder contentBuilder = new StringBuilder();
            contentBuilder.append("## 搜索结果\n\n");
            contentBuilder.append("以下是关于「").append(query).append("」的搜索结果：\n\n");

            for (int i = 0; i < results.size(); i++) {
                SearchResult result = results.get(i);
                contentBuilder.append("### ").append(i + 1).append(". ").append(result.title).append("\n");
                contentBuilder.append(result.snippet).append("\n");
                contentBuilder.append("来源: ").append(result.url).append("\n\n");
            }

            // 构建额外数据
            Map<String, Object> data = new HashMap<>();
            data.put("query", query);
            data.put("resultCount", results.size());
            data.put("searchEngine", searchEngine);
            data.put("results", results);

            return PluginResult.success(contentBuilder.toString(), data);
        } catch (Exception e) {
            log.error("Web搜索失败: {}", query, e);
            return PluginResult.fail("搜索失败: " + e.getMessage());
        }
    }

    /**
     * 使用 SerpApi HTTP API 进行搜索
     * SerpApi 文档：https://serpapi.com/search-api
     */
    private List<SearchResult> searchWithSerpApi(String query) throws Exception {
        List<SearchResult> results = new ArrayList<>();

        // 构建请求 URL
        String encodedQuery = URLEncoder.encode(query, StandardCharsets.UTF_8);
        StringBuilder urlBuilder = new StringBuilder(SERPAPI_ENDPOINT);
        urlBuilder.append("?api_key=").append(apiKey);
        urlBuilder.append("&q=").append(encodedQuery);
        urlBuilder.append("&engine=").append(searchEngine);
        urlBuilder.append("&num=").append(maxResults);
        urlBuilder.append("&hl=zh-CN");  // 搜索语言
        urlBuilder.append("&gl=cn");      // 搜索地区

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(urlBuilder.toString()))
                .timeout(Duration.ofMillis(timeout))
                .GET()
                .build();

        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());

        if (response.statusCode() != 200) {
            log.error("SerpApi 调用失败，状态码: {}, 响应: {}", response.statusCode(), response.body());
            throw new RuntimeException("SerpApi 调用失败，状态码: " + response.statusCode());
        }

        // 解析响应
        JSONObject responseJson = JSONUtil.parseObj(response.body());

        // 检查错误
        if (responseJson.containsKey("error")) {
            String error = responseJson.getStr("error");
            log.error("SerpApi 返回错误: {}", error);
            throw new RuntimeException("SerpApi 错误: " + error);
        }

        // 解析搜索结果
        JSONArray organicResults = responseJson.getJSONArray("organic_results");
        if (organicResults != null) {
            for (int i = 0; i < organicResults.size() && i < maxResults; i++) {
                JSONObject item = organicResults.getJSONObject(i);

                String title = item.getStr("title", "");
                String url = item.getStr("link", "");
                String snippet = item.getStr("snippet", "");

                if (StrUtil.isNotBlank(title) && StrUtil.isNotBlank(url)) {
                    results.add(new SearchResult(title, url, snippet));
                }
            }
        }

        log.info("SerpApi搜索完成，找到 {} 条结果", results.size());
        return results;
    }

    /**
     * 搜索结果数据类
     */
    public record SearchResult(String title, String url, String snippet) {
    }
}
