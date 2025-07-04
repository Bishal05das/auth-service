package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port         string
		Host         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}

	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}

	JWT struct {
		Secretkey string
		TokenExpiry time.Duration
		RefreshExpiry time.Duration
	}

	Environment string
}

func Load() (*Config, error) {
	godotenv.Load()

	cfg := &Config{}

	//server config
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	cfg.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	cfg.Server.ReadTimeout = time.Second * 15
	cfg.Server.WriteTimeout = time.Second * 15
	cfg.Environment = getEnv("ENV", "development")

	//databae config
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5434")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "password")
	cfg.Database.DBName = getEnv("DB_NAME", "auth_service")
	cfg.Database.SSLMode = getEnv("DB_SSLmode", "disable")

	//JWT config
	cfg.JWT.Secretkey= getEnv("JWT_SECRET", "mysecretkey")
	cfg.JWT.TokenExpiry = time.Hour * 24  // 24 hours
	cfg.JWT.RefreshExpiry = time.Hour * 168 // 7 days
	return cfg, nil

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
