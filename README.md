# Yu-AI-Router 🚀

<p align="center">
  <strong>企业级 AI 网关平台 | Enterprise-Level AI Gateway</strong>
</p>

<p align="center">
  统一接口调用多个主流 AI 模型，支持智能路由、自动故障转移、实时监控
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Spring%20Boot-3.5.9-brightgreen" alt="Spring Boot"/>
  <img src="https://img.shields.io/badge/Spring%20AI-1.1.2-blue" alt="Spring AI"/>
  <img src="https://img.shields.io/badge/Vue-3.5-green" alt="Vue"/>
  <img src="https://img.shields.io/badge/JDK-21+-orange" alt="JDK"/>
  <img src="https://img.shields.io/badge/License-MIT-yellow" alt="License"/>
</p>

---

## 📖 项目简介

Yu-AI-Router 是一个参考 [OpenRouter.ai](https://openrouter.ai/) 设计的**企业级 AI 网关平台**，提供统一的 API 接口来访问多个主流 AI 模型（通义千问、智谱AI、DeepSeek 等），实现智能路由、自动故障转移、完善的监控体系和 Java SDK 支持。

### 🎯 核心价值

| 特性 | 说明 | 价值 |
|------|------|------|
| 🔌 **即插即用** | 兼容 OpenAI SDK 格式，一行代码切换模型 | 零学习成本 |
| 🛡️ **高可用设计** | 智能重试 + 健康检查 + 自动 Fallback | 99.9% 可用性 |
| 🎯 **智能路由** | 支持 auto 模式（成本/速度优先） | 成本降低 30%+ |
| 📊 **全链路监控** | Prometheus + Grafana + TraceId 追踪 | 秒级定位问题 |
| 🚀 **开发者友好** | Java SDK + 完整文档 + 示例代码 | 5 分钟上手 |
| 💰 **消耗透明** | 实时统计 Token 消耗与费用 | 成本可控 |

---

## ✨ 功能特性

### 核心功能

- **🤖 多模型接入**：支持通义千问、智谱AI、DeepSeek 等主流大模型
- **💬 在线对话**：支持流式响应、多轮对话、上下文管理
- **🔑 API Key 管理**：创建、查看、撤销 API Key，支持调用统计
- **📊 Token 统计**：实时统计每次请求的 Token 消耗
- **🖼️ AI 绘图**：支持通义万相、智谱 CogView 等文生图模型
- **💳 在线充值**：Stripe 支付集成，支持余额管理

### 高级特性

- **🔀 智能路由**：成本优先 / 速度优先 / 轮询策略
- **♻️ 自动 Fallback**：模型故障自动切换到备用模型
- **⏱️ 智能重试**：指数退避重试策略
- **🔒 安全防护**：IP 黑名单、请求限流、TraceId 追踪
- **🔌 插件系统**：Web 搜索、PDF 解析、图片识别
- **🗝️ BYOK 支持**：用户自带 API Key，直连模型

### 监控告警

- **📈 实时监控**：QPS、成功率、错误率、延迟分布
- **📉 Prometheus 指标**：完整的业务指标导出
- **📊 Grafana 大盘**：预置监控仪表盘

---

## 🛠️ 技术栈

### 后端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| Spring Boot | 3.5.9 | 核心框架 |
| Spring AI | 1.1.2 | AI 模型集成框架 |
| Spring WebFlux | - | 流式响应支持 |
| MyBatis-Flex | 1.11.1 | ORM 框架 |
| MySQL | 8.0+ | 主数据库 |
| Redis + Redisson | 7+ | 缓存 / Session / 分布式锁 |
| Resilience4j | 2.2.0 | 熔断器 / 限流 |
| Micrometer | - | 监控指标收集 |
| Knife4j | 4.4.0 | API 文档 |
| Stripe | 31.2.0 | 在线支付 |

### 前端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.5.17 | 前端框架 |
| Vite | 7.0.0 | 构建工具 |
| Ant Design Vue | 4.2.6 | UI 组件库 |
| Pinia | 3.0.3 | 状态管理 |
| ECharts | 5.5.0 | 图表组件 |
| Axios | 1.11.0 | HTTP 客户端 |
| TypeScript | 5.8.0 | 类型安全 |

### 基础设施

| 组件 | 说明 |
|------|------|
| Docker + Docker Compose | 容器化部署 |
| Nginx | 反向代理 / 负载均衡 |
| Prometheus | 指标存储 |
| Grafana | 监控可视化 |

---

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        客户端层                                  │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│   │ Java SDK │  │ 前端应用  │  │ HTTP API │  │  其他SDK  │        │
│   └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘        │
└────────┼─────────────┼─────────────┼─────────────┼──────────────┘
         │             │             │             │
         ▼             ▼             ▼             ▼
┌─────────────────────────────────────────────────────────────────┐
│                     API 网关层                                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                   认证 & 鉴权层                           │    │
│  │     • API Key 验证    • 权限检查    • 配额控制              │   │
│  └──────────────────────────┬──────────────────────────────┘    │
│                              │                                  │
│  ┌──────────────────────────▼──────────────────────────────┐    │
│  │                    智能路由层                             │    │
│  │  • 成本优先路由    • 速度优先路由    • 健康检查              │    │
│  │  • 自动 Fallback   • 负载均衡      • 模型映射              │    │
│  └──────────────────────────┬──────────────────────────────┘    │
│                              │                                  │
│  ┌──────────────────────────▼──────────────────────────────┐    │
│  │                   模型调用层 (Spring AI)                  │    │
│  │  • 统一接口封装    • 流式响应    • 智能重试                 │    │
│  │  • 超时控制        • 响应缓存    • 错误处理                 │    │
│  └───────┬────────────────┬────────────────┬───────────────┘    │
└──────────┼────────────────┼────────────────┼────────────────────┘
           │                │                │
           ▼                ▼                ▼
    ┌──────────┐     ┌──────────┐     ┌──────────┐
    │ 通义千问  │     │  智谱AI   │     │ DeepSeek │  ...更多模型
    └──────────┘     └──────────┘     └──────────┘
```

---

## 📁 项目结构

```
yu-ai-router/
├── frontend/                        # 前端项目
│   ├── src/
│   │   ├── api/                     # API 接口定义
│   │   ├── components/              # 公共组件
│   │   ├── layouts/                 # 布局组件
│   │   ├── pages/                   # 页面
│   │   │   ├── admin/               # 管理员页面
│   │   │   └── user/                # 用户页面
│   │   ├── router/                  # 路由配置
│   │   └── stores/                  # Pinia 状态管理
│   ├── Dockerfile                   # 前端 Docker 配置
│   └── package.json
├── sql/                             # SQL 脚本
│   └── create_table.sql             # 建表语句
├── src/main/java/com/yupi/airouter/
│   ├── adapter/                     # 模型适配器
│   │   ├── DashscopeAdapter.java    # 通义千问适配器
│   │   ├── DeepSeekAdapter.java     # DeepSeek 适配器
│   │   └── ZhipuAIAdapter.java      # 智谱AI 适配器
│   ├── annotation/                  # 自定义注解
│   ├── aop/                         # AOP 切面
│   ├── config/                      # 配置类
│   ├── controller/                  # 控制器
│   │   ├── ChatController.java      # 对话接口 (OpenAI 兼容)
│   │   ├── ImageController.java     # AI 绘图接口
│   │   └── ...
│   ├── filter/                      # 过滤器
│   ├── mapper/                      # MyBatis Mapper
│   ├── metrics/                     # 监控指标
│   ├── model/                       # 数据模型
│   │   ├── dto/                     # 数据传输对象
│   │   ├── entity/                  # 实体类
│   │   ├── enums/                   # 枚举
│   │   └── vo/                      # 视图对象
│   ├── plugin/                      # 插件系统
│   ├── service/                     # 业务逻辑
│   ├── strategy/                    # 路由策略
│   └── utils/                       # 工具类
├── yu-ai-router-sdk/                # Java SDK
├── docker-compose.yml               # Docker Compose 配置
├── Dockerfile                       # 后端 Docker 配置
├── grafana-dashboard.json           # Grafana 监控大盘
├── prometheus.yml                   # Prometheus 配置
└── pom.xml                          # Maven 配置
```

---

## 🚀 快速开始

### 环境要求

- JDK 21+
- Node.js 22+
- MySQL 8.0+
- Redis 7.0+

### 本地开发

#### 1. 克隆项目

```bash
git clone https://github.com/your-repo/yu-ai-router.git
cd yu-ai-router
```

#### 2. 初始化数据库

```bash
# 执行 SQL 脚本创建数据库和表
mysql -u root -p < sql/create_table.sql
```

#### 3. 配置 AI 模型密钥

修改 `src/main/resources/application-local.yml`：

```yaml
# 通义千问配置
spring:
  ai:
    dashscope:
      api-key: your-qwen-api-key

# DeepSeek 配置
    deepseek:
      api-key: your-deepseek-api-key

# 智谱AI 配置
    zhipuai:
      api-key: your-zhipu-api-key
```

#### 4. 启动后端服务

```bash
./mvnw spring-boot:run
```

访问接口文档：http://localhost:8123/api/doc.html

#### 5. 启动前端

```bash
cd frontend
npm install
npm run dev
```

访问前端页面：http://localhost:5173

---

## 🐳 Docker 部署

### 一键部署（推荐）

```bash
docker-compose up -d
```

服务访问地址：

| 服务 | 地址 |
|------|------|
| 前端页面 | http://localhost |
| 后端 API | http://localhost:8123/api |
| 接口文档 | http://localhost:8123/api/doc.html |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3000 (admin/admin) |

### 单独构建

```bash
# 构建后端镜像
docker build -t yu-ai-backend .

# 构建前端镜像
cd frontend
docker build -t yu-ai-frontend .
```

---

## 📡 API 使用

### 对话接口（兼容 OpenAI 格式）

```bash
curl -X POST "http://localhost:8123/api/v1/chat/completions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "qwen-plus",
    "messages": [
      {"role": "user", "content": "你好"}
    ],
    "stream": true
  }'
```

### 支持的模型

| 模型标识 | 提供商 | 说明 |
|----------|--------|------|
| `qwen-plus` | 通义千问 | 增强版，性能更强 |
| `qwen-turbo` | 通义千问 | 快速版，响应更快 |
| `qwen-max` | 通义千问 | 旗舰版，支持深度思考 |
| `glm-4.7` | 智谱AI | 高智能旗舰模型 |
| `glm-4.6` | 智谱AI | 超强性能模型 |
| `glm-4.7-flash` | 智谱AI | 免费模型 |
| `deepseek-chat` | DeepSeek | 对话模型 |
| `deepseek-reasoner` | DeepSeek | 支持深度思考 |
| `deepseek-coder` | DeepSeek | 代码模型 |

### 智能路由

```bash
# 使用 auto 模式，自动选择最优模型
curl -X POST "http://localhost:8123/api/v1/chat/completions" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "auto",
    "messages": [{"role": "user", "content": "你好"}],
    "routing_strategy": "cost_first"
  }'
```

路由策略：
- `cost_first`：成本优先，选择最便宜的健康模型
- `latency_first`：速度优先，选择响应最快的模型
- `round_robin`：轮询所有健康模型

---

## 📦 Java SDK

### 安装

```xml
<dependency>
    <groupId>com.yupi</groupId>
    <artifactId>yu-ai-router-sdk</artifactId>
    <version>1.0.0</version>
</dependency>
```

### 使用示例

```java
// 创建客户端
YuAIClient client = YuAIClient.builder()
    .apiKey("sk-your-api-key")
    .baseUrl("http://localhost:8123/api")
    .build();

// 同步调用
String response = client.chat("你好");
System.out.println(response);

// 指定模型
String response = client.chat("你好", "qwen-plus");

// 流式调用
client.chatStream("讲一个故事", new StreamCallback() {
    @Override
    public void onMessage(String content) {
        System.out.print(content);
    }
    
    @Override
    public void onComplete() {
        System.out.println("\n完成");
    }
});
```

---

## 📊 监控指标

### Prometheus 指标

| 指标名称 | 类型 | 说明 |
|----------|------|------|
| `ai_gateway_requests_total` | Counter | 请求总数 |
| `ai_gateway_request_duration` | Histogram | 请求延迟分布 |
| `ai_gateway_tokens_total` | Counter | Token 消耗总数 |
| `ai_gateway_model_health` | Gauge | 模型健康状态 |

### Grafana 仪表盘

项目提供预置的 Grafana 仪表盘（`grafana-dashboard.json`），包含：

- 实时 QPS 曲线
- 各模型成功率/错误率
- Token 消耗趋势
- 平均延迟分布
- 模型健康状态

---

## 🔧 配置说明

### 应用配置

```yaml
# application.yml
server:
  port: 8123
  servlet:
    context-path: /api

spring:
  datasource:
    url: jdbc:mysql://localhost:3306/yu_ai_router
    username: root
    password: 123456
  
  data:
    redis:
      host: localhost
      port: 6379

# AI 模型配置
ai:
  routing:
    default-strategy: cost_first  # 默认路由策略
    fallback-enabled: true        # 启用自动 Fallback
    retry-times: 3                # 重试次数
    retry-interval: 1000          # 重试间隔（毫秒）
```

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `SPRING_PROFILES_ACTIVE` | 激活配置 | local |
| `SPRING_DATASOURCE_URL` | 数据库连接 | - |
| `SPRING_DATA_REDIS_HOST` | Redis 地址 | localhost |
| `AI_DASHSCOPE_API_KEY` | 通义千问 Key | - |
| `AI_DEEPSEEK_API_KEY` | DeepSeek Key | - |
| `AI_ZHIPU_API_KEY` | 智谱 AI Key | - |

---

## 🗃️ 数据库设计

### 核心表

| 表名 | 说明 |
|------|------|
| `user` | 用户表 |
| `api_key` | API Key 表 |
| `model_provider` | 模型提供者表 |
| `model` | 模型表 |
| `request_log` | 请求日志表 |
| `plugin_config` | 插件配置表 |
| `user_provider_key` | 用户自带密钥表（BYOK） |

---

## 🔐 安全特性

- **API Key 认证**：所有 API 请求需携带有效的 API Key
- **请求限流**：支持 API Key 级别和 IP 级别的限流
- **IP 黑名单**：可配置 IP 黑名单，防止恶意请求
- **TraceId 追踪**：全链路追踪，便于问题排查
- **HTTPS 支持**：建议生产环境启用 HTTPS

---

## 📝 开发计划

- [x] 基础功能（用户系统、对话、API Key 管理）
- [x] 多模型接入（通义千问、智谱AI、DeepSeek）
- [x] 智能路由与 Fallback
- [x] 请求限流与 IP 黑名单
- [x] Java SDK
- [x] Prometheus + Grafana 监控
- [x] Stripe 在线支付
- [x] AI 绘图功能
- [x] 插件系统（Web 搜索、PDF 解析）
- [x] BYOK（用户自带密钥）
- [ ] Function Calling 支持
- [ ] Go/Python SDK

---

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 发起 Pull Request

---

## 📄 开源协议

本项目采用 MIT 协议开源，详见 [LICENSE](LICENSE) 文件。

---

## 🙏 致谢

- [Spring AI](https://spring.io/projects/spring-ai) - AI 模型集成框架
- [OpenRouter](https://openrouter.ai/) - 产品设计参考
- [Ant Design Vue](https://antdv.com/) - UI 组件库
- [编程导航](https://codefather.cn) - 项目孵化平台

---

<p align="center">
  Made with ❤️ by <a href="https://codefather.cn">编程导航学习圈</a>
</p>
