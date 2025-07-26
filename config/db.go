package config

import "os"

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
	partitioning := false
	if os.Getenv("DB_PARTITIONING_ENABLED") == "true" {
		partitioning = true
	}
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
