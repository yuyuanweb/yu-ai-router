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
    userAccount  varchar(256)                           not null comment '账号',
    userPassword varchar(512)                           not null comment '密码',
    userName     varchar(256)                           null comment '用户昵称',
    userAvatar   varchar(1024)                          null comment '用户头像',
    userProfile  varchar(512)                           null comment '用户简介',
    userRole     varchar(256) default 'user'            not null comment '用户角色：user/admin',
    editTime     datetime     default CURRENT_TIMESTAMP not null comment '编辑时间',
    createTime   datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    updateTime   datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    isDelete     tinyint      default 0                 not null comment '是否删除',
    UNIQUE KEY uk_userAccount (userAccount),
    INDEX idx_userName (userName)
) comment '用户' collate = utf8mb4_unicode_ci;

-- API Key 表
create table if not exists api_key
(
    id           bigint auto_increment comment 'id' primary key,
    userId       bigint                                 not null comment '用户id',
    keyValue     varchar(128)                           not null comment 'API Key值（sk-xxx格式）',
    keyName      varchar(128)                           null comment 'Key名称/备注',
    status       varchar(32)  default 'active'          not null comment '状态：active/inactive/revoked',
    totalTokens  bigint       default 0                 not null comment '已使用Token总数',
    lastUsedTime datetime                               null comment '最后使用时间',
    createTime   datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    updateTime   datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    isDelete     tinyint      default 0                 not null comment '是否删除',
    UNIQUE KEY uk_keyValue (keyValue),
    INDEX idx_userId (userId),
    INDEX idx_status (status)
) comment 'API Key' collate = utf8mb4_unicode_ci;

-- 请求日志表（用于记录每次请求和Token消耗）
create table if not exists request_log
(
    id              bigint auto_increment comment 'id' primary key,
    userId          bigint                                 null comment '用户id',
    apiKeyId        bigint                                 null comment 'API Key id',
    modelName       varchar(128)                           not null comment '使用的模型名称',
    promptTokens    int          default 0                 not null comment '输入Token数',
    completionTokens int         default 0                 not null comment '输出Token数',
    totalTokens     int          default 0                 not null comment '总Token数',
    duration        int          default 0                 not null comment '请求耗时（毫秒）',
    status          varchar(32)  default 'success'         not null comment '状态：success/failed',
    errorMessage    text                                   null comment '错误信息',
    createTime      datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    updateTime      datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    INDEX idx_userId (userId),
    INDEX idx_apiKeyId (apiKeyId),
    INDEX idx_createTime (createTime)
) comment '请求日志' collate = utf8mb4_unicode_ci;

-- 初始化用户数据
-- 密码是 12345678（MD5 加密 + 盐值 yupi）
INSERT INTO user (id, userAccount, userPassword, userName, userAvatar, userProfile, userRole) VALUES
(1, 'admin', '10670d38ec32fa8102be6a37f8cb52bf', '管理员', 'https://www.codefather.cn/logo.png', '系统管理员', 'admin'),
(2, 'user', '10670d38ec32fa8102be6a37f8cb52bf', '普通用户', 'https://www.codefather.cn/logo.png', '我是一个普通用户', 'user'),
(3, 'test', '10670d38ec32fa8102be6a37f8cb52bf', '测试账号', 'https://www.codefather.cn/logo.png', '这是一个测试账号', 'user');
