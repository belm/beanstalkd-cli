#!/bin/bash

# Beanstalkd CLI 测试运行脚本

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Beanstalkd CLI 测试套件 ===${NC}"
echo ""

# 检查 beanstalkd 是否运行
echo -e "${YELLOW}检查 beanstalkd 服务...${NC}"
if ! nc -z 127.0.0.1 11300 2>/dev/null; then
    echo -e "${RED}✗ 无法连接到 beanstalkd (127.0.0.1:11300)${NC}"
    echo "请确保 beanstalkd 服务正在运行"
    exit 1
fi
echo -e "${GREEN}✓ beanstalkd 服务运行正常${NC}"
echo ""

# 运行测试
echo -e "${YELLOW}运行测试用例...${NC}"
echo ""

# 1. 连接测试
echo -e "${YELLOW}1. 连接测试${NC}"
go test -v -run TestConnection || true
echo ""

# 2. 任务操作测试
echo -e "${YELLOW}2. 任务操作测试${NC}"
go test -v -run "TestPut|TestReserve|TestDelete|TestRelease|TestBury|TestTouch|TestKick" || true
echo ""

# 3. 查看操作测试
echo -e "${YELLOW}3. 查看操作测试${NC}"
go test -v -run TestPeek || true
echo ""

# 4. Tube 操作测试
echo -e "${YELLOW}4. Tube 操作测试${NC}"
go test -v -run "TestListTubes|TestMultipleTubes|TestTubeStats|TestTubeIsolation" || true
echo ""

# 5. 统计信息测试
echo -e "${YELLOW}5. 统计信息测试${NC}"
go test -v -run TestStats || true
echo ""

# 6. 集成测试
echo -e "${YELLOW}6. 集成测试${NC}"
go test -v -run "TestProducerConsumer|TestRetry|TestPriority" || true
echo ""

# 7. 性能基准测试
echo -e "${YELLOW}7. 性能基准测试${NC}"
go test -bench=. -benchmem || true
echo ""

# 测试覆盖率
echo -e "${YELLOW}生成测试覆盖率报告...${NC}"
go test -coverprofile=coverage.out
go tool cover -func=coverage.out
echo ""
echo -e "${GREEN}测试覆盖率报告已生成: coverage.out${NC}"
echo -e "使用以下命令查看 HTML 报告:"
echo -e "  go tool cover -html=coverage.out"
echo ""

echo -e "${GREEN}=== 测试完成 ===${NC}"
