# Beanstalkd CLI 测试套件

[🇨🇳 中文版](README.md) | [🇺🇸 English](README-EN.md)

完整的测试用例集合，用于验证 beanstalkd-cli 的所有功能。

## 测试文件说明

### 1. `connection_test.go` - 连接测试
- 基本连接功能
- 连接超时
- 连接失败处理
- 多个并发连接
- 连接性能基准测试

### 2. `job_operations_test.go` - 任务操作测试
- Put：插入任务（普通、高优先级、延迟、中文、JSON）
- Reserve：预留任务
- Delete：删除任务
- Release：释放任务
- Bury：埋葬任务
- Touch：触摸任务
- Kick：踢出任务
- 性能基准测试

### 3. `peek_operations_test.go` - 查看操作测试
- Peek：查看指定任务
- PeekReady：查看就绪任务
- PeekDelayed：查看延迟任务
- PeekBuried：查看被埋葬任务
- 错误处理测试

### 4. `tube_operations_test.go` - Tube 操作测试
- 列出所有 tubes
- 多 tube 操作
- Tube 统计信息
- Tube 隔离性
- 并发 tube 访问

### 5. `stats_test.go` - 统计信息测试
- 服务器统计
- 任务统计
- Tube 统计

### 6. `integration_test.go` - 集成测试
- 完整的生产者-消费者流程
- 任务重试机制
- 优先级队列

## 运行测试

### 前提条件
确保 beanstalkd 服务运行在 `127.0.0.1:11300`

### 运行所有测试
```bash
cd tests
go test -v
```

### 运行特定测试文件
```bash
go test -v -run TestConnection
go test -v -run TestPutJob
go test -v -run TestProducerConsumerFlow
```

### 运行性能基准测试
```bash
go test -bench=. -benchmem
```

### 查看测试覆盖率
```bash
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 测试覆盖范围

✅ 连接管理
✅ 任务生命周期（put, reserve, delete, release, bury, touch, kick）
✅ 查看操作（peek 系列）
✅ Tube 管理
✅ 统计信息
✅ 优先级处理
✅ 延迟队列
✅ 错误处理
✅ 并发操作
✅ 集成流程

## 测试最佳实践

1. **隔离性**: 每个测试使用独立的 tube 名称
2. **清理**: 使用 defer 确保测试资源被清理
3. **跳过**: 无法连接服务器时跳过测试而不是失败
4. **日志**: 使用 t.Log 输出详细的测试信息
5. **并发**: 测试并发场景确保线程安全

## 故障排查

### 测试失败
- 确认 beanstalkd 服务正在运行
- 检查端口 11300 是否可访问
- 查看测试日志了解详细错误

### 测试跳过
如果看到 "SKIP" 信息，通常是因为无法连接到 beanstalkd 服务器。

## 持续集成

这些测试可以集成到 CI/CD 流程中：

```yaml
# .github/workflows/test.yml 示例
test:
  steps:
    - name: Start beanstalkd
      run: beanstalkd -l 127.0.0.1 -p 11300 &
    - name: Run tests
      run: cd tests && go test -v
```
