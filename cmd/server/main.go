package server

import (
	"fmt"
	"log"
	"os"

	_ "ArticleServer/docs"
	"ArticleServer/internal/article"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBConfig 定义数据库配置
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

// 组装 DSN 字符串
func (c *DBConfig) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}


// @title AIA 社团博客 API
// @version 1.0
// @description AIA 论坛后端。
// @host localhost:8080
// @BasePath /api
func Run() {
	
	config := DBConfig{
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASS", "061112"),
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     getEnv("DB_PORT", "3306"),
		DBName:   getEnv("DB_NAME", "ArticleData"),
	}

	log.Printf("正在尝试连接数据库: %s:%s/%s", config.Host, config.Port, config.DBName)

	db, err := gorm.Open(mysql.Open(config.getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

    repo := article.NewRepository(db)
	svc := article.NewService(repo)
	handler := article.NewHandler(svc)

	r := gin.Default()
	v1 := r.Group("/api")
	{
		v1.POST("/articles", handler.CreateArticle)
		v1.GET("/articles/:id", handler.GetArticle)
		v1.GET("/articles", handler.ListArticles)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}