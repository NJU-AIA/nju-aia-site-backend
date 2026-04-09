// cmd/server/main.go
package server

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "ArticleServer/docs" // 导入 Swagger 生成的代码
	"ArticleServer/internal/article"
	"ArticleServer/internal/asset"
	"ArticleServer/internal/auth"
	"ArticleServer/internal/livecode"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Run() {
	// 0. 自动加载 .env（文件不存在时静默忽略，但必须的环境变量不能为空）
	_ = godotenv.Load()

	// 1. 初始化数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASS", "000000"), // 建议：生产环境下数据库密码也应使用 mustGetEnv
		getEnv("DB_HOST", "127.0.0.1"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "ArticleDB"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 2. 静态资源存储引擎选择
	var storageEngine asset.Storage
	if getEnv("STORAGE_ENGINE", "") == "local" {
		storageEngine = &asset.LocalStorage{
			RootPath: "./storage",
			BaseURL:  "http://localhost:8080/assets",
		}
	} else {
		storageEngine = &asset.COSStorage{
			BucketURL: getEnv("COS_BUCKET_URL", ""),
			SecretID:  getEnv("COS_SECRET_ID", ""),
			SecretKey: getEnv("COS_SECRET_KEY", ""),
		}
	}

	// 3. 模块初始化
	// --- 文章模块 ---
	articleRepo := article.NewRepository(db)
	articleSvc := article.NewService(articleRepo)
	articleHandler := article.NewHandler(articleSvc)

	// --- 静态资源模块 ---
	assetRepo := asset.NewRepository(db)
	assetSvc := asset.NewService(assetRepo, storageEngine)
	assetHandler := asset.NewHandler(assetSvc)

	// --- Livecode 模块 ---
	livecodeRepo := livecode.NewRepository(db)
	livecodeSvc := livecode.NewService(livecodeRepo)
	livecodeHandler := livecode.NewHandler(livecodeSvc)

	// --- 鉴权模块 (核心修改区) ---
	// 强制读取安全密钥，没有任何默认值，缺失直接崩溃
	jwtSecret := mustGetEnv("JWT_SECRET")
	dynamicSecret := mustGetEnv("DYNAMIC_SECRET") // 替换了之前的 TOTP_ECC_KEY

	authSvc := auth.NewService(
		jwtSecret,
		dynamicSecret,
		2*time.Hour,
	)
	authHandler := auth.NewHandler(authSvc)
	authMiddleware := auth.NewMiddleware(authSvc)

	// 4. 设置路由
	r := gin.Default()

	if _, ok := storageEngine.(*asset.LocalStorage); ok {
		r.Static("/assets", "./storage")
	}

	v1 := r.Group("/api")
	{
		// 登录接口 (现在只接收强动态口令)
		v1.POST("/auth/login", authHandler.LoginAdmin)

		// --- 公开访问 (无需鉴权) ---
		v1.GET("/articles", articleHandler.ListArticles)      // 获取文章列表
		v1.GET("/articles/:id", articleHandler.GetArticle)    // 获取文章详情
		v1.GET("/assets", assetHandler.ListAssets)            // 获取资源列表
		v1.GET("/livecodes", livecodeHandler.ListDocuments)   // 获取 livecode 列表
		v1.GET("/livecodes/:id", livecodeHandler.GetDocument) // 获取 livecode 详情

		// --- 敏感操作 ---
		authorized := v1.Group("/")
		authorized.Use(authMiddleware.RequireAdmin()) // 使用瘦身后的单点角色拦截器
		{
			// 文章管理
			authorized.POST("/articles", articleHandler.CreateArticle)       // 创建文章
			authorized.PUT("/articles/:id", articleHandler.UpdateArticle)    // 更新文章
			authorized.DELETE("/articles/:id", articleHandler.DeleteArticle) // 删除文章

			// Livecode 管理
			authorized.POST("/livecodes", livecodeHandler.CreateDocument)       // 创建 livecode
			authorized.PUT("/livecodes/:id", livecodeHandler.UpdateDocument)    // 更新 livecode
			authorized.DELETE("/livecodes/:id", livecodeHandler.DeleteDocument) // 删除 livecode

			// 资源管理
			authorized.POST("/assets", assetHandler.UploadFile)    // 上传资源
			authorized.DELETE("/assets", assetHandler.DeleteAsset) // 删除资源
		}
	}

	// 5. Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 6. 启动服务
	log.Println("启动服务，监听端口 :8080")
	r.Run(":8080")
}

// getEnv 用于非致命配置（允许有默认值）
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// mustGetEnv 用于核心安全配置（不允许有默认值，缺失即崩溃）
func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		log.Fatalf("致命安全错误: 必须在环境变量或 .env 文件中设置 %s", key)
	}
	return value
}
