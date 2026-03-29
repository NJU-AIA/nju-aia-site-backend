package main

import "ArticleServer/cmd/server"

// @title AIA 社团博客 API
// @version 1.0
// @description 这是 AIA 社团论坛的后端服务，支持 Markdown 文章管理与静态资源分类存储。
// @host localhost:8080
// @BasePath /api
func main() {
	// 调用 cmd/server 包中的 Run 函数启动程序
	server.Run()
}
