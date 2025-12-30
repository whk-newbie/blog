# 构建阶段
FROM golang:1.25-alpine AS builder

# 设置 Go 代理（加速下载）
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# 配置 Alpine 镜像源（加速 apk 下载）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 设置工作目录
WORKDIR /build

# 安装必要的构建工具（合并RUN减少层数）
RUN apk add --no-cache git make wget

# 复制go mod文件（利用缓存层）
COPY go.mod go.sum ./

# 下载依赖（利用缓存层）
RUN go mod download

# 安装swag工具（利用缓存层）
RUN go install github.com/swaggo/swag/cmd/swag@latest

# 复制源代码
COPY . .

# 生成Swagger文档
RUN swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal --outputFileName doc.json

# 构建应用（优化构建参数，减小二进制体积）
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -trimpath \
    -o blog-server ./cmd/server

# 运行阶段
FROM alpine:latest

# 配置 Alpine 镜像源（加速 apk 下载）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装必要的运行时依赖（合并RUN减少层数）
RUN apk --no-cache add ca-certificates tzdata wget && \
    rm -rf /var/cache/apk/*

# 设置时区
ENV TZ=Asia/Shanghai

# 创建应用用户和目录（合并操作）
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    mkdir -p /app/config /app/uploads /app/backups /app/logs && \
    chown -R appuser:appgroup /app

# 设置工作目录
WORKDIR /app

# 从构建阶段复制文件（合并COPY减少层数）
COPY --from=builder --chown=appuser:appgroup /build/blog-server ./
COPY --from=builder --chown=appuser:appgroup /build/docs ./docs
COPY --from=builder --chown=appuser:appgroup /build/migrations ./migrations

# 切换到应用用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查（优化检查命令）
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./blog-server"]

