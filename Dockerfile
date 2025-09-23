# 多阶段构建 - 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gotify2bark .

# 运行阶段 - 使用最小镜像
FROM alpine:latest

# 安装ca-certificates用于HTTPS请求
RUN apk --no-cache add ca-certificates tzdata

# 创建非root用户和应用目录
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    mkdir -p /app && \
    chown appuser:appgroup /app

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/gotify2bark .

# 复制启动脚本
COPY entrypoint.sh .

# 设置执行权限并更改文件所有者
RUN chmod +x entrypoint.sh gotify2bark && \
    chown appuser:appgroup gotify2bark entrypoint.sh

# 切换到非root用户
USER appuser

# 暴露端口（如果需要的话）
# EXPOSE 8080

# 设置环境变量
ENV GOTIFY_HOST=""
ENV GOTIFY_TOKEN=""
ENV BARK_URL=""
ENV ICON_URL=""

# 运行应用
ENTRYPOINT ["./entrypoint.sh"]