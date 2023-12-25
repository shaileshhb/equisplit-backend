package db

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
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

var Ctx = context.Background()

func InitCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_DB_ADDR"),
		Password: os.Getenv("REDIS_DB_PASS"),
		DB:       0,
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}

	return rdb
}
