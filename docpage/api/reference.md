# API 参考

`deskcli serve` 启动后提供 REST API，供脚本或外部系统调用。

- **Base URL**：`http://<host>:<port>/api/v1`

## 认证方式

接口支持两种认证方式，中间件会依次尝试：

| 方式 | 请求头 | 说明 |
|------|--------|------|
| JWT | `Authorization: Bearer <token>` | 推荐。通过登录接口获取，有效期 24 小时，携带用户身份与角色 |
| 静态 Token | `X-Token: <token>` | 兼容旧版脚本，取 `config.yaml` 中的 `token` 值，视为 admin 身份 |

以下三个接口无需认证：`POST /auth/login`、`POST /auth/register`、`GET /status`。

标注「仅 admin」的接口要求 JWT 中的角色为 `admin`，否则返回 403。

## 响应格式

除 SSE 流式接口外，所有响应为统一信封结构：

```json
{ "success": true, "data": { } }
```

失败时：

```json
{ "success": false, "message": "invalid credentials" }
```

---

## 认证

### 登录

**POST** `/auth/login`

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

响应：

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": { "id": 1, "username": "admin", "role": "admin" }
  }
}
```

### 注册首个管理员

**POST** `/auth/register`

仅在系统中尚无任何用户时可用，否则返回 403。用户名 3–32 字符，密码至少 8 位。创建的账号角色为 `admin`。

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "operator", "password": "s3cure-pass"}'
```

### 获取当前用户

**GET** `/auth/me`

```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/auth/me
```

### 刷新 Token

**POST** `/auth/refresh`

以当前身份签发一个新的 24 小时 Token。

### 修改密码

**POST** `/auth/change-password`

```bash
curl -X POST http://localhost:8080/api/v1/auth/change-password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"old_password": "admin123", "new_password": "new-s3cure-pass"}'
```

新密码至少 8 位。旧密码不匹配返回 401。

---

## 用户管理

### 用户列表

**GET** `/users` —— 仅 admin

### 创建用户

**POST** `/users` —— 仅 admin

| 字段 | 类型 | 说明 |
|------|------|------|
| `username` | string | 3–32 字符，必填 |
| `password` | string | 至少 8 位，必填 |
| `role` | string | `admin` 或 `user`，非法值回退为 `user` |

用户名已存在返回 409。

---

## 系统状态

### 公开状态

**GET** `/status`

无需认证，返回引擎类型、端口与容器计数，可用于健康检查。

---

## 容器列表

**GET** `/containers`

列出所有受管容器。

```bash
curl -H "X-Token: your_token" http://localhost:8080/api/v1/containers
```

响应示例：

```json
[
  {
    "name": "my-desktop",
    "image": "pikachuim/debian:12-xfce4l",
    "status": "running",
    "ports": ["4000:4000", "5900:5900", "2222:22"]
  }
]
```

---

## 创建容器

**POST** `/containers`

```bash
curl -X POST http://localhost:8080/api/v1/containers \
  -H "X-Token: your_token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-desktop",
    "image": "pikachuim/debian:12-xfce4l",
    "ports": ["4000:4000", "5900:5900", "2222:22"]
  }'
```

请求体字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 容器名称（必填） |
| `image` | string | 镜像名称（必填） |
| `ports` | string[] | 端口映射列表，格式 `"宿主:容器"` |
| `volumes` | string[] | 挂载卷列表 |
| `env` | object | 环境变量 |

---

## 获取容器详情

**GET** `/containers/:name`

```bash
curl -H "X-Token: your_token" http://localhost:8080/api/v1/containers/my-desktop
```

响应示例：

```json
{
  "name": "my-desktop",
  "image": "pikachuim/debian:12-xfce4l",
  "status": "running",
  "created": "2024-01-01T00:00:00Z",
  "ports": ["4000:4000", "5900:5900", "2222:22"]
}
```

---

## 删除容器

**DELETE** `/containers/:name`

```bash
curl -X DELETE -H "X-Token: your_token" \
  http://localhost:8080/api/v1/containers/my-desktop
```

---

## 启动容器

**POST** `/containers/:name/start`

```bash
curl -X POST -H "X-Token: your_token" \
  http://localhost:8080/api/v1/containers/my-desktop/start
```

---

## 停止容器

**POST** `/containers/:name/stop`

```bash
curl -X POST -H "X-Token: your_token" \
  http://localhost:8080/api/v1/containers/my-desktop/stop
```

---

## 重启容器

**POST** `/containers/:name/restart`

```bash
curl -X POST -H "X-Token: your_token" \
  http://localhost:8080/api/v1/containers/my-desktop/restart
```

---

## 执行命令

**POST** `/containers/:name/exec`

```bash
curl -X POST http://localhost:8080/api/v1/containers/my-desktop/exec \
  -H "X-Token: your_token" \
  -H "Content-Type: application/json" \
  -d '{"cmd": "ls -la /home/user"}'
```

请求体：

| 字段 | 类型 | 说明 |
|------|------|------|
| `cmd` | string | 要执行的命令 |
| `detach` | bool | 是否后台执行（默认 false） |

响应示例：

```json
{
  "output": "total 32\ndrwxr-xr-x ...",
  "exit_code": 0
}
```

---

## 修改密码

**POST** `/containers/:name/passwd`

```bash
curl -X POST http://localhost:8080/api/v1/containers/my-desktop/passwd \
  -H "X-Token: your_token" \
  -H "Content-Type: application/json" \
  -d '{"password": "newpassword"}'
```

---

## 镜像管理

### 镜像列表

**GET** `/images`

```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/images
```

响应：

```json
{
  "success": true,
  "data": [
    {
      "id": "sha256:abc123",
      "tags": ["pikachuim/debian:12-xfce4l"],
      "size": "1.2GB",
      "created": "2024-01-01 00:00:00 +0800 CST"
    }
  ]
}
```

原生 LXC 不支持该接口。

### 拉取镜像

**POST** `/images/pull`

以 Server-Sent Events 流式返回拉取进度，`Content-Type: text/event-stream`。

```bash
curl -N -X POST http://localhost:8080/api/v1/images/pull \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"image": "pikachuim/debian:12-xfce4l"}'
```

每行输出对应引擎的一行原始进度；结束时发送 `data: {"status":"done"}`，失败时发送 `data: {"error":"..."}`。仅支持 Docker 与 Podman。

### 删除镜像

**DELETE** `/images/:id`

### 搜索镜像

**GET** `/images/search?q=<关键词>`

检索 Docker Hub，返回名称、星标数与描述。仅支持 Docker 与 Podman。

---

## 模板管理

### 模板列表

**GET** `/templates`

支持三个可选查询参数进行过滤：`os_type`、`de_name`、`category`。

```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/templates?os_type=debian&de_name=xfce4l"
```

响应：

```json
{
  "success": true,
  "data": [
    {
      "id": 4,
      "os_type": "debian",
      "os_version": "12",
      "de_name": "xfce4l",
      "display_name": "Debian 12 · Xfce 4",
      "description": "Debian 12 with fast, lightweight Xfce 4 desktop.",
      "image_full": "pikachuim/debian:12-xfce4l",
      "icon": "zap",
      "category": "desktop",
      "is_active": true
    }
  ]
}
```

### 可选筛选值

**GET** `/templates/os-types` —— 返回所有发行版名称数组

**GET** `/templates/de-names` —— 返回所有桌面环境名称数组

### 一键部署

**POST** `/templates/:id/deploy`

```bash
curl -X POST http://localhost:8080/api/v1/templates/4/deploy \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "my-desktop"}'
```

请求体字段全部可选：

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 容器名称，留空则自动生成 `{os}-{de}` |
| `engine` | string | 覆盖配置中的默认引擎 |
| `ports` | string[] | 端口映射，留空则从 `port_range` 自动分配 4 个 |
| `softwares` | string[] | 部署后后台安装的软件包 |

自动分配时依次映射到容器的 22（SSH）、3389（RDP）、4000（NX）、5900（VNC）。

响应：

```json
{
  "success": true,
  "data": {
    "name": "my-desktop",
    "image": "pikachuim/debian:12-xfce4l",
    "ports": ["50000:22", "50001:3389", "50002:4000", "50003:5900"],
    "engine": "docker"
  }
}
```

---

## 监控

### 引擎统计

**GET** `/monitor/stats`

```json
{
  "success": true,
  "data": {
    "containers_running": 3,
    "containers_total": 5,
    "images_total": 12,
    "engine_version": "24.0.7"
  }
}
```

### 引擎可用性

**GET** `/monitor/engines`

探测 Docker、Podman、LXC、LXD 四种引擎，返回各自的 `name`、`available` 与 `version`。

---

## 配置

### 读取配置

**GET** `/config` —— 仅 admin

返回脱敏后的配置，不含 `token` 与 `jwt_secret`。

### 更新配置

**PUT** `/config` —— 仅 admin

```bash
curl -X PUT http://localhost:8080/api/v1/config \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"port_range": "40000-45000"}'
```

可更新 `engine`、`port_range`、`allow_exec`、`static_dir`、`jwt_secret`，省略的字段保持原值。改动会写回 `config.yaml`，部分项需重启服务生效。

::: warning 更换 jwt_secret
修改 `jwt_secret` 会使所有已签发的 Token 立即失效，全部用户需重新登录。
:::

---

## WebSocket 终端

**GET** `/ws/terminal/:name`

WebSocket 握手无法自定义请求头，因此认证 Token 通过查询参数传递：

```
ws://<host>:<port>/api/v1/ws/terminal/my-desktop?token=<JWT>
```

旧版静态 Token 使用 `?x_token=<token>`。

消息均为 JSON 文本帧，其中 `data` 字段经 base64 编码：

客户端发往服务端：

```json
{ "type": "stdin", "data": "bHMgLWxhCg==" }
{ "type": "resize", "cols": 120, "rows": 40 }
```

服务端发往客户端：

```json
{ "type": "stdout", "data": "dG90YWwgMzIK" }
{ "type": "exit", "code": 0 }
{ "type": "error", "message": "exec failed: ..." }
```

::: tip 已知限制
仅支持 Docker 与 Podman 引擎。由于未使用 PTY，`resize` 消息目前为空操作，依赖终端尺寸的全屏程序显示可能异常。
:::

---

## 错误码

| HTTP 状态码 | 说明 |
|------------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 认证失败（Token 无效或凭据错误） |
| 403 | 权限不足（需要 admin 角色，或接口已关闭） |
| 404 | 资源不存在 |
| 409 | 资源冲突（如用户名已被占用） |
| 500 | 服务器内部错误 |

错误响应格式：

```json
{
  "success": false,
  "message": "container not found: my-desktop"
}
```
