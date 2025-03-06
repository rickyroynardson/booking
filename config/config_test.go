package config

import (
	"os"
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

func TestValidateConfig(t *testing.T) {
	cases := []struct {
		name        string
		setupEnv    func()
		expectedErr string
	}{
		{
			name: "Missing APP_PORT",
			setupEnv: func() {
				os.Clearenv()
				t.Setenv("POSTGRES_HOST", "localhost")
				t.Setenv("POSTGRES_PORT", "5432")
				t.Setenv("POSTGRES_USER", "postgres")
				t.Setenv("POSTGRES_PASSWORD", "postgres")
				t.Setenv("POSTGRES_DB", "test_db")
			},
			expectedErr: "APP_PORT is required",
		},
		{
			name: "Missing POSTGRES_HOST",
			setupEnv: func() {
				os.Clearenv()
				t.Setenv("APP_PORT", "8080")
				t.Setenv("POSTGRES_PORT", "5432")
				t.Setenv("POSTGRES_USER", "postgres")
				t.Setenv("POSTGRES_PASSWORD", "postgres")
				t.Setenv("POSTGRES_DB", "test_db")
			},
			expectedErr: "POSTGRES_HOST is required",
		},
		{
			name: "All environment variables set",
			setupEnv: func() {
				os.Clearenv()
				t.Setenv("APP_PORT", "8080")
				t.Setenv("POSTGRES_HOST", "localhost")
				t.Setenv("POSTGRES_PORT", "5432")
				t.Setenv("POSTGRES_USER", "postgres")
				t.Setenv("POSTGRES_PASSWORD", "postgres")
				t.Setenv("POSTGRES_DB", "test_db")
			},
			expectedErr: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			config = nil
			c.setupEnv()
			err := validateConfig()
			if c.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, c.expectedErr)
			}
		})
	}
}

func TestInitConfig(t *testing.T) {
	os.Clearenv()
	t.Setenv("APP_PORT", "8080")
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_USER", "postgres")
	t.Setenv("POSTGRES_PASSWORD", "postgres")
	t.Setenv("POSTGRES_DB", "test_db")

	c, err := initConfig()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "8080", c.App.Port)
	assert.Equal(t, "localhost", c.DB.Host)
	assert.Equal(t, "5432", c.DB.Port)
	assert.Equal(t, "postgres", c.DB.User)
	assert.Equal(t, "postgres", c.DB.Password)
	assert.Equal(t, "test_db", c.DB.DBName)
}
