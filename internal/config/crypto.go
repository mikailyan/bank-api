package config

import (
	"github.com/sirupsen/logrus"
)

type CryptoConfig struct {
	PGPKey  string
	HMACKey string
}

func LoadCrypto() CryptoConfig {
	cfg := CryptoConfig{
		PGPKey:  getEnv("BANK_PGP_KEY", "bankDefaultPGPKey2024"),
		HMACKey: getEnv("BANK_HMAC_KEY", "bankDefaultHMACKey2024"),
	}

	logrus.Info("Конфигурация криптографических ключей загружена")

	return cfg
}
