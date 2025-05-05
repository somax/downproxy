#!/bin/bash

# 检查是否以 root 权限运行
if [ "$EUID" -ne 0 ]; then
  echo "请使用 sudo 运行此脚本"
  exit 1
fi

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
EXECUTABLE="$SCRIPT_DIR/downproxy"

# 检查可执行文件是否存在
if [ ! -f "$EXECUTABLE" ]; then
  echo "错误: 在 $EXECUTABLE 找不到可执行文件"
  exit 1
fi

# 确保可执行文件有执行权限
chmod +x "$EXECUTABLE"

# 创建 systemd 服务文件
cat > /tmp/downproxy.service << EOF
[Unit]
Description=Download Proxy Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$SCRIPT_DIR
ExecStart=$EXECUTABLE
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# 复制服务文件到 systemd 目录
cp /tmp/downproxy.service /etc/systemd/system/
rm /tmp/downproxy.service

echo "服务文件已创建并安装到 /etc/systemd/system/downproxy.service"

# 重新加载 systemd 配置
systemctl daemon-reload

echo "是否立即启用并启动服务? (y/n)"
read -r answer
if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
  systemctl enable downproxy.service
  systemctl start downproxy.service
  echo "服务已启用并启动"
  systemctl status downproxy.service
else
  echo "您可以稍后使用以下命令启用并启动服务:"
  echo "  sudo systemctl enable downproxy.service"
  echo "  sudo systemctl start downproxy.service"
fi

echo "安装完成!"