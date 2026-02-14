#!/bin/sh

set -eu

echo "======================================"
echo "  鱼 AI 网关（Go）- Docker 一键启动"
echo "======================================"
echo ""

if ! docker info >/dev/null 2>&1; then
  echo "Docker 未运行，请先启动 Docker Desktop"
  exit 1
fi

echo "清理旧容器..."
docker compose down
echo ""

echo "启动服务（MySQL / Redis / Go Backend / Frontend）..."
docker compose up -d --build
echo ""

echo "等待服务启动..."
sleep 10
echo ""

echo "服务状态："
docker compose ps
echo ""

echo "访问地址："
echo "  前端页面: http://localhost"
echo "  后端 API: http://localhost:8123/api"
echo "  健康检查: http://localhost:8123/api/health/"
echo ""
echo "查看后端日志：docker compose logs -f backend"
echo "停止服务：docker compose down"
echo ""
echo "说明：本部署不包含 Prometheus / Grafana。"
