package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Invalid value for %s: %s (defaulting to %d)", key, valStr, defaultVal)
		return defaultVal
	}
	return val
}

func parseEnvDuration(key string, defaultVal time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := time.ParseDuration(valStr)
	if err != nil {
		log.Printf("Invalid duration for %s: %s (defaulting to %s)", key, valStr, defaultVal)
		return defaultVal
	}
	return val
}

func readEnvOrFile(key string) string {
	// check plain variables
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	// check for docker secrets
	filePath := os.Getenv(key + "_FILE")
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading secret file %s: %v", filePath, err)
			return ""
		}
		return strings.TrimSpace(string(data))
	}
	log.Printf("Environment variable %s or %s_FILE not set", key, key)
	return ""
}
