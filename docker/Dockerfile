FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum 文件（如果有）
COPY go.mod ./
# 如果有 go.sum 文件，取消下面这行的注释
# COPY go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o downproxy .

# 使用轻量级的基础镜像
FROM alpine:latest

WORKDIR /root/

# 安装 ca-certificates，用于 HTTPS 请求
RUN apk --no-cache add ca-certificates

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/downproxy .

# 暴露端口
EXPOSE 9527

# 运行应用
ENTRYPOINT ["./downproxy"]