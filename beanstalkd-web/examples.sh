#!/bin/bash

# Beanstalkd Web UI 配置示例脚本

echo "═══════════════════════════════════════════════════════"
echo "  Beanstalkd Web UI - 配置示例"
echo "═══════════════════════════════════════════════════════"
echo ""

# 示例 1: 使用默认配置
echo "【示例 1】使用默认配置"
echo "  命令: go run server.go"
echo "  Beanstalkd: 127.0.0.1:11300"
echo "  Web 端口: 8080"
echo ""

# 示例 2: 自定义 Beanstalkd 地址
echo "【示例 2】自定义 Beanstalkd 地址"
echo "  命令: go run server.go -beanstalkd 192.168.1.100:11300"
echo "  Beanstalkd: 192.168.1.100:11300"
echo "  Web 端口: 8080 (默认)"
echo ""

# 示例 3: 自定义 Web 端口
echo "【示例 3】自定义 Web 端口"
echo "  命令: go run server.go -port 9090"
echo "  Beanstalkd: 127.0.0.1:11300 (默认)"
echo "  Web 端口: 9090"
echo ""

# 示例 4: 同时自定义
echo "【示例 4】同时自定义多个配置"
echo "  命令: go run server.go -beanstalkd 192.168.1.100:11300 -port 9090"
echo "  Beanstalkd: 192.168.1.100:11300"
echo "  Web 端口: 9090"
echo ""

# 示例 5: 使用环境变量
echo "【示例 5】使用环境变量"
echo "  命令: export BEANSTALKD_HOST=192.168.1.100:11300"
echo "        export WEB_PORT=9090"
echo "        go run server.go"
echo "  Beanstalkd: 192.168.1.100:11300"
echo "  Web 端口: 9090"
echo ""

# 示例 6: 一次性环境变量
echo "【示例 6】一次性设置环境变量"
echo "  命令: BEANSTALKD_HOST=192.168.1.100:11300 WEB_PORT=9090 go run server.go"
echo "  Beanstalkd: 192.168.1.100:11300"
echo "  Web 端口: 9090"
echo ""

# 示例 7: 使用启动脚本
echo "【示例 7】使用启动脚本"
echo "  命令: ./start.sh"
echo "  或: BEANSTALKD_HOST=192.168.1.100:11300 ./start.sh"
echo "  或: ./start.sh -beanstalkd 192.168.1.100:11300 -port 9090"
echo ""

echo "═══════════════════════════════════════════════════════"
echo "  选择一个示例运行"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "请选择示例 (1-7) 或按 q 退出:"
read -r choice

case $choice in
    1)
        echo "运行示例 1..."
        go run server.go
        ;;
    2)
        echo "运行示例 2..."
        go run server.go -beanstalkd 192.168.1.100:11300
        ;;
    3)
        echo "运行示例 3..."
        go run server.go -port 9090
        ;;
    4)
        echo "运行示例 4..."
        go run server.go -beanstalkd 192.168.1.100:11300 -port 9090
        ;;
    5)
        echo "运行示例 5..."
        export BEANSTALKD_HOST=192.168.1.100:11300
        export WEB_PORT=9090
        go run server.go
        ;;
    6)
        echo "运行示例 6..."
        BEANSTALKD_HOST=192.168.1.100:11300 WEB_PORT=9090 go run server.go
        ;;
    7)
        echo "运行示例 7..."
        ./start.sh
        ;;
    q|Q)
        echo "退出"
        exit 0
        ;;
    *)
        echo "无效选择"
        exit 1
        ;;
esac
