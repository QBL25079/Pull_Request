package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	HTTPPort string
	DBUser   string
	DBPass   string
	DBHost   string
	DBName   string
	DBPort   int
}

func Load() (*Config, error) {
	httpPort := os.Getenv("HTTP_PORT")

	if httpPort == "" {
		httpPort = "8080"
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dbPortStr := os.Getenv("DB_PORT")
	portDb, err := strconv.Atoi(dbPortStr)
	if err != nil {
		portDb = 5432
	}

	return &Config{
		HTTPPort: httpPort,
		DBUser:   dbUser,
		DBPass:   dbPass,
		DBHost:   dbHost,
		DBName:   dbName,
		DBPort:   portDb,
	}, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName)
}
