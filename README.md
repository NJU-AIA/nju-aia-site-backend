# ArticleServer 使用说明

## 1. 环境依赖

*   **运行时**: Go 1.20+
*   **数据库**: MySQL 8.0+
*   **文档工具**: Swag CLI (用于自动生成 Swagger 契约)
    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

## 2. 快速启动

### 步骤 A：初始化数据库
在 MySQL 中创建逻辑库（表结构由程序启动时自动迁移）：
```sql
CREATE DATABASE IF NOT EXISTS BlogData;
```

### 步骤 B：配置环境变量
复制环境变量模板并填写实际值：
```bash
cp .env.example .env
```
首次使用需生成密钥：
```bash
# 生成 TOTP 密钥（运行两次，分别填入 TOTP_ECC_KEY 和 JWT_SECRET）
go run cmd/totp/main.go -generate-key
```

### 步骤 C：编译依赖与同步契约
```bash
# 整理依赖
go mod tidy

# 生成 Swagger 文档
swag init
```

### 步骤 D：启动服务
```bash
go run main.go
```
程序会自动加载项目根目录下的 `.env` 文件。

## 3. 环境变量配置表 (Environment Variables)

| 变量名 | 说明 | 默认值 | 安全敏感 |
| :--- | :--- | :--- | :--- |
| `JWT_SECRET` | JWT 签名密钥（base64） | `change-me-in-production` | **是** |
| `TOTP_ECC_KEY` | ECC P-256 TOTP 动态码密钥（base64） | 无 | **是** |
| `DB_USER` | 数据库用户名 | `root` | 是 |
| `DB_PASS` | 数据库密码 | `114514` | 是 |
| `DB_HOST` | 数据库物理地址 | `127.0.0.1` | 否 |
| `DB_PORT` | 数据库监听端口 | `3306` | 否 |
| `DB_NAME` | 逻辑数据库名称 | `BlogData` | 否 |

> ⚠️ `.env` 文件已被 `.gitignore` 排除，**不会推送到 GitHub**。仅 `.env.example`（不含实际密钥）会被提交。

## 4. 鉴权说明 (Authentication)

### 登录流程
1. 运行 `go run cmd/totp/main.go` 获取当前 6 位动态码（每 30 秒刷新）
2. 调用 `POST /api/auth/login`，Body 为 `{"password": "<动态码>"}`
3. 返回 JWT Token，2 小时有效

### 动态码原理
基于 ECC P-256 私钥 + HMAC-SHA256 + 当前时间窗口（30 秒），生成 6 位一次性密码。
持有相同密钥的团队成员均可生成有效动态码。

## 5. 接口契约说明 (API Contract)

### 鉴权模块 (Auth)
| 功能 | 方法 | 路径 | 鉴权 | 成功码 |
| :--- | :--- | :--- | :--- | :--- |
| **登录** | `POST` | `/api/auth/login` | 否 | `200` |

### 文章模块 (Articles)
| 功能 | 方法 | 路径 | 鉴权 | 成功码 |
| :--- | :--- | :--- | :--- | :--- |
| **文章列表** | `GET` | `/api/articles` | 否 | `200` |
| **文章详情** | `GET` | `/api/articles/:id` | 否 | `200` |
| **发布文章** | `POST` | `/api/articles` | 是 | `201` |
| **编辑文章** | `PUT` | `/api/articles/:id` | 是 | `200` |
| **删除文章** | `DELETE` | `/api/articles/:id` | 是 | `204` |

### 资源模块 (Assets)
| 功能 | 方法 | 路径 | 鉴权 | 成功码 |
| :--- | :--- | :--- | :--- | :--- |
| **资源列表** | `GET` | `/api/assets` | 否 | `200` |
| **上传资源** | `POST` | `/api/assets` | 是 | `201` |
| **删除资源** | `DELETE` | `/api/assets` | 是 | `204` |

*注：受限接口需在 Header 中携带 `Authorization: Bearer {token}`。*

## 6. 存储与访问规范

### 物理存储
*   **根路径**: `./storage/assets/`
*   **分卷逻辑**: `/{scope}/{filename}`
*   **自动清理**: 删除接口会同步移除磁盘文件。

### 静态访问
*   **映射规则**: 访问 `http://{host}:8080/assets/` 将映射至 `./storage/assets/`。

## 7. 项目物理结构

```text
├── cmd/
│   ├── server/          # 环境调度中心、依赖注入、静态资源路由映射
│   └── totp/            # TOTP 动态码工具（生成密钥 / 获取当前码）
├── internal/
│   ├── article/         # 文章模块：支持小驼峰契约、UUID标识、Markdown存储
│   ├── asset/           # 资源模块：支持本地/OSS双引擎切换、物理文件管理
│   └── auth/            # 鉴权模块：ECC-TOTP 动态码 + JWT + RBAC 中间件
├── storage/             # 物理存储根目录 (运行时自动创建)
├── docs/                # OpenAPI/Swagger 静态文档
├── .env.example         # 环境变量模板（安全，可提交）
├── .env                 # 实际环境变量（已 gitignore，不提交）
├── main.go              # 程序唯一入口
├── go.mod               # 模块依赖描述
└── README.md
```