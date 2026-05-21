package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func getEnvAsDuration(key string, defaultValue interface{}) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		switch v := defaultValue.(type) {
		case string:
			if duration, err := time.ParseDuration(v); err == nil {
				return duration
			}
		case time.Duration:
			return v
		}
		return 0
	}

	if duration, err := time.ParseDuration(valueStr); err == nil {
		return duration
	}
	return 0
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAny(keys []string, defaultValue string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsIntAny(keys []string, defaultValue int) int {
	for _, key := range keys {
		valueStr := strings.TrimSpace(os.Getenv(key))
		if valueStr == "" {
			continue
		}
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		rawItems := strings.Split(value, ",")
		items := make([]string, 0, len(rawItems))
		for _, item := range rawItems {
			trimmed := strings.TrimSpace(item)
			if trimmed != "" {
				items = append(items, trimmed)
			}
		}
		if len(items) > 0 {
			return items
		}
	}
	return defaultValue
}
