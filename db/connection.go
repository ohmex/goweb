package db

import (
	"fmt"
	"goweb/config"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) *gorm.DB {
	var db *gorm.DB
	var err error
	var dataSourceName string

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level: Silent, Error, Warn, Info
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	switch cfg.DB.Driver {
	case "mysql":
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.Name)
		db, err = gorm.Open(mysql.Open(dataSourceName), &gorm.Config{
			Logger: newLogger,
			// Performance optimizations
			SkipDefaultTransaction: true, // Disable default transaction for better performance
			PrepareStmt:            true, // Enable prepared statements
		})
	case "postgres":
		dataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Name)
		db, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
			Logger: newLogger,
			// Performance optimizations
			SkipDefaultTransaction: true, // Disable default transaction for better performance
			PrepareStmt:            true, // Enable prepared statements
		})
	case "yugabytedb":
		// YugabyteDB is PostgreSQL-compatible, so we can use the PostgreSQL driver
		dataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Name)
		db, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
			Logger: newLogger,
			// Performance optimizations
			SkipDefaultTransaction: true, // Disable default transaction for better performance
			PrepareStmt:            true, // Enable prepared statements
		})
	default:
		panic("unsupported database driver: " + cfg.DB.Driver)
	}

	if err != nil {
		panic(err.Error())
	}

	// Configure connection pool for better performance
	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error())
	}

	// Set connection pool settings for optimal performance
	sqlDB.SetMaxIdleConns(25)                  // Increased idle connections for better reuse
	sqlDB.SetMaxOpenConns(200)                 // Increased max connections for high concurrency
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Shorter lifetime for better connection freshness
	sqlDB.SetConnMaxIdleTime(15 * time.Minute) // Shorter idle time for better resource management

	return db
}

func InitRedis(cfg *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		// Optimized connection pool settings for better performance
		PoolSize:     50, // Increased pool size for high concurrency
		MinIdleConns: 10, // Increased minimum idle connections
		MaxRetries:   3,  // Maximum number of retries
		// Optimized timeout settings
		DialTimeout:  3 * time.Second, // Reduced dial timeout
		ReadTimeout:  2 * time.Second, // Reduced read timeout
		WriteTimeout: 2 * time.Second, // Reduced write timeout
		// Additional performance settings
		PoolTimeout: 4 * time.Second, // Timeout for getting connection from pool
	})
}
