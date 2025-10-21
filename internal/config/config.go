package config

import (
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds application configuration
type Config struct {
	Server struct {
		Host string
		Port string
	}
	Database struct {
		Driver   string
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{}

	// Server config
	cfg.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	// Database config
	cfg.Database.Driver = getEnv("DB_DRIVER", "postgres")
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvInt("DB_PORT", 5432)
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.DBName = getEnv("DB_NAME", "myapp")

	return cfg
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *Config) (*gorm.DB, error) {
	switch cfg.Database.Driver {
	case "mysql":
		return connectMySQL(cfg)
	case "postgres":
		return connectPostgres(cfg)
	default:
		return connectPostgres(cfg)
	}
}

func connectPostgres(cfg *Config) (*gorm.DB, error) {
	dsn := "host=" + cfg.Database.Host +
		" port=" + strconv.Itoa(cfg.Database.Port) +
		" user=" + cfg.Database.User +
		" password=" + cfg.Database.Password +
		" dbname=" + cfg.Database.DBName +
		" sslmode=disable"

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func connectMySQL(cfg *Config) (*gorm.DB, error) {
	dsn := cfg.Database.User + ":" + cfg.Database.Password +
		"@tcp(" + cfg.Database.Host + ":" + strconv.Itoa(cfg.Database.Port) + ")/" +
		cfg.Database.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}