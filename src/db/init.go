package db

import (
	"github.com/samber/lo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=equisplit port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db := lo.Must(gorm.Open(postgres.Open(dsn), config))
	return db
}
