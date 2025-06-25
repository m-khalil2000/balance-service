package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
}

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int           // Maximum number of open connections to the database
	MaxIdleConns    int           // Maximum number of idle connections to the database
	ConnMaxLifetime time.Duration // Maximum amount of time(s) a connection may be reused
	ConnMaxIdleTime time.Duration // Maximum amount of time(s) a connection may be idle
}

type ServerConfig struct {
	Port string
}

// Loads config from docker-compose environment variables
func LoadFromEnv() *Config {

	maxOpen, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if err != nil {
		log.Fatalf("Invalid DB_MAX_OPEN_CONNS: %v", err)
	}
	lifetime, err := time.ParseDuration(os.Getenv("DB_CONN_MAX_LIFETIME"))
	if err != nil {
		log.Fatalf("Invalid DB_CONN_MAX_LIFETIME: %v", err)
	}

	idleTime, err := time.ParseDuration(os.Getenv("DB_CONN_MAX_IDLE_TIME"))
	if err != nil {
		log.Fatalf("Invalid DB_CONN_MAX_IDLE_TIME: %v", err)
	}

	maxIdle, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	if err != nil {
		log.Fatalf("Invalid DB_MAX_IDLE_CONNS: %v", err)
	}
	return &Config{
		DB: DBConfig{
			Host:            os.Getenv("DB_HOST"),
			Port:            os.Getenv("DB_PORT"),
			User:            os.Getenv("DB_USER"),
			Password:        os.Getenv("DB_PASSWORD"),
			Name:            os.Getenv("DB_NAME"),
			SSLMode:         os.Getenv("DB_SSLMODE"),
			MaxOpenConns:    maxOpen,
			MaxIdleConns:    maxIdle,
			ConnMaxLifetime: lifetime,
			ConnMaxIdleTime: idleTime,
		},
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
	}
}

// Constructs the PostgreSQL connection string from the DBConfig.
func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func (c *ServerConfig) GetPort() string {
	if c.Port == "" {
		return "8080"
	}
	return c.Port
}
