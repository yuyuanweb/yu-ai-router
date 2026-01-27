# 鱼 AI 网关 - Docker 部署指南

## 一、快速开始

### 1. 一键启动（推荐）

```bash
chmod +x start.sh
./start.sh
```

### 2. 手动启动

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f backend
```

### 3. 停止服务

```bash
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

## 二、服务列表

| 服务 | 容器名 | 端口 | 说明 |
|------|--------|------|------|
| MySQL | yu-ai-mysql | 3306 | 数据库 |
| Redis | yu-ai-redis | 6379 | 缓存 |
| Backend | yu-ai-backend | 8123 | 后端应用 |
| Frontend | yu-ai-frontend | 80 | 前端应用 |
| Prometheus | yu-ai-prometheus | 9090 | 监控数据采集 |
| Grafana | yu-ai-grafana | 3000 | 监控可视化 |

## 三、访问地址

### 应用访问
- **前端页面**：http://localhost
- **后端 API**：http://localhost:8123/api
- **接口文档**：http://localhost:8123/api/doc.html

### 监控访问
- **Prometheus**：http://localhost:9090
- **Grafana**：http://localhost:3000
  - 默认账号：`admin`
  - 默认密码：`admin`

### 健康检查
- **后端健康**：http://localhost:8123/api/actuator/health
- **Prometheus 指标**：http://localhost:8123/api/actuator/prometheus

## 四、环境变量配置

可以通过修改 `docker-compose.yml` 中的环境变量来配置：

```yaml
environment:
  # 数据库配置
  - SPRING_DATASOURCE_URL=jdbc:mysql://mysql:3306/yu_ai_router
  - SPRING_DATASOURCE_USERNAME=root
  - SPRING_DATASOURCE_PASSWORD=123456
  
  # Redis 配置
  - SPRING_DATA_REDIS_HOST=redis
  - SPRING_DATA_REDIS_PORT=6379
  
  # 通义千问 API Key（必须配置）
  - QWEN_API_KEY=your_qwen_api_key_here
  
  # JVM 配置
  - JAVA_OPTS=-Xms512m -Xmx1024m
```

## 五、数据持久化

使用 Docker Volume 持久化数据：

| Volume | 说明 |
|--------|------|
| `mysql-data` | MySQL 数据文件 |
| `redis-data` | Redis 数据文件 |
| `prometheus-data` | Prometheus 监控数据 |
| `grafana-data` | Grafana 配置和仪表盘 |

## 六、故障排查

### 1. 查看服务日志

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f mysql
docker-compose logs -f redis
```

### 2. 进入容器

```bash
# 进入后端容器
docker exec -it yu-ai-backend sh

# 进入 MySQL 容器
docker exec -it yu-ai-mysql mysql -uroot -p123456
```

### 3. 重启特定服务

```bash
docker-compose restart backend
docker-compose restart frontend
```

### 4. 检查网络连接

```bash
# 查看网络
docker network ls

# 检查容器网络
docker network inspect yu-ai-router_app-network
```

## 七、生产环境部署建议

### 1. 修改默认密码

修改 `docker-compose.yml` 中的：
- MySQL root 密码
- Grafana admin 密码

### 2. 配置域名和 HTTPS

使用 Nginx 反向代理 + Let's Encrypt SSL 证书。

### 3. 数据备份

```bash
# 备份 MySQL 数据
docker exec yu-ai-mysql mysqldump -uroot -p123456 yu_ai_router > backup.sql

# 恢复数据
docker exec -i yu-ai-mysql mysql -uroot -p123456 yu_ai_router < backup.sql
```

### 4. 资源限制

在 `docker-compose.yml` 中添加资源限制：

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '0.5'
          memory: 512M
```

## 八、常用命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 查看日志
docker-compose logs -f

# 查看状态
docker-compose ps

# 清理所有容器和数据
docker-compose down -v --remove-orphans

# 重新构建镜像
docker-compose build --no-cache

# 拉取最新镜像
docker-compose pull
```

---

**一键启动脚本和完整的 Docker 配置已完成！**
