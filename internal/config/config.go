package config

import (
	"bufio"
	"os"
	"strings"
)

var envMap map[string]string

func LoadEnv() error {
	envMap = make(map[string]string)
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		envMap[key] = value
	}

	return scanner.Err()
}

func GetEnvValue(key string) string {
	if val, ok := envMap[key]; ok {
		return val
	}
	return ""
}
