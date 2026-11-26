package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv               string
	AppPort              string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	DBSSLMode            string
	JWTSecret            string
	JWTExpiry            time.Duration
	AdminUsername        string
	AdminPassword        string
	AdminName            string
	SchoolTimezone       string
	ClockFutureTolerance time.Duration
	MaxEventAge          time.Duration
	SyncMaxBatch         int
}

func LoadConf() (*Config, error) {
	defaults := map[string]any{
		"APP_ENV":                "development",
		"APP_PORT":               "8080",
		"DB_HOST":                "localhost",
		"DB_PORT":                "5432",
		"DB_USER":                "",
		"DB_PASSWORD":            "",
		"DB_NAME":                "",
		"DB_SSLMODE":             "disable",
		"JWT_SECRET":             "",
		"JWT_EXPIRY":             "24h",
		"ADMIN_USERNAME":         "",
		"ADMIN_PASSWORD":         "",
		"ADMIN_NAME":             "",
		"SCHOOL_TIMEZONE":        "Asia/Jakarta",
		"CLOCK_FUTURE_TOLERANCE": "5m",
		"MAX_EVENT_AGE":          "720h",
		"SYNC_MAX_BATCH":         50,
	}
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFound) && !os.IsNotExist(err) {
			return nil, fmt.Errorf("read .env config: %w", err)
		}
	}

	for _, key := range []string{"DB_USER", "DB_NAME", "JWT_SECRET"} {
		if strings.TrimSpace(viper.GetString(key)) == "" {
			return nil, fmt.Errorf("%s is required", key)
		}
	}

	jwtExpiry, err := parseDurationWithDays(viper.GetString("JWT_EXPIRY"))
	if err != nil {
		return nil, fmt.Errorf("parse JWT_EXPIRY: %w", err)
	}
	clockFutureTolerance, err := parseDurationWithDays(viper.GetString("CLOCK_FUTURE_TOLERANCE"))
	if err != nil {
		return nil, fmt.Errorf("parse CLOCK_FUTURE_TOLERANCE: %w", err)
	}
	maxEventAge, err := parseDurationWithDays(viper.GetString("MAX_EVENT_AGE"))
	if err != nil {
		return nil, fmt.Errorf("parse MAX_EVENT_AGE: %w", err)
	}

	return &Config{
		AppEnv:               viper.GetString("APP_ENV"),
		AppPort:              viper.GetString("APP_PORT"),
		DBHost:               viper.GetString("DB_HOST"),
		DBPort:               viper.GetString("DB_PORT"),
		DBUser:               viper.GetString("DB_USER"),
		DBPassword:           viper.GetString("DB_PASSWORD"),
		DBName:               viper.GetString("DB_NAME"),
		DBSSLMode:            viper.GetString("DB_SSLMODE"),
		JWTSecret:            viper.GetString("JWT_SECRET"),
		JWTExpiry:            jwtExpiry,
		AdminUsername:        viper.GetString("ADMIN_USERNAME"),
		AdminPassword:        viper.GetString("ADMIN_PASSWORD"),
		AdminName:            viper.GetString("ADMIN_NAME"),
		SchoolTimezone:       viper.GetString("SCHOOL_TIMEZONE"),
		ClockFutureTolerance: clockFutureTolerance,
		MaxEventAge:          maxEventAge,
		SyncMaxBatch:         viper.GetInt("SYNC_MAX_BATCH"),
	}, nil
}

func parseDurationWithDays(raw string) (time.Duration, error) {
	raw = strings.TrimSpace(raw)
	if days, ok := strings.CutSuffix(raw, "d"); ok {
		dayCount, err := strconv.ParseInt(days, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid day duration %q: %w", raw, err)
		}
		return time.Duration(dayCount) * 24 * time.Hour, nil
	}

	duration, err := time.ParseDuration(raw)
	if err != nil {
		return 0, err
	}
	return duration, nil
}
