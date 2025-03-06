package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	App struct {
		Port string
	}
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
}

var config *Config

func Get() *Config {
	if config != nil {
		return config
	}
	c, err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	config = c
	return c
}

func validateConfig() error {
	if os.Getenv("APP_PORT") == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	if os.Getenv("POSTGRES_HOST") == "" {
		return fmt.Errorf("POSTGRES_HOST is required")
	}
	if os.Getenv("POSTGRES_PORT") == "" {
		return fmt.Errorf("POSTGRES_PORT is required")
	}
	if os.Getenv("POSTGRES_USER") == "" {
		return fmt.Errorf("POSTGRES_USER is required")
	}
	if os.Getenv("POSTGRES_PASSWORD") == "" {
		return fmt.Errorf("POSTGRES_PASSWORD is required")
	}
	if os.Getenv("POSTGRES_DB") == "" {
		return fmt.Errorf("POSTGRES_DB is required")
	}
	return nil
}

func initConfig() (*Config, error) {
	if err := validateConfig(); err != nil {
		return nil, err
	}

	c := Config{}
	c.App.Port = os.Getenv("APP_PORT")

	c.DB.Host = os.Getenv("POSTGRES_HOST")
	c.DB.Port = os.Getenv("POSTGRES_PORT")
	c.DB.User = os.Getenv("POSTGRES_USER")
	c.DB.Password = os.Getenv("POSTGRES_PASSWORD")
	c.DB.DBName = os.Getenv("POSTGRES_DB")
	return &c, nil
}
