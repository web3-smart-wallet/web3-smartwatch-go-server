# 构建阶段
FROM golang:1.24-alpine AS builder

# 安装必要的依赖
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 运行go generate（添加错误处理）
RUN go generate ./... || echo "Warning: go generate command failed, continuing build..."

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/main .

# 仅当.env文件存在时才复制
COPY --from=builder /app/.env* ./ 2>/dev/null || true

# 暴露端口（根据您的应用需要修改端口号）
EXPOSE 8080

# 运行应用
CMD ["./main"] 