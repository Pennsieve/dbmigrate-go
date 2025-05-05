package config

import (
	"fmt"
	"strconv"
)

type DefaultSettings map[string]string

func NewDefaultSettings() DefaultSettings {
	return make(DefaultSettings, 20)
}

const VerboseLoggingKey = "VERBOSE_LOGGING"

type Config struct {
	PostgresDB     PostgresDBConfig
	VerboseLogging bool
}

func LoadConfig(defaultSettings DefaultSettings) (Config, error) {
	verboseStr := GetEnvOrDefault(VerboseLoggingKey, "false")
	isVerbose, err := strconv.ParseBool(verboseStr)
	if err != nil {
		return Config{}, fmt.Errorf("error converting %q value %s to bool: %w",
			VerboseLoggingKey, verboseStr, err)
	}
	return Config{
		PostgresDB:     LoadPostgresDBConfig(defaultSettings),
		VerboseLogging: isVerbose,
	}, nil
}
