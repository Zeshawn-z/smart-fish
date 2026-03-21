#!/bin/bash
#
# 智钓蓝海 - 一键启动脚本
# 使用方式: chmod +x start.sh && ./start.sh
#

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$PROJECT_DIR/back_end"
FRONTEND_DIR="$PROJECT_DIR/front_end"

# 日志函数
info()    { echo -e "${BLUE}[INFO]${NC}  $1"; }
success() { echo -e "${GREEN}[OK]${NC}    $1"; }
warn()    { echo -e "${YELLOW}[WARN]${NC}  $1"; }
error()   { echo -e "${RED}[ERROR]${NC} $1"; }

echo -e "${CYAN}"
echo "  ╔══════════════════════════════════════╗"
echo "  ║        智钓蓝海 - 启动脚本            ║"
echo "  ║   Smart Fishing Information Platform  ║"
echo "  ╚══════════════════════════════════════╝"
echo -e "${NC}"

# ===== 1. 环境检查 =====
info "正在检查运行环境..."

# 检查 Go
if command -v go &> /dev/null; then
  GO_VERSION=$(go version | awk '{print $3}')
  success "Go 已安装: $GO_VERSION"
else
  error "未检测到 Go 环境，请先安装 Go (>= 1.21)"
  echo "  安装方法: brew install go  或  https://go.dev/dl/"
  exit 1
fi

# 检查 Node.js
if command -v node &> /dev/null; then
  NODE_VERSION=$(node --version)
  success "Node.js 已安装: $NODE_VERSION"
else
  error "未检测到 Node.js 环境，请先安装 Node.js (>= 18)"
  echo "  安装方法: brew install node  或  https://nodejs.org/"
  exit 1
fi

# 检查 npm
if command -v npm &> /dev/null; then
  NPM_VERSION=$(npm --version)
  success "npm 已安装: $NPM_VERSION"
else
  error "未检测到 npm"
  exit 1
fi

# 检查 MySQL
MYSQL_AVAILABLE=false
if command -v mysql &> /dev/null; then
  if mysql -u root -e "SELECT 1" &>/dev/null 2>&1; then
    success "MySQL 已安装且可连接（root 无密码）"
    MYSQL_AVAILABLE=true
  elif mysql -u root -p'' -e "SELECT 1" &>/dev/null 2>&1; then
    success "MySQL 已安装且可连接"
    MYSQL_AVAILABLE=true
  else
    warn "MySQL 已安装但无法连接（可能需要密码或未启动）"
    warn "后端将使用默认配置尝试连接，如失败请检查 back_end/config/config.go"
  fi
else
  warn "未检测到 MySQL 客户端"
  warn "后端需要 MySQL 数据库，请确保 MySQL 服务已启动"
  warn "如仅需查看前端效果，前端有 Mock 数据可独立运行"
fi

echo ""

# ===== 2. 选择启动模式 =====
echo -e "${CYAN}请选择启动模式:${NC}"
echo "  1) 前后端一起启动（推荐）"
echo "  2) 仅启动前端（使用 Mock 数据）"
echo "  3) 仅启动后端"
echo "  4) 前后端一起启动 + 重置种子数据"
echo ""
read -p "请输入选项 [1/2/3/4，默认 1]: " MODE
MODE=${MODE:-1}

# 清理函数
cleanup() {
  echo ""
  info "正在停止服务..."
  if [ -n "$BACKEND_PID" ]; then
    kill $BACKEND_PID 2>/dev/null && success "后端已停止 (PID: $BACKEND_PID)"
  fi
  if [ -n "$FRONTEND_PID" ]; then
    kill $FRONTEND_PID 2>/dev/null && success "前端已停止 (PID: $FRONTEND_PID)"
  fi
  exit 0
}

trap cleanup SIGINT SIGTERM

BACKEND_PID=""
FRONTEND_PID=""

# ===== 3. 启动后端 =====
start_backend() {
  local SEED_FLAG="$1"

  info "正在构建后端..."
  cd "$BACKEND_DIR"

  # 下载依赖
  if [ ! -d "vendor" ] && [ ! -d "$HOME/go/pkg/mod/github.com" ]; then
    info "首次运行，正在下载 Go 依赖..."
    go mod download
  fi

  # 编译
  go build -o smart-fish-server . 2>&1
  if [ $? -ne 0 ]; then
    error "后端编译失败！"
    exit 1
  fi
  success "后端编译成功"

  # 启动
  if [ "$SEED_FLAG" = "--seed" ]; then
    info "正在启动后端服务 (端口 8080) [含种子数据]..."
    ./smart-fish-server --seed &
  else
    info "正在启动后端服务 (端口 8080)..."
    ./smart-fish-server &
  fi
  BACKEND_PID=$!
  sleep 2

  if kill -0 $BACKEND_PID 2>/dev/null; then
    success "后端已启动 (PID: $BACKEND_PID)"
    echo -e "  ${GREEN}→ http://localhost:8080${NC}"
  else
    error "后端启动失败，请检查日志"
    warn "可能原因: MySQL 未启动 / 端口被占用 / 配置错误"
  fi
}

# ===== 4. 启动前端 =====
start_frontend() {
  info "正在准备前端..."
  cd "$FRONTEND_DIR"

  # 安装依赖
  if [ ! -d "node_modules" ]; then
    info "首次运行，正在安装前端依赖..."
    npm install
  fi

  # 启动开发服务器
  info "正在启动前端开发服务器 (端口 5173)..."
  npm run dev &
  FRONTEND_PID=$!
  sleep 3

  if kill -0 $FRONTEND_PID 2>/dev/null; then
    success "前端已启动 (PID: $FRONTEND_PID)"
    echo -e "  ${GREEN}→ http://localhost:5173${NC}"
  else
    error "前端启动失败"
  fi
}

# ===== 5. 执行 =====
echo ""
case $MODE in
  1)
    info "模式: 前后端同时启动"
    echo ""
    start_backend
    echo ""
    start_frontend
    ;;
  2)
    info "模式: 仅启动前端（Mock 数据模式）"
    warn "前端将使用内置的模拟数据运行，无需后端和数据库"
    echo ""
    start_frontend
    ;;
  3)
    info "模式: 仅启动后端"
    echo ""
    start_backend
    ;;
  4)
    info "模式: 前后端同时启动 + 重置种子数据"
    echo ""
    start_backend --seed
    echo ""
    start_frontend
    ;;
  *)
    error "无效选项: $MODE"
    exit 1
    ;;
esac

echo ""
echo -e "${CYAN}════════════════════════════════════════${NC}"
echo -e "${GREEN}  服务已启动！按 Ctrl+C 停止所有服务${NC}"
echo -e "${CYAN}════════════════════════════════════════${NC}"
echo ""

if [ "$MODE" = "1" ] || [ "$MODE" = "2" ] || [ "$MODE" = "4" ]; then
  echo -e "  前端地址: ${GREEN}http://localhost:5173${NC}"
fi
if [ "$MODE" = "1" ] || [ "$MODE" = "3" ] || [ "$MODE" = "4" ]; then
  echo -e "  后端地址: ${GREEN}http://localhost:8080${NC}"
fi
echo ""

# 等待退出信号
wait
