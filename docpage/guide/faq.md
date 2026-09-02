# 常见问题

## 安装与启动

### 服务启动失败，日志提示 `store init` 错误

`deskcli` 启动时会在 `data_dir`（默认 `/opt/deskcli/data`）创建 SQLite 数据库。若目录不存在或权限不足会启动失败：

```bash
sudo mkdir -p /opt/deskcli/data
sudo chown -R root:root /opt/deskcli
sudo systemctl restart deskcli
```

查看完整日志：

```bash
journalctl -u deskcli -n 50 --no-pager
```

### 浏览器打开是 404 或空白页

Web 界面依赖静态资源目录。确认 `config.yaml` 中的 `static_dir` 指向实际存在的构建产物：

```yaml
static_dir: /opt/deskcli/web/dist
```

该目录下应有 `index.html` 与 `assets/`。若为空，说明安装时未部署前端，可手动构建：

```bash
cd web && npm install && npm run build
sudo cp -r ../deskcli/web/dist/* /opt/deskcli/web/dist/
```

若 `static_dir` 为空或路径不存在，服务只提供 REST API，不注册任何静态路由。

### 只能从本机访问，外部访问不通

默认 `listen_addr` 为 `127.0.0.1`。要允许外部访问需改为 `0.0.0.0` 并重启：

```yaml
listen_addr: 0.0.0.0
```

同时确认防火墙已放行端口：

```bash
sudo ufw allow 8080/tcp          # Debian/Ubuntu
sudo firewall-cmd --add-port=8080/tcp --permanent && sudo firewall-cmd --reload   # Fedora
```

## 账号与认证

### 忘记管理员密码

停止服务后直接操作数据库，删除用户记录，重启后系统会重新创建默认的 `admin` / `admin123`：

```bash
sudo systemctl stop deskcli
sudo sqlite3 /opt/deskcli/data/deskapi.db "DELETE FROM users;"
sudo systemctl start deskcli
```

### 注册接口返回 403

`POST /api/v1/auth/register` 只在系统中**尚无任何用户**时开放，用于创建第一个管理员。由于服务首次启动已自动创建 `admin`，该接口通常处于关闭状态。后续新增用户请以管理员身份调用 `POST /api/v1/users`。

### 旧的 X-Token 还能用吗

可以。认证中间件同时接受 JWT（`Authorization: Bearer <token>`）与旧版 `X-Token` 请求头，已有脚本无需改动。新集成建议使用 JWT，因为它携带用户身份与角色信息。

## 容器与端口

### 部署时提示 `no available ports`

端口池已耗尽。默认范围 `50000-60000` 可容纳 2500 个容器（每个占 4 个端口），出现该错误通常是范围被改窄了。调整配置：

```yaml
port_range: 50000-60000
```

已被数据库记录占用的端口不会重复分配，但**手动**用 `docker run` 起的容器所占端口不在记录内，可能造成冲突。

### 端口冲突导致容器启动失败

检查目标端口是否已被占用：

```bash
sudo ss -tlnp | grep <端口号>
```

若被无关进程占用，删除该容器后重新部署，系统会分配新端口。

### 容器创建了但连不上桌面

先确认容器确实在运行，再检查端口映射是否生效：

```bash
docker ps -a --filter name=<容器名>
docker port <容器名>
```

带桌面环境的镜像首次启动需要初始化 X11 或 Wayland 会话，等待十几秒再连接。Server 类模板不含桌面，只能 SSH。

## 镜像

### 镜像拉取缓慢

配置国内镜像加速器，编辑 `/etc/docker/daemon.json`：

```json
{
  "registry-mirrors": ["https://docker.mirrors.ustc.edu.cn"]
}
```

重启 Docker 生效：

```bash
sudo systemctl restart docker
```

也可以先用 `docker pull` 手动拉取，再在平台上部署，此时镜像已在本地。

### 删除镜像报错 image is being used

有容器仍在引用该镜像。先删除相关容器，再删镜像。

## 引擎差异

不同引擎支持的功能有差异：

| 功能 | Docker / Podman | LXD | 原生 LXC |
|------|:---------------:|:---:|:--------:|
| 容器增删改查 | 支持 | 支持 | 支持 |
| 镜像列表 / 删除 | 支持 | 支持 | 不支持 |
| 镜像搜索 | 支持 | 不支持 | 不支持 |
| SSE 流式拉取 | 支持 | 不支持 | 不支持 |
| Web 终端 | 支持 | 不支持 | 不支持 |

不支持的操作会返回明确的错误消息（如 `not supported for native LXC`），而非静默失败。LXC/LXD 下请改用 `deskcli exec` 进入容器。

## Windows 能用吗

管理平台本身依赖 Linux 容器运行时与 systemd，不提供原生 Windows 部署。Windows 用户请在 WSL2 中安装，或连接远程 Linux 主机。浏览器访问 Web 界面不受操作系统限制。
