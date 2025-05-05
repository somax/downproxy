# DownProxy

DownProxy 是一个简单的文件下载代理服务器，支持 HTTP/HTTPS 协议的文件下载。它可以帮助你绕过网络限制，支持断点续传，并提供文件名解析等功能。

## 特性

- 支持 HTTP/HTTPS 协议
- 支持断点续传
- 自动解析文件名
- 支持大文件下载（默认最大 1GB）
- 可配置下载超时
- 支持多架构部署（amd64/arm64）
- 提供 Docker 容器支持

## 快速开始

### 直接运行

```bash
# 编译
go build -o downproxy

# 运行（默认端口 9527）
./downproxy

# 指定端口运行
./downproxy -port 8080
```

## Docker 运行
```bash
# 使用 Docker Compose
docker-compose up -d

# 或直接使用 Docker
docker run -d -p 9527:9527 yourusername/downproxy
# 使用 Docker Composedocker-compose up -d# 或直接使用 Dockerdocker run -d -p 9527:9527 yourusername/downproxy
```

## 使用方法
下载文件：

```plaintext
http://localhost:9527/download?url=https://example.com/file.zip
```

## 配置说明
主要配置参数：

- `port`: 服务监听端口（默认：9527）
- `maxContentSize`: 最大下载文件大小（默认：1GB）
- `timeout`: 下载超时时间（默认：300秒）
- `userAgent`: 请求 User-Agent
- `allowedProtocols`: 允许的协议（默认：http, https）

## 构建
### 多架构二进制构建
```bash
运行
# 运行多架构构建脚本
./scripts/build-go-multiarch.sh
# 运行多架构构建脚本./scripts/build-go-multiarch.sh
```

### Docker 多架构镜像构建

```bash
# 运行 Docker 多架构构建脚本
./scripts/build-multiarch.sh
# 运行 Docker 多架构构建脚本./scripts/build-multiarch.sh
```

## 系统服务安装
在 Linux 系统上安装为系统服务：

```bash
运行
sudo ./scripts/install-service.sh
```

## 许可证
MIT License

## 贡献
欢迎提交 Issue 和 Pull Request！
