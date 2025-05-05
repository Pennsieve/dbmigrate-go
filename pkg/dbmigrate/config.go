package dbmigrate

import (
	"fmt"
	config2 "github.com/pennsieve/dbmigrate-go/pkg/shared/config"
	"strconv"
)

const VerboseLoggingKey = "VERBOSE_LOGGING"

type Config struct {
	PostgresDB     config2.PostgresDBConfig
	VerboseLogging bool
}

func LoadConfig(defaultSettings config2.DefaultSettings) (Config, error) {
	verboseStr := config2.GetEnvOrDefault(VerboseLoggingKey, "false")
	isVerbose, err := strconv.ParseBool(verboseStr)
	if err != nil {
		return Config{}, fmt.Errorf("error converting %q value %s to bool: %w",
			VerboseLoggingKey, verboseStr, err)
	}
	return Config{
		PostgresDB:     config2.LoadPostgresDBConfig(defaultSettings),
		VerboseLogging: isVerbose,
	}, nil
}
