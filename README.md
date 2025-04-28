# Gotify to Bark Forwarder

一个将 Gotify 消息转发到 Bark 服务器的工具。这个程序会监听 Gotify 服务器的 WebSocket 连接，并将接收到的消息实时转发到 Bark 服务器。

## 功能特点

- 实时转发 Gotify 消息到 Bark
- 自动重连机制
- 详细的日志记录
- 支持自定义消息格式
- 支持 HTTPS/WSS 连接
- 可自定义通知图标

## 安装

1. 确保已安装 Go 1.21 或更高版本
2. 克隆仓库：
   ```bash
   git clone https://github.com/yourusername/gotify2bark.git
   cd gotify2bark
   ```
3. 安装依赖：
   ```bash
   go mod download
   ```

## 使用方法

### 命令行参数

```bash
go run main.go -host <gotify-host> -token <gotify-token> -target <bark-url> [-icon <icon-url>]
```

参数说明：
- `-host`: Gotify 服务器地址（例如：https://gotify.example.com）
- `-token`: Gotify 客户端 token
- `-target`: Bark 服务器地址（例如：https://api.day.app/your-device-key）
- `-icon`: 通知图标 URL（可选，默认为 https://day.app/assets/images/avatar.jpg）

### 示例

```bash
# 使用默认图标
go run main.go -host https://gotify.example.com -token ABC123 -target https://api.day.app/your-device-key

# 使用自定义图标
go run main.go -host https://gotify.example.com -token ABC123 -target https://api.day.app/your-device-key -icon https://example.com/icon.png
```

### 构建可执行文件

```bash
go build -o gotify2bark
```

然后运行：
```bash
# 使用默认图标
./gotify2bark -host https://gotify.example.com -token ABC123 -target https://api.day.app/your-device-key

# 使用自定义图标
./gotify2bark -host https://gotify.example.com -token ABC123 -target https://api.day.app/your-device-key -icon https://example.com/icon.png
```

## 消息格式

程序会将 Gotify 消息转换为以下 Bark 格式：

```json
{
    "title": "Gotify消息标题",
    "body": "Gotify消息内容",
    "badge": 1,
    "sound": "minuet",
    "group": "Gotify",
    "icon": "自定义图标URL",
    "url": "Gotify服务器地址"
}
```

## 日志输出

程序会输出详细的日志信息，包括：
- WebSocket 连接状态
- 接收到的消息内容
- 转发状态
- 错误信息（如果有）

## 开发

### 调试

使用 VS Code 调试配置（.vscode/launch.json）：
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Gotify Forwarder",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "args": [
                "-host",
                "https://gotify.example.com",
                "-token",
                "your-token",
                "-target",
                "https://api.day.app/your-device-key",
                "-icon",
                "https://example.com/icon.png"
            ],
            "env": {},
            "showLog": true
        }
    ]
}
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！ 