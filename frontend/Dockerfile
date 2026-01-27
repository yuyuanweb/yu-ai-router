# 前端 Dockerfile
FROM node:22-alpine AS build
WORKDIR /app

# 设置淘宝镜像源
RUN npm config set registry https://registry.npmmirror.com

# 复制依赖文件并安装（利用 Docker 缓存）
COPY package*.json ./
RUN npm ci --legacy-peer-deps || npm install --legacy-peer-deps

# 复制源码并构建
COPY . .
RUN npm run build

# Nginx 运行阶段
FROM nginx:alpine
WORKDIR /usr/share/nginx/html

# 复制构建产物
COPY --from=build /app/dist .

# 复制 Nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]
