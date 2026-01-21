/**
 * HTTP 客户端
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.ai.sdk.http;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.yupi.ai.sdk.callback.StreamCallback;
import com.yupi.ai.sdk.config.ClientConfig;
import com.yupi.ai.sdk.exception.AuthException;
import com.yupi.ai.sdk.exception.RateLimitException;
import com.yupi.ai.sdk.exception.YuAIException;
import com.yupi.ai.sdk.model.ChatChunk;
import com.yupi.ai.sdk.model.ChatRequest;
import com.yupi.ai.sdk.model.ChatResponse;
import com.yupi.ai.sdk.model.StreamResponse;
import lombok.extern.slf4j.Slf4j;
import okhttp3.*;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.util.concurrent.TimeUnit;

/**
 * HTTP 客户端封装
 */
@Slf4j
public class HttpClient {

    private final ClientConfig config;
    private final OkHttpClient okHttpClient;
    private final Gson gson;
    private static final MediaType JSON = MediaType.get("application/json; charset=utf-8");

    public HttpClient(ClientConfig config) {
        this.config = config;
        this.gson = new GsonBuilder().create();
        this.okHttpClient = new OkHttpClient.Builder()
                .connectTimeout(config.getConnectTimeout(), TimeUnit.MILLISECONDS)
                .readTimeout(config.getReadTimeout(), TimeUnit.MILLISECONDS)
                .writeTimeout(config.getWriteTimeout(), TimeUnit.MILLISECONDS)
                .retryOnConnectionFailure(true)
                .build();
    }

    /**
     * 同步聊天请求
     */
    public ChatResponse chat(ChatRequest request) {
        request.setStream(false);
        String url = config.getBaseUrl() + "/v1/chat/completions";
        String jsonBody = gson.toJson(request);

        Request httpRequest = buildRequest(url, jsonBody);

        int retryCount = 0;
        Exception lastException = null;

        while (retryCount <= config.getMaxRetries()) {
            try (Response response = okHttpClient.newCall(httpRequest).execute()) {
                return handleResponse(response);
            } catch (RateLimitException | AuthException e) {
                // 认证和限流异常不重试
                throw e;
            } catch (Exception e) {
                lastException = e;
                retryCount++;
                if (retryCount <= config.getMaxRetries()) {
                    log.warn("Request failed, retrying... ({}/{})", retryCount, config.getMaxRetries(), e);
                    try {
                        Thread.sleep(config.getRetryDelay() * retryCount);
                    } catch (InterruptedException ie) {
                        Thread.currentThread().interrupt();
                        throw new YuAIException("Request interrupted", ie);
                    }
                }
            }
        }

        throw new YuAIException("Request failed after " + config.getMaxRetries() + " retries", lastException);
    }

    /**
     * 流式聊天请求（解析 OpenAI SSE 格式）
     */
    public void chatStream(ChatRequest request, StreamCallback callback) {
        request.setStream(true);
        String url = config.getBaseUrl() + "/v1/chat/completions";
        String jsonBody = gson.toJson(request);

        Request httpRequest = buildRequest(url, jsonBody);

        try {
            Response response = okHttpClient.newCall(httpRequest).execute();
            if (!response.isSuccessful()) {
                handleErrorResponse(response);
            }

            ResponseBody body = response.body();
            if (body == null) {
                throw new YuAIException("Response body is null");
            }

            InputStream inputStream = body.byteStream();
            BufferedReader reader = new BufferedReader(new InputStreamReader(inputStream));

            String line;
            while ((line = reader.readLine()) != null) {
                // 跳过空行
                if (line.trim().isEmpty()) {
                    continue;
                }

                if (line.startsWith("data:")) {
                    // 去掉 "data: " 前缀
                    String data = line.substring(5);
                    if (data.trim().isEmpty()) {
                        continue;
                    }

                    try {
                        // 解析 JSON 为 StreamResponse
                        StreamResponse streamResponse = gson.fromJson(data, StreamResponse.class);

                        // 提取内容并转换为 ChatChunk
                        if (streamResponse.getChoices() != null && !streamResponse.getChoices().isEmpty()) {
                            StreamResponse.StreamChoice choice = streamResponse.getChoices().getFirst();

                            // 检查是否结束（finishReason 为 "stop"）
                            if ("stop".equals(choice.getFinishReason())) {
                                callback.onComplete();
                                break;
                            }

                            if (choice.getDelta() != null) {
                                StreamResponse.Delta delta = choice.getDelta();

                                // 如果既没有内容也没有思考内容，跳过
                                if (delta.getContent() == null && delta.getReasoningContent() == null) {
                                    continue;
                                }

                                ChatChunk chunk = ChatChunk.builder()
                                        .content(delta.getContent())
                                        .reasoningContent(delta.getReasoningContent())
                                        .model(streamResponse.getModel())
                                        .done(false)
                                        .build();
                                callback.onMessage(chunk);
                            }
                        }
                    } catch (Exception e) {
                        log.error("Failed to parse stream response: {}", data, e);
                    }
                }
            }

        } catch (Exception e) {
            log.error("Stream request failed", e);
            callback.onError(e);
        }
    }

    /**
     * 构建 HTTP 请求
     */
    private Request buildRequest(String url, String jsonBody) {
        RequestBody body = RequestBody.create(jsonBody, JSON);
        return new Request.Builder()
                .url(url)
                .header("Authorization", "Bearer " + config.getApiKey())
                .header("Content-Type", "application/json")
                .post(body)
                .build();
    }

    /**
     * 处理响应
     */
    private ChatResponse handleResponse(Response response) throws IOException {
        if (!response.isSuccessful()) {
            handleErrorResponse(response);
        }

        ResponseBody body = response.body();
        if (body == null) {
            throw new YuAIException("Response body is null");
        }

        String responseBody = body.string();
        return gson.fromJson(responseBody, ChatResponse.class);
    }

    /**
     * 处理错误响应
     */
    private void handleErrorResponse(Response response) throws IOException {
        int code = response.code();
        String message = "Request failed with code: " + code;

        ResponseBody body = response.body();
        if (body != null) {
            message = body.string();
        }

        switch (code) {
            case 401:
                throw new AuthException(message);
            case 429:
                throw new RateLimitException(message);
            default:
                throw new YuAIException(code, message);
        }
    }

    /**
     * 关闭客户端
     */
    public void close() {
        okHttpClient.dispatcher().executorService().shutdown();
        okHttpClient.connectionPool().evictAll();
    }
}
