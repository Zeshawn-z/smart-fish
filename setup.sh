#!/usr/bin/env bash
# ============================================================================
# 智钓蓝海 (Smart Fish) 一键初始化脚本
# 功能：检查/安装 MySQL → 用户输入密码 → 生成 .env → 连接并创建数据库
# 用法：chmod +x setup.sh && ./setup.sh
# ============================================================================

set -e

# ===== 颜色 =====
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { echo -e "${CYAN}[INFO]${NC}  $1"; }
ok()    { echo -e "${GREEN}[OK]${NC}    $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $1"; }
fail()  { echo -e "${RED}[FAIL]${NC}  $1"; exit 1; }

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ENV_FILE="$SCRIPT_DIR/.env"

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}   🐟 智钓蓝海 — 环境初始化脚本${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# ===== 1. 检测操作系统 =====
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command -v apt-get &>/dev/null; then
            OS="debian"
        elif command -v yum &>/dev/null; then
            OS="rhel"
        elif command -v dnf &>/dev/null; then
            OS="rhel"
        elif command -v pacman &>/dev/null; then
            OS="arch"
        else
            OS="linux_unknown"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
    else
        OS="unknown"
    fi
    info "检测到操作系统: $OS ($OSTYPE)"
}

# ===== 2. 检查 MySQL 是否安装 =====
check_mysql() {
    if command -v mysql &>/dev/null; then
        MYSQL_VERSION=$(mysql --version 2>/dev/null || true)
        ok "MySQL 客户端已安装: $MYSQL_VERSION"
        return 0
    else
        warn "未检测到 MySQL 客户端"
        return 1
    fi
}

# ===== 3. 检查 MySQL 服务是否运行 =====
check_mysql_service() {
    # 尝试多种方式检测
    if command -v systemctl &>/dev/null; then
        if systemctl is-active --quiet mysql 2>/dev/null || systemctl is-active --quiet mysqld 2>/dev/null || systemctl is-active --quiet mariadb 2>/dev/null; then
            ok "MySQL 服务正在运行"
            return 0
        fi
    fi
    if command -v brew &>/dev/null && brew services list 2>/dev/null | grep -q 'mysql.*started'; then
        ok "MySQL 服务正在运行 (Homebrew)"
        return 0
    fi
    # 尝试直接连接检测
    if mysqladmin ping -h 127.0.0.1 --silent 2>/dev/null; then
        ok "MySQL 服务正在运行"
        return 0
    fi
    warn "MySQL 服务未运行"
    return 1
}

# ===== 4. 安装 MySQL =====
install_mysql() {
    echo ""
    echo -e "${YELLOW}是否要自动安装 MySQL？(y/n)${NC}"
    read -r INSTALL_CHOICE
    if [[ "$INSTALL_CHOICE" != "y" && "$INSTALL_CHOICE" != "Y" ]]; then
        fail "请手动安装 MySQL 5.7+ 后重新运行此脚本"
    fi

    info "正在安装 MySQL ..."
    case "$OS" in
        debian)
            sudo apt-get update
            sudo apt-get install -y mysql-server mysql-client
            sudo systemctl start mysql
            sudo systemctl enable mysql
            ;;
        rhel)
            if command -v dnf &>/dev/null; then
                sudo dnf install -y mysql-server mysql
                sudo systemctl start mysqld
                sudo systemctl enable mysqld
            else
                sudo yum install -y mysql-server mysql
                sudo systemctl start mysqld
                sudo systemctl enable mysqld
            fi
            ;;
        arch)
            sudo pacman -Sy --noconfirm mariadb
            sudo mariadb-install-db --user=mysql --basedir=/usr --datadir=/var/lib/mysql
            sudo systemctl start mariadb
            sudo systemctl enable mariadb
            ;;
        macos)
            if ! command -v brew &>/dev/null; then
                fail "请先安装 Homebrew: https://brew.sh"
            fi
            brew install mysql
            brew services start mysql
            ;;
        *)
            fail "不支持自动安装 MySQL，请手动安装后重试"
            ;;
    esac

    # 验证安装
    sleep 2
    if check_mysql; then
        ok "MySQL 安装成功"
    else
        fail "MySQL 安装失败，请手动安装"
    fi
}

# ===== 5. 启动 MySQL 服务 =====
start_mysql_service() {
    echo ""
    echo -e "${YELLOW}是否要尝试启动 MySQL 服务？(y/n)${NC}"
    read -r START_CHOICE
    if [[ "$START_CHOICE" != "y" && "$START_CHOICE" != "Y" ]]; then
        fail "请手动启动 MySQL 服务后重新运行此脚本"
    fi

    info "正在启动 MySQL 服务 ..."
    case "$OS" in
        debian|rhel|arch)
            sudo systemctl start mysql 2>/dev/null || sudo systemctl start mysqld 2>/dev/null || sudo systemctl start mariadb 2>/dev/null
            ;;
        macos)
            brew services start mysql 2>/dev/null || true
            ;;
    esac

    sleep 2
    if check_mysql_service; then
        ok "MySQL 服务已启动"
    else
        fail "MySQL 服务启动失败，请手动启动"
    fi
}

# ===== 6. 收集数据库配置 =====
collect_config() {
    echo ""
    echo -e "${CYAN}--- 数据库配置 ---${NC}"
    echo ""

    # 数据库主机
    read -rp "数据库主机 [127.0.0.1]: " DB_HOST
    DB_HOST=${DB_HOST:-127.0.0.1}

    # 数据库端口
    read -rp "数据库端口 [3306]: " DB_PORT
    DB_PORT=${DB_PORT:-3306}

    # 数据库用户
    read -rp "数据库用户名 [root]: " DB_USER
    DB_USER=${DB_USER:-root}

    # 数据库密码（隐藏输入）
    while true; do
        echo -n "数据库密码: "
        read -rs DB_PASSWORD
        echo ""
        if [[ -z "$DB_PASSWORD" ]]; then
            warn "密码不能为空，请重新输入"
            continue
        fi
        echo -n "确认密码: "
        read -rs DB_PASSWORD_CONFIRM
        echo ""
        if [[ "$DB_PASSWORD" != "$DB_PASSWORD_CONFIRM" ]]; then
            warn "两次密码不一致，请重新输入"
            continue
        fi
        break
    done

    # 数据库名
    read -rp "数据库名称 [smart_fish]: " DB_NAME
    DB_NAME=${DB_NAME:-smart_fish}

    # 服务端口
    read -rp "后端服务端口 [8080]: " SERVER_PORT
    SERVER_PORT=${SERVER_PORT:-8080}

    echo ""
    info "配置预览:"
    echo "  数据库: ${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
    echo "  服务端口: ${SERVER_PORT}"
    echo ""
}

# ===== 7. 测试数据库连接 =====
test_connection() {
    info "测试数据库连接 ..."
    if mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" &>/dev/null; then
        ok "数据库连接成功"
        return 0
    else
        fail "无法连接到 MySQL（主机: ${DB_HOST}:${DB_PORT}，用户: ${DB_USER}），请检查密码和服务状态"
    fi
}

# ===== 8. 创建数据库 =====
create_database() {
    info "检查数据库 '${DB_NAME}' ..."
    DB_EXISTS=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -N -e \
        "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME='${DB_NAME}';" 2>/dev/null)

    if [[ -n "$DB_EXISTS" ]]; then
        ok "数据库 '${DB_NAME}' 已存在"
        echo -e "${YELLOW}是否要清空重建数据库？这会删除所有数据！(y/n)${NC}"
        read -r REBUILD_CHOICE
        if [[ "$REBUILD_CHOICE" == "y" || "$REBUILD_CHOICE" == "Y" ]]; then
            info "正在重建数据库 ..."
            mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -e \
                "DROP DATABASE \`${DB_NAME}\`; CREATE DATABASE \`${DB_NAME}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
            ok "数据库已重建"
        fi
    else
        info "正在创建数据库 '${DB_NAME}' ..."
        mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -e \
            "CREATE DATABASE \`${DB_NAME}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
        ok "数据库 '${DB_NAME}' 创建成功"
    fi
}

# ===== 9. 生成 .env 文件 =====
generate_env() {
    # 生成随机 JWT Secret
    JWT_SECRET=$(head -c 32 /dev/urandom | base64 | tr -d '/+=' | head -c 32)

    if [[ -f "$ENV_FILE" ]]; then
        warn ".env 文件已存在: $ENV_FILE"
        echo -e "${YELLOW}是否要覆盖？(y/n)${NC}"
        read -r OVERWRITE_CHOICE
        if [[ "$OVERWRITE_CHOICE" != "y" && "$OVERWRITE_CHOICE" != "Y" ]]; then
            info "保留现有 .env 文件"
            return
        fi
    fi

    cat > "$ENV_FILE" << EOF
# 智钓蓝海 — 环境配置文件
# 由 setup.sh 自动生成于 $(date '+%Y-%m-%d %H:%M:%S')

# 数据库配置
DB_HOST=${DB_HOST}
DB_PORT=${DB_PORT}
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=${DB_NAME}

# JWT 配置
JWT_SECRET=${JWT_SECRET}
JWT_ACCESS_EXPIRE_HOURS=1
JWT_REFRESH_EXPIRE_DAYS=3

# 服务配置
SERVER_PORT=${SERVER_PORT}
GIN_MODE=release
EOF

    ok ".env 文件已生成: $ENV_FILE"
}

# ===== 10. 导入示例数据（可选） =====
import_seed_data() {
    SEED_FILE="$SCRIPT_DIR/seed.sql"
    if [[ ! -f "$SEED_FILE" ]]; then
        info "未找到 seed.sql，跳过示例数据导入"
        return
    fi

    echo ""
    echo -e "${YELLOW}是否导入示例数据？这会清空现有数据并填充演示用数据。(y/n)${NC}"
    read -r SEED_CHOICE
    if [[ "$SEED_CHOICE" == "y" || "$SEED_CHOICE" == "Y" ]]; then
        info "正在导入示例数据 ..."
        mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$SEED_FILE"
        ok "示例数据导入完成"
    else
        info "跳过示例数据导入"
    fi
}

# ============================================================================
# 主流程
# ============================================================================

detect_os

# 步骤 1: 检查 MySQL
echo ""
info "步骤 1/5: 检查 MySQL"
if ! check_mysql; then
    install_mysql
fi

# 步骤 2: 检查 MySQL 服务
echo ""
info "步骤 2/5: 检查 MySQL 服务"
if ! check_mysql_service; then
    start_mysql_service
fi

# 步骤 3: 收集配置
echo ""
info "步骤 3/5: 配置数据库信息"
collect_config

# 步骤 4: 测试连接 + 创建数据库
echo ""
info "步骤 4/5: 连接数据库"
test_connection
create_database

# 步骤 5: 生成 .env
echo ""
info "步骤 5/5: 生成配置文件"
generate_env

# 可选: 导入示例数据
import_seed_data

# ===== 完成 =====
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   ✅ 初始化完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "  配置文件: ${CYAN}${ENV_FILE}${NC}"
echo -e "  数据库:   ${CYAN}${DB_NAME}${NC}"
echo ""
echo -e "  启动服务:"

# 检测可用的二进制文件
COMBO_BIN=$(ls "$SCRIPT_DIR"/smart-fish-combo-* 2>/dev/null | head -1)
SERVER_BIN=$(ls "$SCRIPT_DIR"/smart-fish-server-* 2>/dev/null | head -1)

if [[ -n "$COMBO_BIN" ]]; then
    COMBO_NAME=$(basename "$COMBO_BIN")
    echo -e "    ${CYAN}./${COMBO_NAME}${NC}"
elif [[ -n "$SERVER_BIN" ]]; then
    SERVER_NAME=$(basename "$SERVER_BIN")
    echo -e "    ${CYAN}./${SERVER_NAME}${NC}"
else
    echo -e "    ${CYAN}./smart-fish-combo-linux-amd64${NC}  (请先下载对应平台的二进制文件)"
fi

echo ""
echo -e "  首次运行加 ${YELLOW}--seed${NC} 参数可通过 Go 代码填充基础数据"
echo ""
