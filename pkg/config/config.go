package config

// VerboseLoggingKey is the env var that determines migrator's logging level
const VerboseLoggingKey = "VERBOSE_LOGGING"

type Config struct {
	PostgresDB     PostgresDBConfig
	VerboseLogging bool
}

func LoadConfig(defaultSettings DefaultSettings) (Config, error) {
	isVerbose, err := getEnvBoolOrDefault(VerboseLoggingKey, defaultSettings.getWithFallback(VerboseLoggingKey, "false"))
	if err != nil {
		return Config{}, err
	}
	postgresDBConfig, err := LoadPostgresDBConfig(defaultSettings)
	if err != nil {
		return Config{}, err
	}
	return Config{
		PostgresDB:     postgresDBConfig,
		VerboseLogging: isVerbose,
	}, nil
}
