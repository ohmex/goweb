package db

import (
	"fmt"
	"goweb/config"
	"goweb/db/seeders"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init(cfg *config.Config) *gorm.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name)

	fmt.Println(dataSourceName)

	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	userSeeder := seeders.NewUserSeeder(db)
	userSeeder.SetUsers()

	return db
}

func InitRedis(cfg *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
