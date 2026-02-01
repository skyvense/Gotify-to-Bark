# Gotify to Bark Forwarder

一个将 Gotify 消息转发到 Bark 服务器的工具。这个程序会监听 Gotify 服务器的 WebSocket 连接，并将接收到的消息实时转发到 Bark 服务器。

## 功能特点

- 实时转发 Gotify 消息到 Bark
- 自动重连机制
- 详细的日志记录
- 支持自定义消息格式
- 支持 HTTPS/WSS 连接
- 可自定义通知图标
- 支持 Bark 消息加密 (AES-128-CBC)
- 支持 Bark 服务器 Basic Auth 认证

## 安装

### 方式一：使用 Docker（推荐）

直接拉取预构建的 Docker 镜像：

```bash
# 从 Docker Hub 拉取
docker pull skyvense/gotify-to-bark:latest

# 或从 GitHub Container Registry 拉取
docker pull ghcr.io/skyvense/gotify-to-bark:latest
```

### 方式二：从源码构建

1. 确保已安装 Go 1.21 或更高版本
2. 克隆仓库：
   ```bash
   git clone https://github.com/skyvense/Gotify-to-Bark.git
   cd Gotify-to-Bark
   ```
3. 安装依赖：
   ```bash
   go mod download
   ```

## 使用方法

### Docker 方式（推荐）

使用 Docker 是最简单的部署方式：

```bash
# 使用 Docker Hub
docker run -d --name gotify2bark \
  -e GOTIFY_HOST="https://gotify.example.com" \
  -e GOTIFY_TOKEN="your-gotify-token" \
  -e BARK_URL="https://api.day.app/your-device-key" \
  -e ICON_URL="https://example.com/icon.png" \
  -e BARK_ENCRYPTION_KEY="your-aes-key-16-bytes" \
  -e BARK_ENCRYPTION_IV="your-aes-iv-16-bytes" \
  -e BARK_SERVER_USER="username" \
  -e BARK_SERVER_PASSWORD="password" \
  skyvense/gotify-to-bark:latest

# 或使用 GitHub Container Registry
docker run -d --name gotify2bark \
  -e GOTIFY_HOST="https://gotify.example.com" \
  -e GOTIFY_TOKEN="your-gotify-token" \
  -e BARK_URL="https://api.day.app/your-device-key" \
  -e ICON_URL="https://example.com/icon.png" \
  -e BARK_ENCRYPTION_KEY="your-aes-key-16-bytes" \
  -e BARK_ENCRYPTION_IV="your-aes-iv-16-bytes" \
  -e BARK_SERVER_USER="username" \
  -e BARK_SERVER_PASSWORD="password" \
  ghcr.io/skyvense/gotify-to-bark:latest
```

#### 环境变量说明：
- `GOTIFY_HOST`: Gotify 服务器地址（例如：https://gotify.example.com）
- `GOTIFY_TOKEN`: Gotify 客户端 token
- `BARK_URL`: Bark 服务器地址（例如：https://api.day.app/your-device-key）
- `ICON_URL`: 通知图标 URL（可选，默认为 https://day.app/assets/images/avatar.jpg）
- `BARK_ENCRYPTION_KEY`: Bark 加密密钥 (AES-128-CBC, 16字节，可选)
- `BARK_ENCRYPTION_IV`: Bark 加密 IV (AES-128-CBC, 16字节, 可选)
- `BARK_SERVER_USER`: Bark 服务器 Basic Auth 用户名 (可选)
- `BARK_SERVER_PASSWORD`: Bark 服务器 Basic Auth 密码 (可选)

#### Docker Compose

创建 `docker-compose.yml` 文件：

```yaml
version: '3.8'
services:
  gotify2bark:
    image: skyvense/gotify-to-bark:latest
    container_name: gotify2bark
    environment:
      - GOTIFY_HOST=https://gotify.example.com
      - GOTIFY_TOKEN=your-gotify-token
      - BARK_URL=https://api.day.app/your-device-key
      - ICON_URL=https://example.com/icon.png
      # 可选配置
      - BARK_ENCRYPTION_KEY=your-aes-key-16-bytes
      - BARK_ENCRYPTION_IV=your-aes-iv-16-bytes
      - BARK_SERVER_USER=username
      - BARK_SERVER_PASSWORD=password
    restart: unless-stopped
```

然后运行：
```bash
docker-compose up -d
```

### 命令行参数

```bash
go run main.go -host <gotify-host> -token <gotify-token> -target <bark-url> [-icon <icon-url>] [-aes-key <key>] [-aes-iv <iv>] [-bark-user <user>] [-bark-password <password>]
```

参数说明：
- `-host`: Gotify 服务器地址（例如：https://gotify.example.com）
- `-token`: Gotify 客户端 token
- `-target`: Bark 服务器地址（例如：https://api.day.app/your-device-key）
- `-icon`: 通知图标 URL（可选，默认为 https://day.app/assets/images/avatar.jpg）
- `-aes-key`: Bark 加密密钥 (AES-128-CBC, 16字节)
- `-aes-iv`: Bark 加密 IV (AES-128-CBC, 16字节, 可选)
- `-bark-user`: Bark 服务器 Basic Auth 用户名
- `-bark-password`: Bark 服务器 Basic Auth 密码

### 示例

```bash
# 使用默认图标
go run main.go -host https://gotify.example.com -token ABC123 -target https://api.day.app/your-device-key

# 使用自定义图标
go run main.go -host https://gotify.example.com -token ABC123 -target https://api.day.app/your-device-key -icon https://example.com/icon.png

# 使用加密和 Basic Auth
go run main.go \
  -host https://gotify.example.com \
  -token ABC123 \
  -target https://api.day.app/your-device-key \
  -aes-key "1234567890123456" \
  -bark-user "admin" \
  -bark-password "secret"
```

### 构建可执行文件

```bash
go build -o gotify2bark
```

然后运行：
```bash
# 基本用法
./gotify2bark -host https://gotify.example.com -token ABC123 -target https://api.day.app/your-device-key
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
    "url": "Gotify服务器地址",
    "ciphertext": "加密后的内容（如果启用了加密）",
    "iv": "加密IV（如果启用了加密）"
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
