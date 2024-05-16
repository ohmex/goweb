package db

import (
	"fmt"
	"goweb/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name)

	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	return db
}

func InitRedis(cfg *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
