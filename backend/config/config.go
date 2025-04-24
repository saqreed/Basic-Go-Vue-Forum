package config

import (
	"os"
)

type Config struct {
	JWTSecret string
}

func LoadConfig() *Config {
	// Получаем секретный ключ из переменных окружения
	// Если переменная не установлена, используем значение по умолчанию
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "ewstrdygihuoio[p'i;ulykewfrawsgrfhdgtuky" // Значение по умолчанию
	}

	return &Config{
		JWTSecret: jwtSecret,
	}
}
