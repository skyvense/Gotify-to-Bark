# Docker 使用指南

本文档介绍如何使用 Docker 运行 Gotify-to-Bark 应用。

## 🐳 Docker 镜像

### 官方镜像源

- **Docker Hub**: `skyvense/gotify-to-bark`
- **GitHub Container Registry**: `ghcr.io/skyvense/gotify-to-bark`

### 支持的架构

- `linux/amd64` (x86_64)
- `linux/arm64` (ARM64/AArch64)

## 🚀 快速开始

### 使用 Docker 运行

```bash
docker run -d --name gotify2bark \
  -e GOTIFY_HOST="https://your-gotify-server.com" \
  -e GOTIFY_TOKEN="your-gotify-token" \
  -e BARK_URL="https://api.day.app/your-bark-key" \
  -e ICON_URL="https://your-icon-url.com/icon.png" \
  skyvense/gotify-to-bark:latest
```

### 使用 Docker Compose

1. 创建 `.env` 文件：

```env
GOTIFY_HOST=https://your-gotify-server.com
GOTIFY_TOKEN=your-gotify-token
BARK_URL=https://api.day.app/your-bark-key
ICON_URL=https://your-icon-url.com/icon.png
```

2. 运行服务：

```bash
docker-compose up -d
```

## 📋 环境变量

| 变量名 | 描述 | 必需 | 示例 |
|--------|------|------|------|
| `GOTIFY_HOST` | Gotify 服务器地址 | ✅ | `https://gotify.example.com` |
| `GOTIFY_TOKEN` | Gotify 客户端令牌 | ✅ | `AbCdEfGhIjKlMnOp` |
| `BARK_URL` | Bark 推送 URL | ✅ | `https://api.day.app/your-key` |
| `ICON_URL` | 通知图标 URL | ❌ | `https://example.com/icon.png` |

## 🔧 高级配置

### 自定义网络

```bash
# 创建自定义网络
docker network create gotify-network

# 在自定义网络中运行
docker run -d --name gotify2bark \
  --network gotify-network \
  -e GOTIFY_HOST="https://your-gotify-server.com" \
  -e GOTIFY_TOKEN="your-gotify-token" \
  -e BARK_URL="https://api.day.app/your-bark-key" \
  skyvense/gotify-to-bark:latest
```

### 资源限制

```bash
docker run -d --name gotify2bark \
  --memory="128m" \
  --cpus="0.5" \
  -e GOTIFY_HOST="https://your-gotify-server.com" \
  -e GOTIFY_TOKEN="your-gotify-token" \
  -e BARK_URL="https://api.day.app/your-bark-key" \
  skyvense/gotify-to-bark:latest
```

### 重启策略

```bash
docker run -d --name gotify2bark \
  --restart unless-stopped \
  -e GOTIFY_HOST="https://your-gotify-server.com" \
  -e GOTIFY_TOKEN="your-gotify-token" \
  -e BARK_URL="https://api.day.app/your-bark-key" \
  skyvense/gotify-to-bark:latest
```

## 📊 监控和日志

### 查看日志

```bash
# 查看实时日志
docker logs -f gotify2bark

# 查看最近的日志
docker logs --tail 100 gotify2bark
```

### 容器状态

```bash
# 查看容器状态
docker ps

# 查看容器详细信息
docker inspect gotify2bark
```

## 🛠️ 本地构建

### 构建镜像

```bash
# 克隆仓库
git clone https://github.com/skyvense/Gotify-to-Bark.git
cd Gotify-to-Bark

# 构建镜像
docker build -t gotify-to-bark:local .
```

### 多架构构建

```bash
# 设置 buildx
docker buildx create --use

# 构建多架构镜像
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t gotify-to-bark:multi-arch \
  --push .
```

## 🔍 故障排除

### 常见问题

1. **连接失败**
   ```bash
   # 检查网络连接
   docker exec gotify2bark ping gotify-server.com
   ```

2. **权限问题**
   ```bash
   # 检查容器用户
   docker exec gotify2bark id
   ```

3. **环境变量未生效**
   ```bash
   # 检查环境变量
   docker exec gotify2bark env | grep GOTIFY
   ```

### 调试模式

```bash
# 以交互模式运行
docker run -it --rm \
  -e GOTIFY_HOST="https://your-gotify-server.com" \
  -e GOTIFY_TOKEN="your-gotify-token" \
  -e BARK_URL="https://api.day.app/your-bark-key" \
  skyvense/gotify-to-bark:latest
```

## 📝 版本管理

### 标签说明

- `latest`: 最新稳定版本
- `v1.0.0`: 特定版本号
- `v1.0`: 主要版本
- `v1`: 大版本

### 更新镜像

```bash
# 拉取最新镜像
docker pull skyvense/gotify-to-bark:latest

# 停止并删除旧容器
docker stop gotify2bark
docker rm gotify2bark

# 运行新容器
docker run -d --name gotify2bark \
  -e GOTIFY_HOST="https://your-gotify-server.com" \
  -e GOTIFY_TOKEN="your-gotify-token" \
  -e BARK_URL="https://api.day.app/your-bark-key" \
  skyvense/gotify-to-bark:latest
```

## 🔐 安全最佳实践

1. **使用非 root 用户**: 镜像已配置为使用非 root 用户运行
2. **最小权限**: 只暴露必要的端口和卷
3. **定期更新**: 保持镜像版本最新
4. **环境变量安全**: 使用 Docker secrets 或环境文件管理敏感信息

```bash
# 使用 Docker secrets (Docker Swarm)
echo "your-gotify-token" | docker secret create gotify_token -
```

## 📚 相关链接

- [项目主页](https://github.com/skyvense/Gotify-to-Bark)
- [Docker Hub](https://hub.docker.com/r/skyvense/gotify-to-bark)
- [GitHub Container Registry](https://github.com/skyvense/Gotify-to-Bark/pkgs/container/gotify-to-bark)
- [发布说明](https://github.com/skyvense/Gotify-to-Bark/releases)