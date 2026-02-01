#!/bin/sh

# 启动脚本：将环境变量转换为命令行参数

# 检查必需的环境变量
if [ -z "$GOTIFY_HOST" ] || [ -z "$GOTIFY_TOKEN" ] || [ -z "$BARK_URL" ]; then
    echo "错误：缺少必需的环境变量"
    echo "请设置以下环境变量："
    echo "  GOTIFY_HOST - Gotify服务器地址 (例如: https://gotify.example.com)"
    echo "  GOTIFY_TOKEN - Gotify客户端令牌"
    echo "  BARK_URL - Bark推送URL"
    echo "  ICON_URL - 图标URL (可选，默认使用Bark图标)"
    exit 1
fi

# 构建命令行参数
ARGS="-host $GOTIFY_HOST -token $GOTIFY_TOKEN -target $BARK_URL"

# 如果设置了图标URL，添加到参数中
if [ -n "$ICON_URL" ]; then
    ARGS="$ARGS -icon $ICON_URL"
fi

# 加密相关参数
if [ -n "$BARK_ENCRYPTION_KEY" ]; then
    ARGS="$ARGS -aes-key $BARK_ENCRYPTION_KEY"
fi

if [ -n "$BARK_ENCRYPTION_IV" ]; then
    ARGS="$ARGS -aes-iv $BARK_ENCRYPTION_IV"
fi

# Basic Auth 相关参数
if [ -n "$BARK_SERVER_USER" ]; then
    ARGS="$ARGS -bark-user $BARK_SERVER_USER"
fi

if [ -n "$BARK_SERVER_PASSWORD" ]; then
    ARGS="$ARGS -bark-password $BARK_SERVER_PASSWORD"
fi

echo "启动 Gotify-to-Bark 转发器..."
echo "Gotify服务器: $GOTIFY_HOST"
echo "Bark URL: $BARK_URL"
if [ -n "$ICON_URL" ]; then
    echo "图标URL: $ICON_URL"
fi

# 执行应用程序
exec ./gotify2bark $ARGS