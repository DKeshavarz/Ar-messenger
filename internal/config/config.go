package config

import (
	"os"

	"github.com/joho/godotenv"
)

var envMap map[string]string

func LoadEnv() error {
	return godotenv.Load(".env")
}

func GetEnvValue(key string) string {
	return os.Getenv(key)
}
