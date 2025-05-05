#!/bin/bash

# 设置变量
APP_NAME="downproxy"
VERSION="1.0.0"
OUTPUT_DIR="./build"

# 创建输出目录
mkdir -p ${OUTPUT_DIR}
echo "创建输出目录: ${OUTPUT_DIR}"

# 定义要构建的平台
PLATFORMS=(
  "linux/amd64"
  # "linux/386"
  # "linux/arm64"
  # "linux/arm/v7"
  # "darwin/amd64"
  "darwin/arm64"
  # "windows/amd64"
  # "windows/386"
)

# 构建函数
build() {
  local platform=$1
  local os=$(echo ${platform} | cut -d/ -f1)
  local arch=$(echo ${platform} | cut -d/ -f2)
  local arm=""
  
  # 处理 ARM 版本
  if [[ ${platform} == *"arm/v"* ]]; then
    arm=$(echo ${platform} | cut -d/ -f3)
    arm=${arm#v}
    arch="arm"
  fi
  
  # 设置输出文件名
  local output="${OUTPUT_DIR}/${APP_NAME}-${os}-${arch}"
  if [ ! -z "${arm}" ]; then
    output="${output}-v${arm}"
  fi
  
  # Windows 平台添加 .exe 后缀
  if [ "${os}" == "windows" ]; then
    output="${output}.exe"
  fi
  
  echo "正在构建 ${platform}..."
  
  # 设置环境变量并编译
  GOOS=${os} GOARCH=${arch} GOARM=${arm} CGO_ENABLED=0 go build -ldflags="-s -w" -o ${output} .
  
  if [ $? -eq 0 ]; then
    echo "✅ 构建成功: ${output}"
  else
    echo "❌ 构建失败: ${platform}"
    return 1
  fi
  
  # 为非 Windows 平台添加执行权限
  if [ "${os}" != "windows" ]; then
    chmod +x ${output}
  fi
  
  return 0
}

# 开始构建
echo "开始为 ${APP_NAME} v${VERSION} 构建多架构二进制文件..."

# 记录成功和失败的平台
success_platforms=()
failed_platforms=()

# 遍历所有平台进行构建
for platform in "${PLATFORMS[@]}"; do
  build ${platform}
  if [ $? -eq 0 ]; then
    success_platforms+=("${platform}")
  else
    failed_platforms+=("${platform}")
  fi
done

# 输出构建结果摘要
echo ""
echo "构建完成!"
echo "成功平台 (${#success_platforms[@]}):"
for platform in "${success_platforms[@]}"; do
  echo "  - ${platform}"
done

if [ ${#failed_platforms[@]} -gt 0 ]; then
  echo "失败平台 (${#failed_platforms[@]}):"
  for platform in "${failed_platforms[@]}"; do
    echo "  - ${platform}"
  done
  exit 1
fi

echo ""
echo "所有二进制文件已保存到 ${OUTPUT_DIR} 目录"