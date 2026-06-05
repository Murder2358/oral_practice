package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"oral_practice/internal/handler"
	"oral_practice/internal/repository"
	"oral_practice/internal/service"
	"oral_practice/pkg/llm"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 加载 .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}

	// 加载 config.yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// 初始化数据库
	dbPath := viper.GetString("database.path")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	repo := repository.NewRepository(db)
	if err := repo.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// 初始化 LLM 客户端
	llmClient := llm.NewClient()

	// 初始化服务
	convService := service.NewConversationService(llmClient)

	// 初始化处理器
	sceneHandler := handler.NewSceneHandler(repo)
	chatHandler := handler.NewChatHandler(convService)

	// 设置 Gin
	mode := viper.GetString("server.mode")
	gin.SetMode(mode)

	r := gin.Default()

	// 静态文件（前端构建产物）
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/", "./web/dist/index.html")

	// API 路由
	api := r.Group("/api")
	{
		api.GET("/scenes", sceneHandler.ListScenes)
		api.POST("/sessions", sceneHandler.CreateSession)
		api.GET("/sessions/:id", sceneHandler.GetSession)
	}

	// WebSocket
	r.GET("/ws", chatHandler.HandleWebSocket)

	// 启动服务
	port := viper.GetString("server.port")
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
