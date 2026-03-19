# 🐟 智钓蓝海 — 智能垂钓信息平台

> **边缘-云协同的全生态智能垂钓系统信息平台**
>
> 实时采集水域环境数据（水温、气温、湿度、水质等），提供垂钓热度排行、智能建议、提醒告警等功能，面向垂钓爱好者和水域管理人员。

## 系统架构

```
┌────────────────────────────────────────────────────────────┐
│                       客户端（浏览器）                       │
│  Vue 3 + TypeScript + Element Plus + ECharts + Pinia       │
│  桌面端 & 移动端响应式 · 内置 Mock 数据可独立运行             │
└────────────────────┬───────────────────────────────────────┘
                     │ HTTP (Vite Proxy / Nginx)
                     ▼
┌────────────────────────────────────────────────────────────┐
│                    Go 后端（Gin + GORM）                     │
│  ┌──────────┐ ┌──────────┐ ┌───────────┐ ┌─────────────┐  │
│  │ 认证模块  │ │ 区域管理  │ │ 水域管理   │ │  设备/网关   │  │
│  │ JWT 双Token│ │ 省/市树形 │ │ CRUD+收藏  │ │  状态监控    │  │
│  └──────────┘ └──────────┘ └───────────┘ └─────────────┘  │
│  ┌──────────┐ ┌──────────┐ ┌───────────┐ ┌─────────────┐  │
│  │ 提醒系统  │ │ 通知中心  │ │ 垂钓建议   │ │  数据上传    │  │
│  │ 告警处理  │ │ 系统公告  │ │ 智能推荐   │ │  IoT 数据    │  │
│  └──────────┘ └──────────┘ └───────────┘ └─────────────┘  │
│  中间件：CORS · JWT 认证 · 三级权限（user/staff/admin）       │
└────────────────────┬───────────────────────────────────────┘
                     │ GORM ORM
                     ▼
┌────────────────────────────────────────────────────────────┐
│                   MySQL 8.0+（11 张表）                      │
│  User · Region · FishingSpot · Device · Gateway            │
│  HistoricalData · EnvironmentData · WaterQualityData       │
│  Reminder · Notice · FishingSuggestion                     │
└────────────────────────────────────────────────────────────┘
```

## 项目结构

```
smart-fish/
├── back_end/                   # Go 后端
│   ├── config/                 # 配置管理（.env 加载）
│   ├── database/               # 数据库连接、自动迁移、种子数据
│   ├── handlers/               # API 处理器（10 个模块）
│   ├── middleware/              # 中间件（CORS、认证、权限）
│   ├── models/                 # 数据模型（11 个）
│   ├── routes/                 # 路由注册 + 前端静态文件服务
│   ├── services/               # 业务服务（JWT Token）
│   ├── static.go               # 组合构建模式（embed 前端）
│   ├── static_stub.go          # 纯后端构建模式
│   └── main.go                 # 入口
├── front_end/                  # Vue 3 前端
│   ├── src/
│   │   ├── components/         # 可复用组件
│   │   │   ├── chart/          # BaseChart 图表封装
│   │   │   ├── data/           # SpotList 数据展示
│   │   │   └── layout/         # Navbar、Footer
│   │   ├── views/              # 7 个页面视图 + 子组件
│   │   │   ├── home/           # 首页模块（Banner、统计、趋势图等）
│   │   │   ├── spots/          # 水域模块（卡片、筛选、区域树等）
│   │   │   ├── reminder/       # 信息中心（提醒、通知、建议）
│   │   │   ├── admin/          # 管理后台（4 个管理表格）
│   │   │   ├── auth/           # 登录/注册表单
│   │   │   └── user/           # 个人中心（资料、密码、收藏）
│   │   ├── stores/             # 7 个 Pinia Store
│   │   ├── services/           # API 服务 + MockDataService
│   │   ├── network/            # Axios 封装
│   │   ├── types/              # TypeScript 类型定义
│   │   ├── router/             # Vue Router + 路由守卫
│   │   └── plugins/            # Element Plus 插件配置
│   └── vite.config.ts          # Vite 配置（代理、别名）
├── .github/workflows/          # GitHub Actions CI/CD
│   └── release.yml             # 自动构建 & 发布流水线
├── start.sh                    # 一键启动脚本
└── README.md
```

## 快速开始

### 环境要求

| 工具 | 最低版本 | 说明 |
|------|---------|------|
| Go | 1.21+ | 后端运行环境 |
| Node.js | 18+ | 前端构建环境 |
| npm | 9+ | 前端包管理 |
| MySQL | 8.0+ | 数据库（仅后端需要） |

### 方式一：使用启动脚本（推荐）

```bash
chmod +x start.sh
./start.sh
```

脚本会自动检测环境、安装依赖并启动服务，支持三种模式：
- **前后端一起启动** — 完整功能
- **仅启动前端** — 使用内置 Mock 数据，无需数据库
- **仅启动后端** — API 调试

### 方式二：手动启动

#### 1. 启动后端

```bash
cd back_end

# 创建 .env 文件（可选，有默认值）
cat > .env << EOF
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=smart_fish
JWT_SECRET=your-secret-key
SERVER_PORT=8080
EOF

# 创建数据库
mysql -u root -e "CREATE DATABASE IF NOT EXISTS smart_fish CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 下载依赖
go mod download

# 首次启动加 --seed 填充示例数据
go run . --seed

# 后续启动
go run .
```

后端默认运行在 `http://localhost:8080`

#### 2. 启动前端

```bash
cd front_end
npm install
npm run dev
```

前端默认运行在 `http://localhost:5173`

> **无后端也能跑**：前端内置了完整的 Mock 数据服务 (`MockDataService.ts`)，当后端不可用时自动回退到模拟数据。

### 方式三：仅查看前端

```bash
cd front_end
npm install
npm run dev
# 打开 http://localhost:5173 即可查看完整前端效果（Mock 数据）
```

## 构建与部署

### 前端构建

```bash
cd front_end
npm run build     # TypeScript 类型检查 + Vite 构建，输出到 dist/
npm run preview   # 本地预览生产构建
```

### 后端构建

```bash
cd back_end

# 纯后端模式（仅 API，不含前端）
go build -tags backend_only -ldflags="-s -w" -o smart-fish-server .

# 组合模式（前端嵌入后端，单文件运行）
# 需要先将前端 dist/ 放到 back_end/dist/
cd ../front_end && npm run build && cp -r dist/ ../back_end/dist/
cd ../back_end && go build -ldflags="-s -w" -o smart-fish-combo .
```

组合模式下，运行单个二进制文件即可同时提供 API 和前端页面，支持 Vue Router SPA 路由。

### CI/CD 自动发布

项目配置了 GitHub Actions 流水线 (`.github/workflows/release.yml`)，推送到 `main` 分支时自动触发：

| 步骤 | 说明 |
|------|------|
| 自动打 Tag | 递增 patch 版本号（v1.0.0 → v1.0.1） |
| 构建前端 | Node.js 20 构建，打包为 tar.gz/zip |
| 构建纯后端 | 交叉编译 5 个平台（Linux/macOS/Windows × amd64/arm64） |
| 构建组合版 | 前端嵌入后端，5 个平台的一体化二进制 |
| 创建 Release | 附带所有构建产物和部署说明 |

**发布产物：**

| 产物 | 说明 |
|------|------|
| `smart-fish-combo-*` | 组合版（前端 + 后端一体，单文件运行） |
| `smart-fish-server-*` | 纯后端 API 服务 |
| `smart-fish-frontend-dist.*` | 前端静态文件包（tar.gz / zip） |

### 生产部署（组合版）

```bash
# 1. 创建 .env
cat > .env << EOF
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=smart_fish
JWT_SECRET=your-production-secret
SERVER_PORT=8080
GIN_MODE=release
EOF

# 2. 首次运行（含示例数据）
chmod +x smart-fish-combo-linux-amd64
./smart-fish-combo-linux-amd64 --seed

# 3. 访问 http://localhost:8080
```

### 生产部署（分离模式）

1. 解压 `smart-fish-frontend-dist.zip` 到 Nginx 等 Web 服务器
2. 配置反向代理，将 `/api/*` 和 `/auth/*` 转发到后端端口
3. 后端运行 `smart-fish-server-*`，仅提供 API 服务

## 技术栈

### 后端

| 技术 | 说明 |
|------|------|
| **Go 1.25** | 编程语言 |
| **Gin** | Web 框架 |
| **GORM** | ORM 框架 |
| **MySQL** | 关系型数据库 |
| **JWT (HS256)** | 双 Token 认证（Access 1h + Refresh 3d） |
| **godotenv** | 环境变量配置 |
| **bcrypt** | 密码加密 |

### 前端

| 技术 | 说明 |
|------|------|
| **Vue 3.5** | UI 框架（Composition API） |
| **TypeScript 5.9** | 类型安全 |
| **Vite 8** | 构建工具 |
| **Element Plus** | UI 组件库 |
| **Pinia** | 状态管理（7 个 Store） |
| **ECharts 6** | 数据可视化 |
| **Axios** | HTTP 客户端 |
| **Vue Router** | 前端路由（History 模式） |

## 功能页面

| 页面 | 路由 | 权限 | 功能 |
|------|------|------|------|
| 首页 | `/` | 公开 | 数据统计卡片、热门水域排行、垂钓建议、趋势图、环境仪表盘、提醒面板、通知公告 |
| 垂钓水域 | `/spots` | 公开 | 区域树导航、卡片/表格双视图、状态筛选搜索、负载率进度条、详情弹窗（含环境折线图） |
| 信息中心 | `/reminders` | 公开 | 提醒/通知/建议三 Tab、级别筛选、处理状态管理 |
| 登录/注册 | `/auth` | 公开 | 表单验证、Tab 切换、密码确认 |
| 个人中心 | `/profile` | 登录 | 个人信息编辑、修改密码、收藏列表 |
| 管理后台 | `/admin` | Staff+ | 区域/水域/设备/网关的 CRUD 管理 |

## API 参考

### 认证

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| POST | `/api/auth/register` | 公开 | 用户注册 |
| POST | `/api/auth/login` | 公开 | 用户登录 |
| POST | `/api/auth/refresh` | 公开 | 刷新 Token |
| GET | `/api/auth/me` | 登录 | 获取当前用户信息 |
| PUT | `/api/auth/me` | 登录 | 更新个人信息 |
| PUT | `/api/auth/password` | 登录 | 修改密码 |

### 区域管理

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/regions` | 公开 | 区域列表（树形） |
| GET | `/api/regions/provinces` | 公开 | 省份列表 |
| GET | `/api/regions/environment` | 公开 | 区域环境概览 |
| GET | `/api/regions/:id` | 公开 | 区域详情 |
| GET | `/api/regions/:id/environment` | 公开 | 区域环境历史 |
| POST | `/api/regions` | Staff+ | 创建区域 |
| PUT | `/api/regions/:id` | Staff+ | 更新区域 |
| DELETE | `/api/regions/:id` | Staff+ | 删除区域 |

### 水域管理

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/spots` | 公开 | 水域列表（分页+筛选） |
| GET | `/api/spots/popular` | 公开 | 热门水域排行 |
| GET | `/api/spots/:id` | 公开 | 水域详情 |
| GET | `/api/spots/:id/historical` | 公开 | 垂钓人数历史 |
| GET | `/api/spots/:id/environment` | 公开 | 环境数据 |
| POST | `/api/spots/:id/favor` | 登录 | 收藏/取消收藏 |
| GET | `/api/spots/favorites` | 登录 | 我的收藏 |
| POST | `/api/spots` | Staff+ | 创建水域 |
| PUT | `/api/spots/:id` | Staff+ | 更新水域 |
| DELETE | `/api/spots/:id` | Staff+ | 删除水域 |

### 设备 & 网关

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/devices` | 公开 | 设备列表 |
| GET | `/api/devices/:id` | 公开 | 设备详情 |
| POST/PUT/DELETE | `/api/devices[/:id]` | Staff+ | 设备 CRUD |
| GET | `/api/gateways` | 公开 | 网关列表 |
| GET | `/api/gateways/:id` | 公开 | 网关详情 |
| POST/PUT/DELETE | `/api/gateways[/:id]` | Staff+ | 网关 CRUD |

### 提醒 & 通知 & 建议

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/reminders` | 公开 | 提醒列表 |
| PATCH | `/api/reminders/:id/resolve` | Staff+ | 处理提醒 |
| GET | `/api/notices` | 公开 | 通知列表 |
| GET | `/api/suggestions` | 公开 | 垂钓建议列表 |
| GET | `/api/suggestions/latest` | 公开 | 最新建议 |

### 其他

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/summary` | 公开 | 系统统计概览 |
| POST | `/api/upload/fishing-data` | 登录 | 上传垂钓数据 |
| POST | `/api/upload/environment` | 登录 | 上传环境数据 |
| POST | `/api/upload/water-quality` | 登录 | 上传水质数据 |
| POST | `/api/upload/device-status` | 登录 | 上传设备状态 |
| GET | `/api/users` | Admin | 用户列表 |
| PATCH | `/api/users/:id/role` | Admin | 修改用户角色 |
| DELETE | `/api/users/:id` | Admin | 删除用户 |

## 数据模型

| 模型 | 说明 | 关键字段 |
|------|------|---------|
| **User** | 用户 | username, password(bcrypt), role(user/staff/admin) |
| **Region** | 垂钓区域 | name, province, city, parent_id |
| **FishingSpot** | 垂钓水域 | name, type, capacity, status, region_id, bound_device |
| **Device** | 监测设备 | name, type(environment/underwater/fish_finder), status |
| **Gateway** | 物联网网关 | name, ip, status, connected_devices |
| **HistoricalData** | 垂钓人数记录 | spot_id, fishing_count, recorded_at |
| **EnvironmentData** | 环境监测 | water_temp, air_temp, humidity, pressure |
| **WaterQualityData** | 水质数据 | ph, dissolved_oxygen, turbidity |
| **Reminder** | 提醒信息 | type(weather/capacity/environment/info), level, status |
| **Notice** | 系统通知 | title, content, type |
| **FishingSuggestion** | 垂钓建议 | content, spot_id, score |

## Mock 数据

前端 `src/services/MockDataService.ts` 提供完整的模拟数据，**无需后端即可体验全部功能**：

- 4 个模拟区域（黑龙江省 — 哈尔滨 / 牡丹江 / 鸡西）
- 6 个模拟水域（松花江 / 太阳岛 / 镜泊湖 / 兴凯湖 / 呼兰河 / 松北鱼塘）
- 6 台模拟设备（环境监测 / 水下感知 / 探鱼辅助）
- 4 条提醒（天气 / 容量 / 环境 / 信息四种类型）
- 3 条通知、2 条垂钓建议
- 24 小时历史数据和环境数据动态生成（含日变化模拟）

每个页面都实现了"先请求真实 API，失败后自动回退到 Mock 数据"的降级机制。

## 环境变量

后端支持的环境变量（通过 `.env` 文件或系统环境变量配置）：

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `DB_HOST` | `127.0.0.1` | MySQL 主机 |
| `DB_PORT` | `3306` | MySQL 端口 |
| `DB_USER` | `root` | MySQL 用户名 |
| `DB_PASSWORD` | _(空)_ | MySQL 密码 |
| `DB_NAME` | `smart_fish` | 数据库名 |
| `JWT_SECRET` | `default-secret-change-me` | JWT 签名密钥（**生产环境务必修改**） |
| `JWT_ACCESS_EXPIRE_HOURS` | `1` | Access Token 有效期（小时） |
| `JWT_REFRESH_EXPIRE_DAYS` | `3` | Refresh Token 有效期（天） |
| `SERVER_PORT` | `8080` | 后端监听端口 |
| `GIN_MODE` | `debug` | Gin 运行模式（debug / release） |

## 许可证

本项目仅供学习交流使用。
