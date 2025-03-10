package main

import (
	"fmt"
	"log"
	"notex/api/router"
	"notex/config"
	"notex/middleware"
	"notex/migrations"
	"notex/pkg/database"
	"notex/pkg/email"
	"path/filepath"
)

func main() {
	// 加载配置文件
	configPath := filepath.Join("config", "config.yaml")
	if err := config.LoadConfig(configPath); err != nil {
		log.Printf("Warning: Failed to load config file: %v", err)
		log.Println("Using default configuration")
	}

	cfg := config.GetConfig()

	// 初始化数据库连接
	if err := database.Initialize(cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 执行数据库迁移
	if err := migrations.RunMigrations(database.GetDB(), "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 初始化邮件发送器
	if err := email.Initialize(cfg.Email); err != nil {
		log.Fatalf("Failed to initialize email system: %v", err)
	}

	// 初始化限流器
	middleware.InitRateLimiters(&cfg.RateLimit)

	// 设置路由
	r := router.SetupRouter(cfg)

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(r.Run(serverAddr))
}
