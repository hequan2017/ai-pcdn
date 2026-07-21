# Docker Compose 一键部署与 SQLite 初始化

- 状态：已完成
- 日期：2026-07-21

## 需求

更新根 README 的部署说明，并增加可重复执行的一键部署脚本。

## 已实现

- 根目录新增 `deploy.sh`。
- 脚本检查 Docker、Docker Compose v2 和 curl，构建并启动前后端容器。
- 首次部署自动调用后端初始化接口，创建本地 SQLite 数据库与基础数据。
- 重复部署通过 `/init/checkdb` 判断状态，不重复初始化或覆盖数据库。
- Compose 使用 `ai-pcdn_server-data` 和 `ai-pcdn_uploads` 持久化数据库及上传文件。
- Nginx 通过 Compose 服务名 `server` 转发 `/api/` 请求，不再依赖固定容器 IP。
- README 记录一键部署、端口覆盖、日志、升级、停止和数据保护方式。

## 验证

- Git Bash `bash -n deploy.sh` 通过。
- `docker compose -f deploy/docker-compose/docker-compose.yaml config --quiet` 通过。
- 当前 Docker Desktop 未运行，因此未执行镜像拉取、容器构建和真实部署。

