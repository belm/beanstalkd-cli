.PHONY: all build test clean install run help

# 变量定义
BINARY_NAME=beanstalkd-cli
GO=go
GOFLAGS=-v

# 默认目标
all: build

# 构建
build:
	@echo "正在构建 $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME)
	@echo "构建完成: $(BINARY_NAME)"

# 运行测试
test:
	@echo "运行测试..."
	cd tests && $(GO) test $(GOFLAGS)

# 运行所有测试（详细模式）
test-verbose:
	@echo "运行测试（详细模式）..."
	cd tests && $(GO) test -v

# 运行特定测试
test-connection:
	cd tests && $(GO) test -v -run TestConnection

test-jobs:
	cd tests && $(GO) test -v -run TestPutJob

test-integration:
	cd tests && $(GO) test -v -run TestProducerConsumerFlow

# 性能基准测试
bench:
	@echo "运行性能基准测试..."
	cd tests && $(GO) test -bench=. -benchmem

# 测试覆盖率
coverage:
	@echo "生成测试覆盖率报告..."
	cd tests && $(GO) test -coverprofile=coverage.out
	cd tests && $(GO) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: tests/coverage.html"

# 清理
clean:
	@echo "清理构建文件..."
	rm -f $(BINARY_NAME)
	rm -f tests/coverage.out tests/coverage.html
	@echo "清理完成"

# 安装依赖
deps:
	@echo "下载依赖..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "依赖下载完成"

# 格式化代码
fmt:
	@echo "格式化代码..."
	$(GO) fmt ./...
	@echo "格式化完成"

# 代码检查
vet:
	@echo "运行代码检查..."
	$(GO) vet ./...
	@echo "检查完成"

# 运行 CLI
run:
	./$(BINARY_NAME) $(ARGS)

# 帮助信息
help:
	@echo "Beanstalkd CLI Makefile"
	@echo ""
	@echo "使用方法:"
	@echo "  make build           - 构建二进制文件"
	@echo "  make test            - 运行测试"
	@echo "  make test-verbose    - 运行测试（详细输出）"
	@echo "  make bench           - 运行性能基准测试"
	@echo "  make coverage        - 生成测试覆盖率报告"
	@echo "  make clean           - 清理构建文件"
	@echo "  make deps            - 下载依赖"
	@echo "  make fmt             - 格式化代码"
	@echo "  make vet             - 代码检查"
	@echo "  make run ARGS='...'  - 运行 CLI"
	@echo ""
	@echo "示例:"
	@echo "  make run ARGS='put \"test data\"'"
	@echo "  make run ARGS='stats'"
	@echo "  make test-connection"
