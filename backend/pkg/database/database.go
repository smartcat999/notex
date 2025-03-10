package database

import (
	"fmt"
	"notex/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB 全局数据库连接
var DB *gorm.DB

// Initialize 初始化数据库连接
func Initialize(cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	DB = db
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}
