package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Load .env file configurations into PROC ENVIRONMENT
func Load(dotEnvConfigFilePath string) error {
	dotEnvConfigFilePath, err := filepath.Abs(dotEnvConfigFilePath)
	if err != nil {
		return err
	}

	if err := godotenv.Load(dotEnvConfigFilePath); err != nil {
		return err
	}
	return nil
}

func lookup(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultVal
}

func GetString(key string, defaultVal string) string {
	return lookup(key, defaultVal)
}

func GetInt(key string, defaultVal int) int {
	value := lookup(key, "")
	if value, err := strconv.Atoi(value); err == nil {
		return value
	}
	return defaultVal
}

func GetBool(key string, defaultVal bool) bool {
	value := lookup(key, "")
	if value, err := strconv.ParseBool(value); err == nil {
		return value
	}
	return defaultVal
}

func GetDuration(key string, defaultVal time.Duration) time.Duration {
	value := lookup(key, "")
	if value, err := time.ParseDuration(value); err == nil {
		return value
	}
	return defaultVal
}
