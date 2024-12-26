package db

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable"+
		" TimeZone=Asia/Kolkata", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	fmt.Println(dsn)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db := lo.Must(gorm.Open(postgres.Open(dsn), config))
	return db
}
