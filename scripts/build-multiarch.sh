#!/bin/bash

# 设置变量
IMAGE_NAME="yourusername/downproxy"
VERSION="1.0.0"

# 确保 Docker Buildx 可用
docker buildx create --name multiarch-builder --use || true

# 构建并推送多架构镜像
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 \
  -t ${IMAGE_NAME}:${VERSION} \
  -t ${IMAGE_NAME}:latest \
  --push .

echo "多架构镜像构建完成并推送到 Docker Hub"