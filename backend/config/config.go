package config

import (
	"os"
	"sync"
)

type Config struct {
	JWTSecret string
}

var (
	instance *Config
	once     sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		// Получаем секретный ключ из переменных окружения
		// Если переменная не установлена, используем значение по умолчанию
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "ewstrdygihuoio[p'i;ulykewfrawsgrfhdgtuky" // Значение по умолчанию
		}

		instance = &Config{
			JWTSecret: jwtSecret,
		}
	})

	return instance
}
