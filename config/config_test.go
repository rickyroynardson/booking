package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	config = nil

	t.Setenv("APP_PORT", "8080")
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_USER", "postgres")
	t.Setenv("POSTGRES_PASSWORD", "postgres")
	t.Setenv("POSTGRES_DB", "test_db")

	cfg := Get()

	assert.Equal(t, "8080", cfg.App.Port)
	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, "5432", cfg.DB.Port)
	assert.Equal(t, "postgres", cfg.DB.User)
	assert.Equal(t, "postgres", cfg.DB.Password)
	assert.Equal(t, "test_db", cfg.DB.DBName)

	cfg2 := Get()
	assert.Same(t, cfg, cfg2)
}
