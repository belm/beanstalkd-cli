#!/bin/bash

# Beanstalkd Web UI 启动脚本

echo "🚀 启动 Beanstalkd Web UI 管理后台..."
echo ""

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到 Go 环境"
    echo "请先安装 Go: https://golang.org/dl/"
    exit 1
fi

# 读取配置（从环境变量或使用默认值）
BEANSTALKD_HOST=${BEANSTALKD_HOST:-127.0.0.1:11300}
WEB_PORT=${WEB_PORT:-8080}

# 解析 Beanstalkd 主机和端口
IFS=':' read -r BEANSTALKD_IP BEANSTALKD_PORT <<< "$BEANSTALKD_HOST"

# 检查 beanstalkd 是否运行
echo "🔍 检查 Beanstalkd 服务..."
if ! nc -z "$BEANSTALKD_IP" "$BEANSTALKD_PORT" 2>/dev/null; then
    echo "⚠️  警告: 无法连接到 Beanstalkd ($BEANSTALKD_HOST)"
    echo "请确保 Beanstalkd 服务正在运行"
    echo ""
    read -p "是否继续启动 Web UI? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    echo "✓ Beanstalkd 服务运行正常"
fi

echo ""
echo "📋 配置信息:"
echo "   Beanstalkd: $BEANSTALKD_HOST"
echo "   Web 端口: $WEB_PORT"
echo ""

# 启动 Web 服务器
echo "🌐 启动 Web 服务器..."
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "💡 提示: 可通过环境变量或命令行参数自定义配置"
echo "   export BEANSTALKD_HOST=192.168.1.100:11300"
echo "   export WEB_PORT=9090"
echo "   或使用: go run server.go -beanstalkd <host> -port <port>"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 传递参数给 server.go（如果提供了命令行参数）
if [ "$#" -gt 0 ]; then
    go run server.go "$@"
else
    go run server.go
fi
