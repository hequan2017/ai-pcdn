# ai-pcdn

> Modified for ai-pcdn on 2026-07-21: 增加产品说明与部署文档。
> 本项目基于 Gin-Vue-Admin 开发，原项目版权、授权声明及 [LICENSE](./LICENSE) 保持有效。

ai-pcdn 是基于 Go、Gin、Vue 3 和 Vite 构建的 PCDN 管理平台。仓库包含：

- `server/`：Go 后端，默认监听 `8888`
- `web/`：Vue 3 前端，开发及容器端口为 `8080`
- `deploy/`：Docker Compose、容器构建和运行配置
- `deploy.sh`：基于 Docker Compose 的一键部署脚本

## 一键部署

### 环境要求

- Docker 24 或更高版本
- Docker Compose v2（使用 `docker compose` 命令）
- `curl`
- Linux、macOS、WSL 或 Git Bash

部署前请确认 Docker 服务已经启动，并确保目标机器的 `8080`、`8888` 端口未被占用。

### 执行部署

推荐在首次部署时显式设置管理员密码：

```bash
AI_PCDN_ADMIN_PASSWORD='请替换为至少6位的强密码' bash ./deploy.sh
```

不设置密码也可以直接运行：

```bash
bash ./deploy.sh
```

此时管理员初始密码为 `123456`，仅适合本地开发或临时验证，首次登录后必须立即修改。

脚本会依次完成以下操作：

1. 检查 Docker、Docker Compose 和 `curl`
2. 创建本地运行配置 `deploy/docker-compose/runtime/config.yaml`
3. 构建并启动前后端容器
4. 等待后端健康检查通过
5. 首次部署时初始化 SQLite 数据库和基础数据
6. 验证前端页面及 `/api` 代理

重复执行脚本会复用已有配置和数据库，不会重新初始化或删除数据。

部署成功后访问：

- 前端：<http://127.0.0.1:8080>
- 后端：<http://127.0.0.1:8888>
- Swagger：<http://127.0.0.1:8888/swagger/index.html>
- 默认管理员账号：`admin`

### 自定义参数

| 环境变量 | 默认值 | 说明 |
| --- | --- | --- |
| `AI_PCDN_ADMIN_PASSWORD` | `123456` | 首次初始化时的管理员密码，至少 6 位 |
| `AI_PCDN_WEB_PORT` | `8080` | 前端宿主机端口 |
| `AI_PCDN_SERVER_PORT` | `8888` | 后端宿主机端口 |
| `AI_PCDN_DB_NAME` | `ai-pcdn` | 首次初始化时的 SQLite 数据库名 |

示例：

```bash
AI_PCDN_WEB_PORT=18080 \
AI_PCDN_SERVER_PORT=18888 \
AI_PCDN_ADMIN_PASSWORD='ChangeMe_2026' \
bash ./deploy.sh
```

## 容器运维

查看运行状态：

```bash
docker compose -f deploy/docker-compose/docker-compose.yaml ps
```

查看日志：

```bash
docker compose -f deploy/docker-compose/docker-compose.yaml logs -f --tail=200
```

重新构建并升级：

```bash
bash ./deploy.sh
```

停止服务但保留数据库和上传文件：

```bash
docker compose -f deploy/docker-compose/docker-compose.yaml down
```

数据库存储在 Docker 卷 `ai-pcdn_server-data`，上传文件存储在 `ai-pcdn_uploads`。不要使用带 `-v` 的 `docker compose down`，否则会删除持久化数据。

## 本地开发

后端：

```bash
cd server
go run .
```

前端：

```bash
cd web
npm install
npm run dev
```

首次本地启动时，可通过初始化页面选择 SQLite；建议数据库名使用 `gva`、数据库路径使用 `.`。默认开发地址与容器部署地址一致。

## 验证

```bash
cd server && go build ./...
cd ../web && npm run build
docker compose -f ../deploy/docker-compose/docker-compose.yaml config
```

## 授权与版权

ai-pcdn 对产品展示名称进行了定制，但 Gin-Vue-Admin 的版权、作者归属、商用授权提示及许可证义务仍须保留。生产环境或组织部署前，请阅读仓库根目录的 [LICENSE](./LICENSE) 并确认已具备相应授权。
