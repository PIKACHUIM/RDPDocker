# 安装方法

## 方式一：直接使用预构建镜像（推荐）

无需安装任何额外工具，直接使用 `docker run` 拉取并启动镜像。

```bash
# 示例：启动 Debian 12 + Xfce4 桌面
docker run -d \
  --name my-desktop \
  -p 4000:4000 \
  -p 5900:5900 \
  -p 3389:3389 \
  -p 2222:22 \
  pikachuim/debian:12-xfce4l
```

启动后通过 NoMachine（端口 4000）、VNC（端口 5900）或 XRDP（端口 3389）连接。

默认凭据请参考镜像说明，建议首次启动后立即修改密码：

```bash
docker exec my-desktop passwd user
```

## 方式二：安装 deskcli 管理工具

`deskcli` 提供统一的 CLI 和 REST API 来管理多种容器后端。

### 一键安装

```bash
curl -fsSL https://raw.githubusercontent.com/PIKACHUIM/RDDocker/master/deskcli/install.sh | bash
```

### 验证安装

```bash
deskcli --version
```

### 初始配置

```bash
# 设置容器引擎（docker / podman / lxc / lxd）
deskcli configs set engine docker

# 设置 API 访问令牌（用于 REST API 认证）
deskcli configs set token your_secret_token

# 设置 API 服务端口（默认 8080）
deskcli configs set port 8080
```

### 创建第一个桌面容器

```bash
deskcli desk new pikachuim/debian:12-xfce4l \
  --name my-desktop \
  --port 4000:4000 \
  --port 5900:5900 \
  --port 2222:22
```

## 方式三：部署 Web 管理平台

在 CLI 之外额外提供 React 管理界面，与 REST API 共用同一端口。若只想尽快跑起来，请直接看[快速开始](/guide/quick-start)。

### 一键安装

```bash
curl -fsSL https://raw.githubusercontent.com/PIKACHUIM/RDDocker/main/scripts/install/install-all.sh | sudo bash
```

脚本会创建目录、部署二进制与前端资源、生成带随机密钥的配置，并注册 systemd 服务。

### 手动安装

若不想执行远程脚本，可以逐步操作。先准备目录：

```bash
sudo mkdir -p /opt/deskcli/{conf,data,web/dist}
```

构建后端（需要 Go 1.21+）：

```bash
cd deskcli
go build -o /usr/local/bin/deskcli .
```

构建前端（需要 Node 18+），产物会输出到 `deskcli/web/dist`：

```bash
cd web
npm install
npm run build
sudo cp -r ../deskcli/web/dist/* /opt/deskcli/web/dist/
```

写入 `/opt/deskcli/conf/config.yaml`：

```yaml
engine: docker
listen_addr: 0.0.0.0
port: 8080
data_dir: /opt/deskcli/data
port_range: 50000-60000
allow_exec: true
static_dir: /opt/deskcli/web/dist
jwt_secret: <随机 32 字节十六进制串>
token: <随机 16 字节十六进制串>
```

两个密钥可以这样生成：

```bash
openssl rand -hex 32   # jwt_secret
openssl rand -hex 16   # token
```

::: warning 务必更换默认密钥
`jwt_secret` 用于签发登录 Token，若使用可预测的值，攻击者可以自行伪造任意用户身份。配置文件权限建议设为 `600`。
:::

### 注册 systemd 服务

写入 `/etc/systemd/system/deskcli.service`：

```ini
[Unit]
Description=RDDocker Management Platform
After=network.target docker.service

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/deskcli serve
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

启用并启动：

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now deskcli
sudo systemctl status deskcli
```

首次启动时若用户表为空，系统会自动创建管理员账号 `admin` / `admin123`，请登录后立即修改。

### 切换容器引擎

修改配置中的 `engine` 字段后重启服务即可。各引擎的准备工作：

```bash
# Docker
sudo systemctl enable --now docker

# Podman（无守护进程，安装即可）
sudo apt install -y podman

# LXD
sudo snap install lxd && sudo lxd init --minimal
```

注意不同引擎支持的功能范围有差异——原生 LXC 不支持镜像列表与删除，LXC/LXD 均不支持 Web 终端，详见[常见问题](/guide/faq)中的对照表。

### 无 root 模式

以非 root 用户运行时，需要把数据目录改到用户可写的位置，并确保当前用户能访问容器引擎：

```yaml
data_dir: ~/.local/share/deskcli
static_dir: ~/.local/share/deskcli/web/dist
```

Docker 需将用户加入 `docker` 组（重新登录后生效）：

```bash
sudo usermod -aG docker $USER
```

Podman 原生支持 rootless，无需额外授权。此外端口转发依赖 iptables 规则，非 root 环境下该功能不可用，请改用引擎自身的端口映射。

## 方式四：从源码构建镜像

适合需要自定义镜像内容的用户。

### 克隆仓库

```bash
git clone https://github.com/PIKACHUIM/RDDocker.git
cd RDDocker
```

### 运行构建脚本

```bash
# 构建所有镜像
bash builds.sh

# 或单独构建某个镜像（示例：Debian 12 GNOME）
docker build \
  -f dockers/debian/desktops/gnome3 \
  -t my-debian-gnome .
```

### 使用 build-arg 自定义软件

```bash
docker build \
  --build-arg INSTALL_FIREFOX=true \
  --build-arg INSTALL_VSCODE=true \
  -f dockers/debian/desktops/gnome3 \
  -t my-debian-gnome-custom .
```

## 构建本文档

```bash
cd docpage
npm install
npm run docs:dev     # 本地预览（http://localhost:5173）
npm run docs:build   # 构建静态文件到 .vitepress/dist/
```
