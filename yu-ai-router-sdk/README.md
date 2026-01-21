# Yu AI Router SDK

Yu AI Router 的 Java SDK，提供简单易用的 API 调用接口。

## 要求

- Java 21+
- Maven 3.6+

## 特性

- ✅ **零依赖核心** - 仅依赖 OkHttp 和 Gson
- ✅ **类型安全** - 完整的 Java 类型定义
- ✅ **同步调用** - 简单直接的同步 API
- ✅ **流式支持** - 支持 SSE 流式响应
- ✅ **自动重试** - 内置指数退避重试机制
- ✅ **异常处理** - 细分的异常类型
- ✅ **超时控制** - 灵活的超时配置

## 快速开始

### Maven 依赖

```xml
<dependency>
    <groupId>com.yupi.ai</groupId>
    <artifactId>yu-ai-router-sdk</artifactId>
    <version>1.0.0</version>
</dependency>
```

### 基础使用

```java
// 1. 创建客户端
YuAIClient client = YuAIClient.builder()
    .apiKey("sk-your-api-key")
    .baseUrl("http://localhost:8123/api")  // 可选，默认本地地址
    .build();

// 2. 同步调用
ChatResponse response = client.chat("你好，请介绍一下自己");
System.out.println(response.getContent());

// 3. 关闭客户端
client.close();
```

## 完整示例

### 1. 简单对话

```java
YuAIClient client = YuAIClient.builder()
    .apiKey("sk-xxx")
    .build();

// 最简单的调用
ChatResponse response = client.chat("讲个笑话");
System.out.println(response.getContent());
```

### 2. 指定模型

```java
// 使用特定模型
ChatResponse response = client.chat("qwen-plus", "解释一下量子计算");
System.out.println(response.getContent());
```

### 3. 高级配置

```java
// 构建完整的请求
ChatRequest request = ChatRequest.builder()
    .model("qwen-plus")
    .messages(Arrays.asList(
        ChatMessage.system("你是一个专业的编程助手"),
        ChatMessage.user("如何学习 Java？")
    ))
    .temperature(0.7)
    .maxTokens(1000)
    .build();

ChatResponse response = client.chat(request);
System.out.println(response.getContent());
```

### 4. 流式调用

```java
client.chatStream("写一篇关于人工智能的文章", new StreamCallback() {
    @Override
    public void onMessage(ChatChunk chunk) {
        System.out.print(chunk.getContent());
    }
    
    @Override
    public void onComplete() {
        System.out.println("\n\n[完成]");
    }
    
    @Override
    public void onError(Throwable error) {
        System.err.println("错误: " + error.getMessage());
    }
});
```

### 5. 自定义配置

```java
YuAIClient client = YuAIClient.builder()
    .apiKey("sk-xxx")
    .baseUrl("https://api.yu-ai.com")
    .connectTimeout(15000)        // 连接超时 15秒
    .readTimeout(60000)           // 读取超时 60秒
    .maxRetries(5)                // 最多重试 5 次
    .retryDelay(2000)             // 重试延迟 2秒
    .build();
```

## API 文档

### YuAIClient

主要客户端类，提供所有 API 调用方法。

#### 创建客户端

```java
YuAIClient client = YuAIClient.builder()
    .apiKey(String apiKey)           // 必需：API Key
    .baseUrl(String baseUrl)         // 可选：基础 URL
    .connectTimeout(Integer timeout) // 可选：连接超时
    .readTimeout(Integer timeout)    // 可选：读取超时
    .writeTimeout(Integer timeout)   // 可选：写入超时
    .maxRetries(Integer retries)     // 可选：最大重试次数
    .retryDelay(Integer delay)       // 可选：重试延迟
    .build();
```

#### 同步调用

```java
// 简单文本
ChatResponse chat(String message)

// 指定模型
ChatResponse chat(String model, String message)

// 完整请求
ChatResponse chat(ChatRequest request)
```

#### 流式调用

```java
// 简单文本
void chatStream(String message, StreamCallback callback)

// 指定模型
void chatStream(String model, String message, StreamCallback callback)

// 完整请求
void chatStream(ChatRequest request, StreamCallback callback)
```

### ChatRequest

聊天请求对象。

```java
ChatRequest request = ChatRequest.builder()
    .model(String model)                      // 模型名称
    .messages(List<ChatMessage> messages)     // 消息列表
    .temperature(Double temperature)          // 温度参数 0-1
    .maxTokens(Integer maxTokens)            // 最大 Token 数
    .enableReasoning(Boolean enable)         // 启用深度思考
    .routingStrategy(String strategy)        // 路由策略
    .build();
```

### ChatResponse

聊天响应对象。

```java
String id             // 响应 ID
String model          // 使用的模型
List<Choice> choices  // 选项列表
Usage usage          // Token 使用统计

// 便捷方法
String getContent()  // 获取响应内容
```

### StreamCallback

流式响应回调接口。

```java
public interface StreamCallback {
    void onMessage(ChatChunk chunk);  // 接收消息块
    void onComplete();                 // 完成
    void onError(Throwable error);    // 错误处理
}
```

## 异常处理

SDK 提供了细分的异常类型：

```java
try {
    ChatResponse response = client.chat("你好");
} catch (AuthException e) {
    // API Key 无效或过期
    System.err.println("认证失败: " + e.getMessage());
} catch (RateLimitException e) {
    // 请求过于频繁
    System.err.println("限流: " + e.getMessage());
} catch (YuAIException e) {
    // 其他错误
    System.err.println("调用失败: " + e.getMessage());
}
```

## 最佳实践

1. **复用客户端** - 创建一次，多次使用
2. **及时关闭** - 使用完毕后调用 `client.close()`
3. **异常处理** - 捕获并处理不同类型的异常
4. **超时设置** - 根据实际需求调整超时时间
5. **流式响应** - 大文本生成使用流式调用

## 示例代码

完整的示例代码请参考 `src/test/java/com/yupi/ai/sdk/example/` 目录。

## License

MIT License

## 支持

- GitHub: https://github.com/liyupi/yu-ai-router
- 文档: https://docs.yu-ai.com
- 问题反馈: https://github.com/liyupi/yu-ai-router/issues
