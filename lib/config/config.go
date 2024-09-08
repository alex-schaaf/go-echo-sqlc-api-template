package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_PATH    string
	PORT       int
	JWT_SECRET string
}

// InitConfig initializes the Config struct with the environment variables
// and sets default parameters
func InitConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	config := &Config{
		DB_PATH:    os.Getenv("DB_PATH"),
		PORT:       port,
		JWT_SECRET: os.Getenv("JWT_SECRET"),
	}

	return config
}
