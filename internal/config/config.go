package config

import (
	"os"
)

func GetEnvValue(key string) string {
	return os.Getenv(key)
}
