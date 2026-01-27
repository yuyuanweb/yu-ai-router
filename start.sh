#!/bin/bash

# 鱼 AI 网关 - 一键启动脚本

echo "======================================"
echo "  鱼 AI 网关 - Docker 一键启动"
echo "======================================"
echo ""

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行，请先启动 Docker Desktop"
    exit 1
fi

echo "✅ Docker 运行正常"
echo ""

# 停止并删除旧容器
echo "🧹 清理旧容器..."
docker-compose down
echo ""

# 构建并启动所有服务
echo "🚀 启动所有服务..."
docker-compose up -d --build

# 等待服务启动
echo ""
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo ""
echo "📊 服务状态检查："
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

docker-compose ps

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "✅ 启动完成！"
echo ""
echo "📱 访问地址："
echo "   前端页面：http://localhost"
echo "   后端 API：http://localhost:8123/api"
echo "   接口文档：http://localhost:8123/api/doc.html"
echo "   监控端点：http://localhost:8123/api/actuator/prometheus"
echo "   Prometheus：http://localhost:9090"
echo "   Grafana：http://localhost:3000 (admin/admin)"
echo ""
echo "📝 查看日志："
echo "   docker-compose logs -f backend"
echo ""
echo "🛑 停止服务："
echo "   docker-compose down"
echo ""
echo "======================================"
