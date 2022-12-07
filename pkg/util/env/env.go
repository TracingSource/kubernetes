package env

import (
	"os"
	"strconv"
)

// GetEnvAsStringOrFallback returns the env variable for the given key
// and falls back to the given defaultValue if not set
func GetEnvAsStringOrFallback(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// GetEnvAsIntOrFallback returns the env variable (parsed as integer) for
// the given key and falls back to the given defaultValue if not set
func GetEnvAsIntOrFallback(key string, defaultValue int) (int, error) {
	if v := os.Getenv(key); v != "" {
		value, err := strconv.Atoi(v)
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}
	return defaultValue, nil
}

// GetEnvAsFloat64OrFallback returns the env variable (parsed as float64) for
// the given key and falls back to the given defaultValue if not set
func GetEnvAsFloat64OrFallback(key string, defaultValue float64) (float64, error) {
	if v := os.Getenv(key); v != "" {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}
	return defaultValue, nil
}
