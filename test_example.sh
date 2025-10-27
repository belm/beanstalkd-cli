#!/bin/bash

echo "=== Beanstalkd CLI 测试示例 ==="
echo ""

# 1. 查看服务器状态
echo "1. 查看服务器统计信息:"
./beanstalkd-cli stats | head -20
echo ""

# 2. 列出所有 tubes
echo "2. 列出所有 tubes:"
./beanstalkd-cli list-tubes
echo ""

# 3. 查看 default tube 统计
echo "3. 查看 default tube 统计:"
./beanstalkd-cli stats-tube default
echo ""

# 4. 插入几个任务
echo "4. 插入任务到 default tube:"
./beanstalkd-cli put "任务 1: 处理订单"
./beanstalkd-cli put "任务 2: 发送邮件" -r 100
./beanstalkd-cli put "任务 3: 延迟任务" -d 10s
echo ""

# 5. 再次查看 tube 统计
echo "5. 插入后的 tube 统计:"
./beanstalkd-cli stats-tube default
echo ""

# 6. 预留一个任务
echo "6. 预留一个任务:"
./beanstalkd-cli reserve
echo ""

# 7. 查看任务详情（假设 job id 为 123）
echo "7. 查看任务统计 (使用上面预留的 job ID):"
echo "   运行: ./beanstalkd-cli stats-job <job-id>"
echo ""

# 8. 删除任务
echo "8. 删除任务 (使用上面预留的 job ID):"
echo "   运行: ./beanstalkd-cli delete <job-id>"
echo ""

echo "=== 测试完成 ==="
echo ""
echo "更多命令示例:"
echo "  查看帮助: ./beanstalkd-cli --help"
echo "  查看命令详情: ./beanstalkd-cli put --help"
echo "  使用其他 tube: ./beanstalkd-cli -t mytube put 'data'"
