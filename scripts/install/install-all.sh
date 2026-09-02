#!/bin/bash
# RDDocker 管理平台一键安装脚本
# 支持 Debian/Ubuntu/Fedora/Arch Linux

set -e

INSTALL_DIR="/opt/deskcli"
BIN_PATH="/usr/local/bin/deskcli"
SERVICE_FILE="/etc/systemd/system/deskcli.service"
GITHUB_REPO="PIKACHUIM/RDDocker"
VERSION="latest"

echo "=========================================="
echo "  RDDocker 管理平台安装脚本"
echo "=========================================="
echo

# 检测系统架构
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "不支持的架构: $ARCH"
        exit 1
        ;;
esac

echo "[1/7] 检测系统架构: $ARCH"

# 检测包管理器
if command -v apt-get &> /dev/null; then
    PKG_INSTALL="apt-get install -y"
    PKG_UPDATE="apt-get update"
elif command -v dnf &> /dev/null; then
    PKG_INSTALL="dnf install -y"
    PKG_UPDATE="dnf check-update || true"
elif command -v yum &> /dev/null; then
    PKG_INSTALL="yum install -y"
    PKG_UPDATE="yum check-update || true"
elif command -v pacman &> /dev/null; then
    PKG_INSTALL="pacman -S --noconfirm"
    PKG_UPDATE="pacman -Sy"
else
    echo "未找到支持的包管理器"
    exit 1
fi

echo "[2/7] 安装依赖..."
$PKG_UPDATE
$PKG_INSTALL curl wget tar

# 创建安装目录
echo "[3/7] 创建目录结构..."
mkdir -p $INSTALL_DIR/{conf,data,web/dist}

# 下载 deskcli 二进制（这里需要实际的 release URL）
echo "[4/7] 下载 deskcli 二进制..."
DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$VERSION/deskcli-linux-$ARCH"
if ! wget -q --show-progress $DOWNLOAD_URL -O $BIN_PATH; then
    echo "下载失败，尝试从备用源..."
    # 备用：从当前构建复制
    if [ -f "./deskcli/deskcli" ]; then
        cp ./deskcli/deskcli $BIN_PATH
    else
        echo "错误：无法获取 deskcli 二进制文件"
        exit 1
    fi
fi
chmod +x $BIN_PATH

# 下载前端静态文件（或从本地复制）
echo "[5/7] 部署前端资源..."
if [ -d "./deskcli/web/dist" ]; then
    cp -r ./deskcli/web/dist/* $INSTALL_DIR/web/dist/
else
    echo "警告：未找到前端构建产物，仅安装 API 服务"
fi

# 生成配置文件
echo "[6/7] 生成配置文件..."
JWT_SECRET=$(openssl rand -hex 32 2>/dev/null || head -c 32 /dev/urandom | base64)
ADMIN_PASSWORD=$(openssl rand -base64 12 2>/dev/null || head -c 12 /dev/urandom | base64)

cat > $INSTALL_DIR/conf/config.yaml <<EOF
# RDDocker 管理平台配置
engine: docker
listen_addr: 0.0.0.0
port: 8080
data_dir: $INSTALL_DIR/data
port_range: 50000-60000
allow_exec: true
static_dir: $INSTALL_DIR/web/dist
jwt_secret: $JWT_SECRET
token: $(openssl rand -hex 16 2>/dev/null || head -c 16 /dev/urandom | base64)
EOF

chmod 600 $INSTALL_DIR/conf/config.yaml

# 创建 systemd 服务
echo "[7/7] 配置 systemd 服务..."
cat > $SERVICE_FILE <<EOF
[Unit]
Description=RDDocker Management Platform
After=network.target docker.service

[Service]
Type=simple
User=root
ExecStart=$BIN_PATH serve
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable deskcli
systemctl start deskcli

echo
echo "=========================================="
echo "  安装完成！"
echo "=========================================="
echo
echo "访问地址: http://$(hostname -I | awk '{print $1}'):8080"
echo "默认用户名: admin"
echo "默认密码: admin123 (首次登录后请立即修改)"
echo
echo "管理命令:"
echo "  systemctl status deskcli   # 查看状态"
echo "  systemctl restart deskcli  # 重启服务"
echo "  systemctl logs deskcli     # 查看日志"
echo "  deskcli --help             # CLI 帮助"
echo
