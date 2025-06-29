package config

import (
	"os"
	"time"
)

func GetJwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

func GetJwtExpirationDuration() time.Duration {
	duration, err := time.ParseDuration(os.Getenv("JWT_EXPIRES_IN"))

	if err != nil {
		return time.Hour * 24
	}

	return duration
}
