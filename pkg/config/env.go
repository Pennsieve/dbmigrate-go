package config

import (
	"fmt"
	"os"
	"strconv"
)

func getEnvOrDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else {
		return defaultValue
	}
}

func getEnvIntOrDefault(key string, defaultValue string) (int, error) {
	strValue := getEnvOrDefault(key, defaultValue)
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, fmt.Errorf("error converting '%s' value '%s' to int: %w",
			key, strValue, err)
	}
	return value, nil
}

func getEnvBoolOrDefault(key string, defaultValue string) (bool, error) {
	strValue := getEnvOrDefault(key, defaultValue)
	value, err := strconv.ParseBool(strValue)
	if err != nil {
		return false, fmt.Errorf("error converting '%s' value '%s' to bool: %w",
			key, strValue, err)
	}
	return value, nil
}
