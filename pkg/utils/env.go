package utils

import (
	"os"
	"strconv"
)

// GetEnv return the value of key environment variable or fallback if the
// variable is not set
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// IsTLSEnabled returns if the key environment variable allows TLS
func IsTLSEnabled(key string) bool {
	enabled, err := strconv.ParseBool(GetEnv(key, "false"))
	return err == nil && enabled
}
