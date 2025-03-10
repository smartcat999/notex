package migrations

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// Migration 表示一个数据库迁移
type Migration struct {
	ID   string
	Up   string
	Down string
}

// RunMigrations 执行数据库迁移
func RunMigrations(db *gorm.DB, migrationsDir string) error {
	// 创建迁移表
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// 获取已应用的迁移
	var appliedMigrations []string
	if err := db.Model(&struct{ ID string }{}).Table("migrations").Pluck("id", &appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	// 读取迁移文件
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	var migrations []Migration
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			id := strings.TrimSuffix(file.Name(), ".up.sql")
			upPath := filepath.Join(migrationsDir, file.Name())
			downPath := filepath.Join(migrationsDir, id+".down.sql")

			upSQL, err := os.ReadFile(upPath)
			if err != nil {
				return fmt.Errorf("failed to read up migration %s: %v", id, err)
			}

			downSQL, err := os.ReadFile(downPath)
			if err != nil {
				return fmt.Errorf("failed to read down migration %s: %v", id, err)
			}

			migrations = append(migrations, Migration{
				ID:   id,
				Up:   string(upSQL),
				Down: string(downSQL),
			})
		}
	}

	// 按ID排序迁移
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID < migrations[j].ID
	})

	// 应用未执行的迁移
	for _, migration := range migrations {
		applied := false
		for _, appliedID := range appliedMigrations {
			if appliedID == migration.ID {
				applied = true
				break
			}
		}

		if !applied {
			log.Printf("Applying migration: %s\n", migration.ID)
			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Exec(migration.Up).Error; err != nil {
					return fmt.Errorf("failed to apply migration %s: %v", migration.ID, err)
				}
				if err := tx.Exec("INSERT INTO migrations (id) VALUES (?)", migration.ID).Error; err != nil {
					return fmt.Errorf("failed to record migration %s: %v", migration.ID, err)
				}
				return nil
			}); err != nil {
				return err
			}
			log.Printf("Successfully applied migration: %s\n", migration.ID)
		}
	}

	return nil
}

// RollbackMigrations 回滚数据库迁移
func RollbackMigrations(db *gorm.DB, migrationsDir string) error {
	// 获取已应用的迁移
	var appliedMigrations []string
	if err := db.Model(&struct{ ID string }{}).Table("migrations").Pluck("id", &appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	// 读取迁移文件
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	var migrations []Migration
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".down.sql") {
			id := strings.TrimSuffix(file.Name(), ".down.sql")
			upPath := filepath.Join(migrationsDir, id+".up.sql")
			downPath := filepath.Join(migrationsDir, file.Name())

			upSQL, err := os.ReadFile(upPath)
			if err != nil {
				return fmt.Errorf("failed to read up migration %s: %v", id, err)
			}

			downSQL, err := os.ReadFile(downPath)
			if err != nil {
				return fmt.Errorf("failed to read down migration %s: %v", id, err)
			}

			migrations = append(migrations, Migration{
				ID:   id,
				Up:   string(upSQL),
				Down: string(downSQL),
			})
		}
	}

	// 按ID倒序排序迁移
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID > migrations[j].ID
	})

	// 回滚已执行的迁移
	for _, migration := range migrations {
		applied := false
		for _, appliedID := range appliedMigrations {
			if appliedID == migration.ID {
				applied = true
				break
			}
		}

		if applied {
			log.Printf("Rolling back migration: %s\n", migration.ID)
			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Exec(migration.Down).Error; err != nil {
					return fmt.Errorf("failed to rollback migration %s: %v", migration.ID, err)
				}
				if err := tx.Exec("DELETE FROM migrations WHERE id = ?", migration.ID).Error; err != nil {
					return fmt.Errorf("failed to remove migration record %s: %v", migration.ID, err)
				}
				return nil
			}); err != nil {
				return err
			}
			log.Printf("Successfully rolled back migration: %s\n", migration.ID)
		}
	}

	return nil
}
