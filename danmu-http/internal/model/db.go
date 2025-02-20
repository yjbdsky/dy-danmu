package model

import (
	"danmu-http/logger"
	"danmu-http/setting"
	"danmu-http/utils"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var dsn = fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s sslmode=disable search_path=%s",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.DBName,
		setting.DatabaseSetting.Port,
		setting.DatabaseSetting.SearchPath)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("db.Setup failure: %v", err)
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get Sql DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(setting.DatabaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(setting.DatabaseSetting.MaxOpenConns)

	// 初始化默认管理员账户
	if setting.AppSetting.AdminEmail != "" && setting.AppSetting.AdminPassword != "" {
		if err := initDefaultAdmin(); err != nil {
			logger.Error().Err(err).Msg("init default admin failed")
		}
	}
}

func initDefaultAdmin() error {

	// 创建默认管理员账户
	hashedPassword, err := utils.HashPassword(setting.AppSetting.AdminPassword)
	if err != nil {
		return fmt.Errorf("hash default admin password failed: %w", err)
	}

	admin := &Auth{
		Email:    setting.AppSetting.AdminEmail,
		Password: hashedPassword,
		Name:     "admin",
		Role:     "admin",
	}

	if err := admin.InsertOrUpdate(); err != nil {
		return fmt.Errorf("insert default admin failed: %w", err)
	}

	logger.Info().
		Str("email", setting.AppSetting.AdminEmail).
		Msg("default admin account created")
	return nil
}

// Close closes the database connection
func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Error getting SQL DB handle: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}
}

func BeginTx() *gorm.DB {
	return DB.Begin()
}

func CommitTx(tx *gorm.DB) error {
	return tx.Commit().Error
}
