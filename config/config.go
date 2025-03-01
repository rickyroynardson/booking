package config

import (
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
	c := initConfig()
	config = c
	return c
}

func initConfig() *Config {
	c := Config{}
	c.App.Port = os.Getenv("APP_PORT")
	if c.App.Port == "" {
		log.Fatal("APP_PORT is required")
	}

	c.DB.Host = os.Getenv("POSTGRES_HOST")
	c.DB.Port = os.Getenv("POSTGRES_PORT")
	c.DB.User = os.Getenv("POSTGRES_USER")
	c.DB.Password = os.Getenv("POSTGRES_PASSWORD")
	c.DB.DBName = os.Getenv("POSTGRES_DB")
	if c.DB.Host == "" {
		log.Fatal("POSTGRES_HOST is required")
	}
	if c.DB.Port == "" {
		log.Fatal("POSTGRES_PORT is required")
	}
	if c.DB.User == "" {
		log.Fatal("POSTGRES_USER is required")
	}
	if c.DB.Password == "" {
		log.Fatal("POSTGRES_PASSWORD is required")
	}
	if c.DB.DBName == "" {
		log.Fatal("POSTGRES_DB is required")
	}
	return &c
}
