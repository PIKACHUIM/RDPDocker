---
layout: home
hero:
  name: deskcli
  text: 开箱即用的桌面容器环境
  tagline: 基于 Docker/Podman/LXC，一键构建带远程桌面的 Linux 容器，并通过 Web 面板统一管理
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/quick-start
    - theme: alt
      text: Web 管理面板
      link: /guide/management
    - theme: alt
      text: 查看镜像
      link: /guide/images
features:
  - icon: 🖥️
    title: 多桌面环境
    details: 支持 GNOME、Xfce4、KDE Plasma、GXDeOS、Niri 等多种桌面环境，满足不同使用场景
  - icon: 🌐
    title: 多种远程访问
    details: 内置 NoMachine NX、x11vnc VNC、XRDP，开箱即用，无需额外配置
  - icon: 🏗️
    title: 多平台支持
    details: 支持 Docker / Podman / LXC / LXD，覆盖 x86_64 和 ARM64 架构
  - icon: 🎛️
    title: Web 管理面板
    details: 内置 React 控制台，涵盖仪表板、容器管理、模板市场与浏览器内终端，无需命令行
  - icon: 📦
    title: 模板一键部署
    details: 预置 30 余个发行版与桌面组合模板，自动分配 SSH/RDP/NX/VNC 端口，填个名字即可开跑
  - icon: ⚙️
    title: CLI 与 REST API
    details: 统一 CLI + JWT 认证的 REST API，支持远程控制、端口转发和容器生命周期管理
---
