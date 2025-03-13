# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# 暴露端口（根据您的应用需要修改端口号）
EXPOSE 8080

# 运行应用
CMD ["./main"] 