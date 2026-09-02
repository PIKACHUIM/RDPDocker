# 快速开始

本文将在 5 分钟内带你完成 RDDocker 管理平台的部署，并创建第一个桌面容器。

## 前置条件

- 一台 Linux 主机（Debian / Ubuntu / Fedora / Arch，x86_64 或 arm64）
- 已安装并启动 Docker（或 Podman / LXD）
- root 权限

## 一键安装

```bash
curl -fsSL https://raw.githubusercontent.com/PIKACHUIM/RDDocker/main/scripts/install/install-all.sh | sudo bash
```

脚本会依次完成：

1. 检测系统架构与包管理器
2. 创建 `/opt/deskcli/{conf,data,web/dist}` 目录
3. 下载 `deskcli` 二进制到 `/usr/local/bin/deskcli`
4. 部署 Web 前端静态资源
5. 生成 `config.yaml`，其中 `jwt_secret` 与 `token` 为随机值
6. 注册并启动 `deskcli` systemd 服务

安装完成后，终端会打印访问地址。

## 首次登录

浏览器打开 `http://<你的服务器IP>:8080`，使用初始账号登录：

| 项目 | 值 |
|------|-----|
| 用户名 | `admin` |
| 密码 | `admin123` |

::: warning 立即修改密码
初始密码是固定的，任何能访问该端口的人都能登录。登录后请立刻在 **系统设置 → 修改密码** 中更换。
:::

初始管理员账号是在服务首次启动、用户表为空时自动创建的。若你希望自行设定账号，可在启动服务前通过 `POST /api/v1/auth/register` 创建第一个用户——该接口仅在系统中尚无任何用户时可用。

## 创建第一个容器

1. 左侧导航进入 **模板市场**
2. 从 30 余个预置模板中挑选，例如 `Debian 12 · Xfce 4`
3. 点击卡片上的 **一键部署**
4. 填写容器名称（如 `my-desktop`），确认部署

平台会自动从 `port_range`（默认 `50000-60000`）中分配 4 个未占用端口，分别映射到容器的 SSH(22)、RDP(3389)、NX(4000)、VNC(5900)。

部署完成后跳转到 **容器管理** 页，即可看到新容器的状态与端口。

## 连接桌面

在容器卡片上查看分配到的端口，然后用对应客户端连接：

```bash
# SSH，假设分配到的外部端口是 50000
ssh root@<服务器IP> -p 50000
```

RDP 客户端（如 Windows 远程桌面、Remmina）连接 `<服务器IP>:<RDP端口>`；VNC 客户端连接 `<服务器IP>:<VNC端口>`。

也可以直接在 Web 界面的 **终端** 页选择容器，获得浏览器内的交互式 Shell，无需额外客户端。

## 下一步

- [Web 管理面板](/guide/management) —— 各页面功能详解
- [安装方法](/guide/installation) —— 手动安装与引擎切换
- [常见问题](/guide/faq) —— 端口冲突、拉取缓慢等排查
