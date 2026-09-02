# Web 管理面板

`deskcli serve` 启动后，除 REST API 外还会在同一端口提供 React 管理界面。本文按页面介绍各项功能。

## 登录与账号

平台使用 JWT 认证，Token 有效期 24 小时，存放在浏览器 `localStorage`，由前端在每次请求时通过 `Authorization: Bearer <token>` 头携带。Token 过期或失效时会自动跳回登录页。

角色分为两种：

| 角色 | 权限 |
|------|------|
| `admin` | 全部功能，包括系统配置与用户管理 |
| `user` | 容器、镜像、模板的日常操作 |

首次启动时若用户表为空，系统会自动创建 `admin` / `admin123` 管理员账号，请登录后立即修改。

## 仪表板

进入后的默认页面，提供四张统计卡片与最近容器列表：

- **运行中容器** / **总容器数** —— 来自 `GET /api/v1/monitor/stats`
- **镜像总数** —— 本地已拉取的镜像数量
- **引擎版本** —— 当前生效的容器引擎及其版本

下方「最近容器」列出最新的 5 个容器，绿色脉冲圆点表示运行中，红色表示已停止。

## 容器管理

以卡片网格展示全部容器，支持按名称或镜像搜索，并可按「全部 / 运行中 / 已停止」筛选。

每张卡片提供的操作：

| 操作 | 可用状态 | 说明 |
|------|----------|------|
| 启动 | 已停止 | 调用 `POST /containers/:name/start` |
| 停止 | 运行中 | 调用 `POST /containers/:name/stop` |
| 重启 | 运行中 | 调用 `POST /containers/:name/restart` |
| 终端 | 运行中 | 跳转终端页并自动选中该容器 |
| 删除 | 任意 | 二次确认后强制删除 |

::: warning 删除不可恢复
删除会执行 `docker rm -f`，容器内未持久化的数据将全部丢失。
:::

## 模板市场

预置 30 余个模板，覆盖 Debian、Ubuntu、Alpine、Fedora、Arch 五种发行版与 Server、X11GUI、GNOME 3、Xfce 4、KDE Plasma、Deepin、Lingmo、Hyprland、Niri 等桌面环境。镜像命名遵循 `pikachuim/{os}:{version}-{de}` 规则。

顶部标签栏可按桌面环境筛选。点击卡片的 **一键部署** 后填写容器名称即可，端口由后端自动分配——从配置的 `port_range` 中挑选 4 个未被占用的端口，依次映射到容器的 22、3389、4000、5900。

若需自定义端口或在部署时预装软件包，可直接调用 `POST /api/v1/templates/:id/deploy`，传入 `ports` 与 `softwares` 字段。

## Web 终端

基于 xterm.js 与 WebSocket，在浏览器内提供交互式 Shell，无需 SSH 客户端。

下拉框仅列出**运行中**的容器。连接建立后右上角状态点变绿。终端通过 `/api/v1/ws/terminal/:name` 建立连接，认证 Token 以查询参数传递（WebSocket 握手无法自定义请求头）。

::: tip 引擎限制
Web 终端目前仅支持 Docker 与 Podman。LXC/LXD 引擎下会返回明确的错误提示，请改用 `deskcli exec` 命令。
:::

由于未引入 PTY 库，终端使用管道转发标准输入输出。绝大多数命令行操作正常，但 `resize` 消息为空操作，依赖终端尺寸的全屏程序（如 `htop`、`vim`）显示可能不理想。

## 系统设置

分为三个标签页：

- **通用配置** —— 只读展示当前引擎、API 端口、端口范围，对应配置文件 `/opt/deskcli/conf/config.yaml`
- **引擎管理** —— 探测 Docker、Podman、LXC、LXD 四种引擎的可用性与版本
- **关于** —— 版本信息与支持的系统/桌面环境列表

修改配置需管理员权限，可通过 `PUT /api/v1/config` 完成，改动会写回 YAML 配置文件。部分项（如监听端口）需重启服务后生效：

```bash
systemctl restart deskcli
```
