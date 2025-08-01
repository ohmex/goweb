package config

import (
	"goweb/util"
	"os"
)

type DBConfig struct {
	User                string
	Password            string
	Driver              string
	Name                string
	Host                string
	Port                string
	PartitioningEnabled bool
}

func LoadDBConfig() DBConfig {
	partitioning := util.IsPartitioningEnabled()
	return DBConfig{
		User:                os.Getenv("DB_USER"),
		Password:            os.Getenv("DB_PASSWORD"),
		Driver:              os.Getenv("DB_DRIVER"),
		Name:                os.Getenv("DB_NAME"),
		Host:                os.Getenv("DB_HOST"),
		Port:                os.Getenv("DB_PORT"),
		PartitioningEnabled: partitioning,
	}
}
