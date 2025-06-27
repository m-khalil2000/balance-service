package config

import (
	"balance-service/pkg/models"
	"fmt"
	"os"
	"time"
)

// Loads config from docker-compose environment variables
func LoadFromEnv() *models.Config {

	maxOpen := parseEnvInt("DB_MAX_OPEN_CONNS", 10)
	maxIdle := parseEnvInt("DB_MAX_IDLE_CONNS", 5)
	lifetime := parseEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute)
	idleTime := parseEnvDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute)
	return &models.Config{
		DB: models.DBConfig{
			Host:            os.Getenv("DB_HOST"),
			Port:            os.Getenv("DB_PORT"),
			User:            readEnvOrFile("DB_USER"),
			Password:        readEnvOrFile("DB_PASSWORD"),
			Name:            os.Getenv("DB_NAME"),
			SSLMode:         os.Getenv("DB_SSLMODE"),
			MaxOpenConns:    maxOpen,
			MaxIdleConns:    maxIdle,
			ConnMaxLifetime: lifetime,
			ConnMaxIdleTime: idleTime,
		},
		Server: models.ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
	}
}

// Constructs the PostgreSQL connection string from the DBConfig.
func DSN(c models.DBConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func GetPort(c *models.ServerConfig) string {
	if c.Port == "" {
		return "8080"
	}
	return c.Port
}
