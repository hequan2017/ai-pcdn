#!/usr/bin/env sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
COMPOSE_FILE="$SCRIPT_DIR/deploy/docker-compose/docker-compose.yaml"
RUNTIME_DIR="$SCRIPT_DIR/deploy/docker-compose/runtime"
RUNTIME_CONFIG="$RUNTIME_DIR/config.yaml"
CONFIG_TEMPLATE="$SCRIPT_DIR/server/config.docker.yaml"

WEB_PORT=${AI_PCDN_WEB_PORT:-8080}
SERVER_PORT=${AI_PCDN_SERVER_PORT:-8888}
ADMIN_PASSWORD=${AI_PCDN_ADMIN_PASSWORD:-123456}
DB_NAME=${AI_PCDN_DB_NAME:-ai-pcdn}

log() {
  printf '[ai-pcdn] %s\n' "$*"
}

fail() {
  printf '[ai-pcdn] ERROR: %s\n' "$*" >&2
  exit 1
}

compose() {
  docker compose -f "$COMPOSE_FILE" "$@"
}

wait_for_url() {
  wait_url=$1
  wait_name=$2
  wait_attempt=0

  while [ "$wait_attempt" -lt 60 ]; do
    if curl --fail --silent --show-error "$wait_url" >/dev/null 2>&1; then
      return 0
    fi
    wait_attempt=$((wait_attempt + 1))
    sleep 2
  done

  fail "$wait_name 在 120 秒内未就绪，请执行 docker compose -f \"$COMPOSE_FILE\" logs 查看日志"
}

json_escape() {
  printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'
}

command -v docker >/dev/null 2>&1 || fail "未安装 Docker"
command -v curl >/dev/null 2>&1 || fail "未安装 curl"
docker compose version >/dev/null 2>&1 || fail "未安装 Docker Compose v2"
docker info >/dev/null 2>&1 || fail "Docker 服务未运行"

if [ "${#ADMIN_PASSWORD}" -lt 6 ]; then
  fail "AI_PCDN_ADMIN_PASSWORD 长度不能少于 6 位"
fi

mkdir -p "$RUNTIME_DIR"
if [ ! -f "$RUNTIME_CONFIG" ]; then
  cp "$CONFIG_TEMPLATE" "$RUNTIME_CONFIG"
  log "已创建本地运行配置：$RUNTIME_CONFIG"
fi

log "构建并启动容器"
compose up -d --build

BACKEND_URL="http://127.0.0.1:$SERVER_PORT"
FRONTEND_URL="http://127.0.0.1:$WEB_PORT"

log "等待后端服务就绪"
wait_for_url "$BACKEND_URL/health" "后端服务"

check_response=$(curl --fail --silent --show-error \
  --request POST \
  --header 'Content-Type: application/json' \
  "$BACKEND_URL/init/checkdb")

initialized_now=false
case "$check_response" in
  *'"needInit":true'*)
    escaped_password=$(json_escape "$ADMIN_PASSWORD")
    escaped_db_name=$(json_escape "$DB_NAME")
    init_payload=$(printf \
      '{"adminPassword":"%s","dbType":"sqlite","dbName":"%s","dbPath":"/data"}' \
      "$escaped_password" \
      "$escaped_db_name")

    log "初始化本地 SQLite 数据库"
    init_response=$(printf '%s' "$init_payload" | curl --fail --silent --show-error \
      --request POST \
      --header 'Content-Type: application/json' \
      --data-binary @- \
      "$BACKEND_URL/init/initdb")

    case "$init_response" in
      *'"code":0'*) initialized_now=true ;;
      *)
        printf '%s\n' "$init_response" >&2
        fail "数据库初始化失败"
        ;;
    esac
    ;;
  *'"needInit":false'*)
    log "数据库已经初始化，跳过初始化步骤"
    ;;
  *)
    printf '%s\n' "$check_response" >&2
    fail "无法识别数据库初始化状态"
    ;;
esac

log "等待前端服务就绪"
wait_for_url "$FRONTEND_URL/" "前端服务"
wait_for_url "$FRONTEND_URL/api/health" "前端 API 代理"

log "部署完成"
printf '前端地址: %s\n' "$FRONTEND_URL"
printf '后端地址: %s\n' "$BACKEND_URL"
printf 'Swagger: %s/swagger/index.html\n' "$BACKEND_URL"
if [ "$initialized_now" = true ]; then
  printf '管理员账号: admin\n'
  printf '管理员初始密码: %s（首次登录后请立即修改）\n' "$ADMIN_PASSWORD"
fi
