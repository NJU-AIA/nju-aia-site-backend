# AIA 后端项目使用说明

## 1. 环境依赖

* **运行时**: Go 1.20+
* **数据库**: MySQL 8.0+
* **文档工具**: Swag CLI (用于自动生成 OpenAPI/Swagger 契约)
    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

## 2. 快速启动

### 步骤 A：初始化数据库
在 MySQL 中创建逻辑库（具体表结构将在程序首次启动时由 GORM 自动迁移）：
```sql
CREATE DATABASE IF NOT EXISTS ArticleDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 步骤 B：配置环境变量 (核心安全步骤)
复制环境变量模板：
```bash
cp .env.example .env
```
你需要按照以下方式生成并填入 `.env` 文件：

1.  **生成 DYNAMIC_SECRET**:
    直接运行动态密码生成器：
    ```bash
    go run cmd/totp/main.go
    ```
    程序会检测到你尚未配置环境变量，并为你自动生成一个 32 字节的高熵 Base64 密钥。将该密钥填入 `.env` 的 `DYNAMIC_SECRET` 字段。
2.  **生成 JWT_SECRET**:
    在终端执行以下命令生成随机高强度密钥，并填入 `.env` 的 `JWT_SECRET` 字段：
    ```bash
    openssl rand -base64 32
    ```

### 步骤 C：编译依赖与同步契约
```bash
# 整理并下载依赖
go mod tidy

# 扫描代码注释，生成最新 Swagger 文档
swag init
```

### 步骤 D：启动服务
```bash
go run main.go
```
程序会自动加载项目根目录下的 `.env` 文件。如果缺失核心密钥，程序将立即崩溃以防止裸奔启动。

---

## 3. 环境变量配置表

| 变量名 | 说明 | 默认值 | 缺失行为 |
| :--- | :--- | :--- | :--- |
| `JWT_SECRET` | JWT 签名密钥（建议 Base64 格式，至少 32 字节） | **无** | 🚨 **致命崩溃** |
| `DYNAMIC_SECRET`| 终极管理员动态口令密钥（必须是 Base64 格式） | **无** | 🚨 **致命崩溃** |
| `DB_USER` | 数据库用户名 | `root` | 采用默认值 |
| `DB_PASS` | 数据库密码 | `000000` | 采用默认值 |
| `DB_HOST` | 数据库物理地址 | `127.0.0.1` | 采用默认值 |
| `DB_PORT` | 数据库监听端口 | `3306` | 采用默认值 |
| `DB_NAME` | 逻辑数据库名称 | `ArticleDB` | 采用默认值 |
| `STORAGE_ENGINE`| 存储引擎选择 (`local` 或 `cos`) | `local` | 采用默认值 |

> **🛡️ 安全提示**: `.env` 文件已被 `.gitignore` 排除，**绝对不会推送到代码仓库**。仅 `.env.example`（不含实际密钥）会被提交。

---

## 4. 鉴权规范 (Authentication)

本项目没有任何用户表，系统仅允许**唯一终极管理员**登录。

### 登录与动态码原理
1.  在本地终端配置好相同的 `DYNAMIC_SECRET` 环境变量。
2.  运行 `go run cmd/totp/main.go`（建议编译为本地可执行文件）。
3.  生成器会基于 HMAC-SHA256 和当前时间戳，每 30 秒严格同步生成一个 **16 位的强动态密码**（例如：`x6+xJ/c5fsWaVyxu`）。
4.  调用 `POST /api/auth/login`，在 Body 中携带该口令：`{"password": "动态密码"}`。
5.  获取 JWT Token 并放入请求进行这样的设计：
代码表只记录一个blocks的id列表，这样就不用维护order了，列表自带顺序
blocks都记录在另一张表，用id标识，不管自己属于哪一个代码。

这样前端修改顺序和添加内容都方便。


### 安全特性
* **0 容差时间同步**: 口令严格按 30 秒窗口计算，过期 1 秒即失效。
* **单设备活跃 (顶号机制)**: 在新设备成功使用动态口令登录后，旧设备的 Token 会被瞬间作废，防止 Token 泄露风险。

---

## 5. 接口契约说明 (API Contract)

### 鉴权模块 (Auth)
| 功能 | 方法 | 路径 | 鉴权要求 | 成功码 |
| :--- | :--- | :--- | :--- | :--- |
| **管理员登录** | `POST` | `/api/auth/login` | 公开 | `200` |

### 文章模块 (Articles)
| 功能 | 方法 | 路径 | 鉴权要求 | 成功码 |
| :--- | :--- | :--- | :--- | :--- |
| **获取文章列表** | `GET` | `/api/articles` | 公开 | `200` |
| **获取文章详情** | `GET` | `/api/articles/:id` | 公开 | `200` |
| **发布文章** | `POST` | `/api/articles` | 🔐 管理员 | `201` |
| **编辑文章** | `PUT` | `/api/articles/:id` | 🔐 管理员 | `200` |
| **删除文章** | `DELETE` | `/api/articles/:id` | 🔐 管理员 | `204` |

### 资源模块 (Assets)
| 功能 | 方法 | 路径 | 鉴权要求 | 成功码 |
| :--- | :--- | :--- | :--- | :--- |
| **获取资源列表** | `GET` | `/api/assets` | 公开 | `200` |
| **上传文件资源** | `POST` | `/api/assets` | 🔐 管理员 | `201` |
| **删除文件资源** | `DELETE` | `/api/assets` | 🔐 管理员 | `204` |

---

## 6. 存储与访问规范

### 物理存储
* **根路径**: `./storage/` (程序运行时若不存在会自动创建)
* **清理逻辑**: 调用删除接口时，将同步移除磁盘/云端的物理文件，防止产生“孤儿文件”。

### 静态访问
* **本地映射**: 当 `STORAGE_ENGINE=local` 时，访问 `http://{host}:8080/assets/` 将由 Gin 路由直接映射至 `./storage/`。

---

## 7. 项目物理结构

```text
├── cmd/
│   ├── server/          # 核心调度中心：依赖注入、路由挂载、服务启动
│   └── totp/            # 本地独立 CLI 工具：终极动态强密码生成器
├── internal/
│   ├── article/         # 文章模块：内容管理、Markdown 存储等
│   ├── asset/           # 资源模块：支持 Local/COS 双引擎切换、物理文件管理
│   └── auth/            # 鉴权模块：高熵动态口令计算、JWT 签发、单点登录防顶号拦截
├── storage/             # 本地物理存储根目录 (运行时自动创建)
├── docs/                # OpenAPI/Swagger 静态文档自动生成目录
├── .env.example         # 环境变量模板（安全，可提交）
├── .env                 # 实际环境变量（已 gitignore，存放真实敏感密钥）
├── go.mod               # Go 模块依赖描述
└── README.md            # 项目说明文档
```