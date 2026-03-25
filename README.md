## 1. 环境依赖

*   **运行时**: Go 1.18+
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

### 步骤 B：编译依赖与同步契约
在项目根目录下执行，确保代码注释与生成的 JSON 文档完全一致：
```bash
# 整理依赖
go mod tidy

# 生成 Swagger 文档 (强制解析 internal 依赖)
swag init -g cmd/server/main.go --parseDependency
```

### 步骤 C：启动服务
支持通过环境变量注入配置。若在本地开发环境，可直接运行：
*   **Windows (PowerShell)**:
    ```powershell
    $env:DB_PASS="你的密码"; go run main.go
    ```
*   **Linux/Mac**:
    ```bash
    DB_PASS=你的密码 go run main.go
    ```
*   **默认配置路径**: `root:114514@tcp(127.0.0.1:3306)/BlogData`

## 3. 环境变量配置表 (Environment Variables)

| 变量名 | 说明 | 默认值 |
| :--- | :--- | :--- |
| `DB_USER` | 数据库用户名 | `root` |
| `DB_PASS` | 数据库密码 | `114514` |
| `DB_HOST` | 数据库物理地址 | `127.0.0.1` |
| `DB_PORT` | 数据库监听端口 | `3306` |
| `DB_NAME` | 逻辑数据库名称 | `BlogData` |

## 4. 接口契约说明 (API Contract)

本服务严格遵循 HTTP 原生状态码。响应体不含冗余 `code` 字段。

| 功能 | 方法 | 路径 | 成功状态码 | 错误码 (4xx) |
| :--- | :--- | :--- | :--- | :--- |
| **发布文章** | `POST` | `/api/v1/blogs` | `201 Created` | `400` (参数缺失) |
| **文章详情** | `GET` | `/api/v1/blogs/:id` | `200 OK` | `400` (ID非数字), `404` (不存在) |
| **文章列表** | `GET` | `/api/v1/blogs` | `200 OK` | `500` (数据库异常) |
| **可视化文档** | `GET` | `/swagger/index.html` | - | - |


## 5. 项目物理结构

```text
├── cmd/
│   └── server/          # 装配中心：数据库连接、CORS中间件、路由挂载
├── internal/
│   └── blog/            # 业务核心：Handler(接口)、Service(逻辑)、Repo(SQL)、Model(实体)
├── docs/                # Swag 自动生成的静态文档 (由 swag init 生成)
├── main.go              # 程序唯一物理入口
├── README.md
├── go.sum
└── go.mod               # 模块依赖描述
```