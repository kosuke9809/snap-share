package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string

	R2AccountID       string
	R2AccessKey       string
	R2SecretAccessKey string
	R2BucketName      string
	R2PublicDomain    string
}

func Load() (*Config, error) {
	config := &Config{
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),

		R2AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKey:       os.Getenv("R2_ACCESS_KEY"),
		R2SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		R2BucketName:      os.Getenv("R2_BUCKET_NAME"),
		R2PublicDomain:    os.Getenv("R2_PUBLIC_DOMAIN"),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.R2AccountID == "" {
		return fmt.Errorf("R2_ACCOUNT_ID is required")
	}
	if c.R2AccessKey == "" {
		return fmt.Errorf("R2_ACCESS_KEY is required")
	}
	if c.R2SecretAccessKey == "" {
		return fmt.Errorf("R2_SECRET_ACCESS_KEY is required")
	}
	if c.R2BucketName == "" {
		return fmt.Errorf("R2_BUCKET_NAME is required")
	}
	if c.R2PublicDomain == "" {
		return fmt.Errorf("R2_PUBLIC_DOMAIN is required")
	}

	// Set default port if not provided
	if c.Port == "" {
		c.Port = "8080"
	}

	return nil
}
