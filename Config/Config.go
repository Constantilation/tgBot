package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

// GetConfig возвращает значение для заданного ключа конфигурации
func GetConfig(key ConfigKey) string {
	// Используйте os.Getenv для получения значения из переменных среды
	value := os.Getenv(string(key))
	if value == "" {
		log.Printf("Warning: No value found for key %s\n", key)
	}
	return value
}

// GetInt64Config пытается преобразовать значение заданного ключа конфигурации в int64
func GetInt64Config(key ConfigKey) (int64, error) {
	valueStr := GetConfig(key)
	if valueStr == "" {
		return 0, nil // или возвращайте ошибку, если отсутствие значения является критическим
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing %s as int64: %v\n", key, err)
		return 0, err
	}
	return value, nil
}
