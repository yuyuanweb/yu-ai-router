# 数据库初始化
# @author <a href="https://codefather.cn">编程导航学习圈</a>

-- 创建库
create database if not exists yu_ai_router;

-- 切换库
use yu_ai_router;

-- 用户表
create table if not exists user
(
    id           bigint auto_increment comment 'id' primary key,
    userAccount  varchar(256)                             not null comment '账号',
    userPassword varchar(512)                             not null comment '密码',
    userName     varchar(256)                             null comment '用户昵称',
    userAvatar   varchar(1024)                            null comment '用户头像',
    userProfile  varchar(512)                             null comment '用户简介',
    userRole     varchar(256)   default 'user'            not null comment '用户角色：user/admin',
    userStatus   varchar(32)    default 'active'          not null comment '用户状态：active/disabled',
    tokenQuota   bigint         default -1                not null comment 'Token配额（-1表示无限制）',
    usedTokens   bigint         default 0                 not null comment '已使用Token数',
    balance      decimal(12, 4) default 0.0000            not null comment '账户余额（元）',
    editTime     datetime       default CURRENT_TIMESTAMP not null comment '编辑时间',
    createTime   datetime       default CURRENT_TIMESTAMP not null comment '创建时间',
    updateTime   datetime       default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    isDelete     tinyint        default 0                 not null comment '是否删除',
    UNIQUE KEY uk_userAccount (userAccount),
    INDEX idx_userName (userName),
    INDEX idx_userStatus (userStatus)
) comment '用户' collate = utf8mb4_unicode_ci;

-- API Key 表
create table if not exists api_key
(
    id           bigint auto_increment comment 'id' primary key,
    userId       bigint                                not null comment '用户id',
    keyValue     varchar(128)                          not null comment 'API Key值（sk-xxx格式）',
    keyName      varchar(128)                          null comment 'Key名称/备注',
    status       varchar(32) default 'active'          not null comment '状态：active/inactive/revoked',
    totalTokens  bigint      default 0                 not null comment '已使用Token总数',
    lastUsedTime datetime                              null comment '最后使用时间',
    createTime   datetime    default CURRENT_TIMESTAMP not null comment '创建时间',
    updateTime   datetime    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    isDelete     tinyint     default 0                 not null comment '是否删除',
    UNIQUE KEY uk_keyValue (keyValue),
    INDEX idx_userId (userId),
    INDEX idx_status (status)
) comment 'API Key' collate = utf8mb4_unicode_ci;

-- 请求日志表（用于记录每次请求和Token消耗）
create table if not exists request_log
(
    id               bigint auto_increment comment 'id' primary key,
    traceId          varchar(64)                              null comment '链路追踪ID',
    userId           bigint                                   null comment '用户id',
    apiKeyId         bigint                                   null comment 'API Key id',
    modelId          bigint                                   null comment '实际调用的模型id',
    requestModel     varchar(128)                             null comment '请求的模型标识',
    modelName        varchar(128)                             null comment '使用的模型名称（兼容字段）',
    requestType      varchar(32)    default 'chat'            not null comment '请求类型：chat/embedding/image',
    source           varchar(32)    default 'web'             not null comment '调用来源：web/api',
    promptTokens     int            default 0                 not null comment '输入Token数',
    completionTokens int            default 0                 not null comment '输出Token数',
    totalTokens      int            default 0                 not null comment '总Token数',
    cost             decimal(12, 6) default 0                 not null comment '本次请求费用（元）',
    duration         int            default 0                 not null comment '请求耗时（毫秒）',
    status           varchar(32)    default 'success'         not null comment '状态：success/failed',
    errorMessage     text                                     null comment '错误信息',
    errorCode        varchar(64)                              null comment '错误码',
    routingStrategy  varchar(32)                              null comment '使用的路由策略',
    isFallback       tinyint        default 0                 not null comment '是否为Fallback请求',
    clientIp         varchar(64)                              null comment '客户端IP',
    userAgent        varchar(512)                             null comment 'User-Agent',
    createTime       datetime       default CURRENT_TIMESTAMP not null comment '创建时间',
    updateTime       datetime       default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    INDEX idx_traceId (traceId),
    INDEX idx_userId (userId),
    INDEX idx_apiKeyId (apiKeyId),
    INDEX idx_modelId (modelId),
    INDEX idx_source (source),
    INDEX idx_createTime (createTime)
) comment '请求日志' collate = utf8mb4_unicode_ci;

-- 模型提供者表
create table if not exists model_provider
(
    id           bigint auto_increment comment 'id' primary key,
    providerName varchar(64)                             not null comment '提供者名称（如：qwen/zhipu/deepseek）',
    displayName  varchar(128)                            not null comment '显示名称（如：通义千问/智谱AI/DeepSeek）',
    baseUrl      varchar(512)                            not null comment 'API基础URL',
    apiKey       varchar(512)                            not null comment 'API密钥',
    status       varchar(32)   default 'active'          not null comment '状态：active/inactive/maintenance',
    healthStatus varchar(32)   default 'unknown'         not null comment '健康状态：healthy/unhealthy/degraded/unknown',
    avgLatency   int           default 0                 not null comment '平均延迟（毫秒）',
    successRate  decimal(5, 2) default 100.00            not null comment '成功率（百分比）',
    priority     int           default 100               not null comment '优先级（越大越优先）',
    config       text                                    null comment '额外配置（JSON格式）',
    createTime   datetime      default CURRENT_TIMESTAMP not null comment '创建时间',
    updateTime   datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    isDelete     tinyint       default 0                 not null comment '是否删除',
    UNIQUE KEY uk_providerName (providerName),
    INDEX idx_status (status),
    INDEX idx_healthStatus (healthStatus)
) comment '模型提供者' collate = utf8mb4_unicode_ci;

-- 模型表
create table if not exists model
(
    id               bigint auto_increment comment 'id' primary key,
    providerId       bigint                                   not null comment '提供者id',
    modelKey         varchar(128)                             not null comment '模型标识（如：qwen-plus）',
    modelName        varchar(128)                             not null comment '模型显示名称',
    modelType        varchar(32)    default 'chat'            not null comment '模型类型：chat/embedding/image/audio',
    description      varchar(512)                             null comment '模型描述',
    contextLength    int            default 4096              not null comment '上下文长度限制',
    inputPrice       decimal(10, 6) default 0                 not null comment '输入价格（元/千Token）',
    outputPrice      decimal(10, 6) default 0                 not null comment '输出价格（元/千Token）',
    status           varchar(32)    default 'active'          not null comment '状态：active/inactive/deprecated',
    healthStatus     varchar(32)    default 'unknown'         not null comment '健康状态：healthy/unhealthy/degraded/unknown',
    avgLatency       int            default 0                 not null comment '平均延迟（毫秒）',
    successRate      decimal(5, 2)  default 100.00            not null comment '成功率（百分比）',
    score            decimal(10, 4) default 0                 not null comment '综合得分（越低越好）',
    priority         int            default 100               not null comment '优先级（越大越优先）',
    defaultTimeout   int            default 60000             not null comment '默认超时时间（毫秒）',
    supportReasoning tinyint        default 0                 not null comment '是否支持深度思考：0=不支持，1=支持',
    capabilities     varchar(512)                             null comment '能力标签（JSON数组）',
    createTime       datetime       default CURRENT_TIMESTAMP not null comment '创建时间',
    updateTime       datetime       default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    isDelete         tinyint        default 0                 not null comment '是否删除',
    UNIQUE KEY uk_modelKey (modelKey),
    INDEX idx_providerId (providerId),
    INDEX idx_modelType (modelType),
    INDEX idx_status (status),
    INDEX idx_healthStatus (healthStatus)
) comment '模型' collate = utf8mb4_unicode_ci;

-- 充值记录表
create table if not exists recharge_record
(
    id            bigint auto_increment comment 'id' primary key,
    userId        bigint                                not null comment '用户id',
    amount        decimal(12, 4)                        not null comment '充值金额（元）',
    paymentMethod varchar(32)                           not null comment '支付方式：stripe/alipay/wechat',
    paymentId     varchar(256)                          null comment '第三方支付ID',
    status        varchar(32) default 'pending'         not null comment '状态：pending/success/failed/refunded',
    description   varchar(512)                          null comment '充值说明',
    createTime    datetime    default CURRENT_TIMESTAMP not null comment '创建时间',
    updateTime    datetime    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    INDEX idx_paymentId (paymentId),
    INDEX idx_status (status),
    INDEX idx_userId (userId)
) comment '充值记录' collate = utf8mb4_unicode_ci;

-- 消费账单表
create table if not exists billing_record
(
    id            bigint auto_increment comment 'id' primary key,
    userId        bigint                                not null comment '用户id',
    requestLogId  bigint                                null comment '关联的请求日志ID',
    amount        decimal(12, 4)                        not null comment '消费金额（元）',
    balanceBefore decimal(12, 4)                        not null comment '消费前余额（元）',
    balanceAfter  decimal(12, 4)                        not null comment '消费后余额（元）',
    description   varchar(512)                          null comment '消费说明',
    billingType   varchar(32) default 'api_call'        not null comment '账单类型：api_call/recharge/refund',
    createTime    datetime    default CURRENT_TIMESTAMP not null comment '创建时间',
    index idx_billingType (billingType),
    index idx_createTime (createTime),
    index idx_requestLogId (requestLogId),
    index idx_userId (userId)
) comment '消费账单' collate = utf8mb4_unicode_ci;

-- 图片生成记录表
create table image_generation_record
(
    id            bigint auto_increment comment 'id'
        primary key,
    userId        bigint                                   not null comment '用户id',
    apiKeyId      bigint                                   null comment 'API Key id',
    modelId       bigint                                   not null comment '使用的模型id',
    modelKey      varchar(128)                             not null comment '模型标识',
    prompt        text                                     not null comment '生成提示词',
    revisedPrompt text                                     null comment '修订后的提示词',
    imageUrl      varchar(1024)                            null comment '图片URL',
    imageData     longtext                                 null comment 'Base64图片数据',
    size          varchar(32)                              null comment '图片尺寸',
    quality       varchar(32)                              null comment '图片质量',
    status        varchar(32)    default 'success'         not null comment '状态：success/failed',
    cost          decimal(12, 4) default 0.0000            not null comment '生成费用（元）',
    duration      int                                      null comment '耗时（毫秒）',
    errorMessage  varchar(512)                             null comment '错误信息',
    clientIp      varchar(128)                             null comment '客户端IP',
    createTime    datetime       default CURRENT_TIMESTAMP not null comment '创建时间',
    INDEX idx_userId (userId),
    INDEX idx_modelId (modelId)
) comment '图片生成记录' collate = utf8mb4_unicode_ci;

-- 插件配置表
create table if not exists plugin_config
(
    id          bigint auto_increment comment 'id' primary key,
    pluginKey   varchar(64)                           not null comment '插件标识：web_search/pdf_parser/image_recognition',
    pluginName  varchar(128)                          not null comment '插件名称',
    pluginType  varchar(32) default 'builtin'         not null comment '插件类型：builtin/custom',
    description varchar(512)                          null comment '插件描述',
    config      text                                  null comment '插件配置（JSON）',
    status      varchar(32) default 'active'          not null comment '状态：active/inactive',
    priority    int         default 100               not null comment '优先级',
    createTime  datetime    default CURRENT_TIMESTAMP not null comment '创建时间',
    updateTime  datetime    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    isDelete    tinyint     default 0                 not null comment '是否删除',
    UNIQUE KEY uk_pluginKey (pluginKey)
) comment '插件配置' collate = utf8mb4_unicode_ci;

-- 用户提供者密钥表（用户自己的 API Key）
create table if not exists user_provider_key
(
    id           bigint auto_increment comment 'id' primary key,
    userId       bigint                                 not null comment '用户 ID',
    providerId   bigint                                 not null comment '提供者 ID',
    providerName varchar(64)                            not null comment '提供者名称（冗余字段，便于查询）',
    apiKey       varchar(512)                           not null comment 'API Key（加密存储）',
    status       varchar(32)  default 'active'          not null comment '状态：active/inactive',
    createTime   datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    updateTime   datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    isDelete     tinyint      default 0                 not null comment '是否删除',
    KEY uk_user_provider (userId, providerId)
) comment '用户提供者密钥（BYOK）' collate = utf8mb4_unicode_ci;

-- 初始化用户数据
-- 密码是 12345678（MD5 加密 + 盐值 yupi）
INSERT INTO user (id, userAccount, userPassword, userName, userAvatar, userProfile, userRole, tokenQuota)
VALUES (1, 'admin', '10670d38ec32fa8102be6a37f8cb52bf', '管理员', 'https://www.codefather.cn/logo.png', '系统管理员', 'admin', -1),
       (2, 'user', '10670d38ec32fa8102be6a37f8cb52bf', '普通用户', 'https://www.codefather.cn/logo.png', '我是一个普通用户', 'user', 100000),
       (3, 'test', '10670d38ec32fa8102be6a37f8cb52bf', '测试账号', 'https://www.codefather.cn/logo.png', '这是一个测试账号', 'user', 50000);

-- 初始化模型提供者数据
INSERT INTO model_provider (providerName, displayName, baseUrl, apiKey, status, priority)
VALUES ('qwen', '通义千问', 'https://dashscope.aliyuncs.com/compatible-mode', 'YOUR_QWEN_API_KEY', 'active', 100),
       ('zhipu', '智谱AI', 'https://open.bigmodel.cn/api/paas', 'YOUR_ZHIPU_API_KEY', 'active', 90),
       ('deepseek', 'DeepSeek', 'https://api.deepseek.com', 'YOUR_DEEPSEEK_API_KEY', 'active', 80);

-- 初始化模型数据
INSERT INTO model (providerId, modelKey, modelName, modelType, description, contextLength, inputPrice, outputPrice, priority, supportReasoning)
VALUES
-- 通义千问模型
(1, 'qwen-plus', 'Qwen Plus', 'chat', '通义千问增强版，性能更强', 32768, 0.004, 0.004, 100, 1),
(1, 'qwen-turbo', 'Qwen Turbo', 'chat', '通义千问快速版，响应更快', 8192, 0.002, 0.002, 90, 0),
(1, 'qwen-max', 'Qwen Max', 'chat', '通义千问旗舰版，能力最强，支持深度思考', 32768, 0.04, 0.04, 110, 1),

-- 智谱AI模型
(2, 'glm-4.7', 'GLM-4.7', 'chat', '智谱AI高智能旗舰模型，通用对话、推理与智能体能力全面升级', 204800, 0.05, 0.05, 100, 1),
(2, 'glm-4.6', 'GLM-4.6', 'chat', '智谱AI超强性能模型，高级编码能力、强大推理以及工具调用能力', 204800, 0.05, 0.05, 90, 1),
(2, 'glm-4.7-flash', 'GLM-4.7 Flash', 'chat', '智谱AI免费模型，最新基座模型的普惠版本', 204800, 0.0001, 0.0001, 80, 0),

-- DeepSeek模型
(3, 'deepseek-reasoner', 'DeepSeek Reasoner', 'chat', 'DeepSeek对话模型，支持深度思考', 32768, 0.01, 0.02, 100, 1),
(3, 'deepseek-chat', 'DeepSeek Chat', 'chat', 'DeepSeek对话模型', 32768, 0.001, 0.002, 100, 0),
(3, 'deepseek-coder', 'DeepSeek Coder', 'chat', 'DeepSeek代码模型', 32768, 0.001, 0.002, 90, 0),

-- 阶段七：绘图模型
(1, 'qwen-image-plus', 'Qwen Image Plus', 'image', '通义万相文生图模型', 0, 0.08, 0, 100, 0),
(2, 'cogview-3-plus', 'CogView-3-Plus', 'image', '智谱AI文生图模型', 0, 0.1, 0, 90, 0);

-- 初始化插件数据
-- 注意：Web搜索插件的 API Key 配置在 application-local.yml 中（plugin.serpapi.api-key）
-- SerpApi 注册地址：https://serpapi.com/
INSERT INTO plugin_config (pluginKey, pluginName, pluginType, description, config, status, priority)
VALUES ('web_search', 'Web搜索', 'builtin', '实时联网搜索（SerpApi）', '{"maxResults":5,"searchEngine":"google","timeout":15000}', 'active', 100),
       ('pdf_parser', 'PDF解析', 'builtin', '解析PDF文档内容，提取文本信息', '{"maxPages":50,"maxTextLength":50000}', 'active', 90),
       ('image_recognition', '图片识别', 'builtin', '识别图片内容，返回图片描述', '{"model":"qwen-vl-plus","maxImageSize":4194304}', 'active', 80);
