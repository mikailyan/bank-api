package config

import (
	"os"
	"time"
)

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

func LoadJWT() JWTConfig {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-bank-api-jwt-secret-key"
	}

	return JWTConfig{
		Secret:    secret,
		ExpiresIn: 24 * time.Hour,
	}
}
