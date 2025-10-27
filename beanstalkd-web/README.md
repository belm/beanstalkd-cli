# Beanstalkd Web UI 管理后台

[🇨🇳 中文版](README.md) | [🇺🇸 English](README-EN.md)

一个现代化、美观、易用的 Beanstalkd Web 管理界面。

## ✨ 特性

- 🎨 **现代化设计** - 使用 Tailwind CSS，响应式布局
- 📊 **实时统计** - 服务器和 Tube 的实时统计信息
- 🔄 **自动刷新** - 每10秒自动更新数据
- 🎯 **完整功能** - 支持所有常用 Beanstalkd 操作
- 🌈 **友好界面** - 直观的操作流程和视觉反馈
- 📱 **移动端适配** - 完美支持移动设备

## 🚀 快速开始

### 1. 确保 Beanstalkd 服务运行

```bash
# 默认连接到 127.0.0.1:11300
```

### 2. 启动 Web 服务器

```bash
cd beanstalkd-web
go run server.go
```

### 3. 访问管理后台

打开浏览器访问: **http://localhost:8080**

## 📋 功能说明

### 1. Tubes 管理
- 查看所有 Tubes 列表
- 查看每个 Tube 的详细统计
- 实时监控各 Tube 的任务状态

### 2. 任务管理
- 选择 Tube 查看任务统计
- 查看就绪、预留、埋葬、延迟任务数量
- 快速切换到操作中心进行操作

### 3. 操作中心

#### 插入任务
- 选择 Tube
- 输入任务数据
- 设置优先级和延迟
- 一键插入

#### 预留任务
- 从指定 Tube 预留任务
- 设置超时时间
- 查看预留的任务详情

#### 删除任务
- 输入任务 ID
- 确认删除

#### 踢出任务
- 批量踢出被埋葬或延迟的任务
- 设置踢出数量

### 4. 服务器统计
- 查看完整的服务器统计信息
- 所有指标一目了然
- 实时更新

## 🎨 界面预览

### 主页面
- **顶部统计卡片**: 快速查看关键指标
- **标签页导航**: Tubes、任务、操作、统计
- **实时刷新**: 自动更新数据

### Tubes 管理
- **卡片式展示**: 每个 Tube 的详细信息
- **颜色编码**: 不同状态使用不同颜色
- **快速操作**: 点击查看详情

### 操作中心
- **四大功能区**: 插入、预留、删除、踢出
- **表单验证**: 友好的输入提示
- **即时反馈**: 操作结果实时显示

## 🔧 配置

### 方式一：命令行参数（推荐）

```bash
# 自定义 Beanstalkd 地址
go run server.go -beanstalkd 192.168.1.100:11300

# 自定义 Web 端口
go run server.go -port 9090

# 同时自定义
go run server.go -beanstalkd 192.168.1.100:11300 -port 9090

# 查看帮助
go run server.go -h
```

### 方式二：环境变量

```bash
# 设置环境变量
export BEANSTALKD_HOST=192.168.1.100:11300
export WEB_PORT=9090

# 启动服务器
go run server.go
```

### 方式三：使用默认值

```bash
# 直接启动（使用默认配置）
# Beanstalkd: 127.0.0.1:11300
# Web 端口: 8080
go run server.go
```

### 配置优先级

**命令行参数 > 环境变量 > 默认值**

## 📦 技术栈

### 前端
- **HTML5** - 结构
- **Tailwind CSS** (CDN) - 样式
- **Font Awesome** (CDN) - 图标
- **原生 JavaScript** - 交互逻辑

### 后端
- **Go** - 服务器语言
- **go-beanstalk** - Beanstalkd 客户端库
- **net/http** - HTTP 服务器

## 🎯 API 端点

### GET /api/stats
获取服务器统计信息

**响应**:
```json
{
  "stats": {
    "current-jobs-ready": "5",
    "current-jobs-reserved": "2",
    ...
  }
}
```

### GET /api/tubes
获取所有 Tubes

**响应**:
```json
{
  "tubes": ["default", "emails", "orders"]
}
```

### GET /api/tubes/{name}/stats
获取指定 Tube 的统计

**响应**:
```json
{
  "stats": {
    "current-jobs-ready": "3",
    "current-jobs-reserved": "1",
    ...
  }
}
```

### POST /api/put
插入任务

**请求**:
```json
{
  "tube": "default",
  "data": "task data",
  "priority": 1024,
  "delay": 0
}
```

**响应**:
```json
{
  "job_id": 123
}
```

### POST /api/reserve
预留任务

**请求**:
```json
{
  "tube": "default",
  "timeout": 5
}
```

**响应**:
```json
{
  "job_id": 123,
  "data": "task data"
}
```

### POST /api/delete
删除任务

**请求**:
```json
{
  "job_id": 123
}
```

### POST /api/kick
踢出任务

**请求**:
```json
{
  "tube": "default",
  "bound": 10
}
```

**响应**:
```json
{
  "kicked": 5
}
```

## 🔐 安全提示

⚠️ **重要**: 此管理后台没有内置认证机制，请注意：

1. 不要暴露到公网
2. 仅在受信任的网络中使用
3. 如需公网访问，请配置反向代理和认证
4. 生产环境建议使用 HTTPS

## 🚀 生产部署

### 使用 systemd

创建服务文件 `/etc/systemd/system/beanstalkd-web.service`:

```ini
[Unit]
Description=Beanstalkd Web UI
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/beanstalkd-cli/beanstalkd-web
ExecStart=/usr/local/go/bin/go run server.go
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动服务:
```bash
sudo systemctl enable beanstalkd-web
sudo systemctl start beanstalkd-web
```

### 使用 Nginx 反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 📝 开发指南

### 添加新功能

1. **添加 API 端点** - 在 `server.go` 中添加处理函数
2. **更新前端** - 在 `app.js` 中添加相应的 JavaScript 函数
3. **更新 UI** - 在 `index.html` 中添加界面元素

### 调试

启用详细日志:
```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
