package config_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/dbmigrate-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestLoadConfig_EmptyDefaultSettings(t *testing.T) {
	// Need to unset vars for CI
	unsetConfigEnvVars(t)

	settings := config.NewDefaultSettings()
	emptyConfig, err := config.LoadConfig(settings)
	require.NoError(t, err)
	assert.Equal(t, config.Config{
		PostgresDB: config.PostgresDBConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "",
			Password: nil,
			Database: "postgres",
			Schema:   "",
		},
		VerboseLogging: false,
	}, emptyConfig)
}

func TestLoadConfig_EnvVars(t *testing.T) {
	expectedPassword := uuid.NewString()
	expected := config.Config{
		PostgresDB: config.PostgresDBConfig{
			Host:     uuid.NewString(),
			Port:     rand.Intn(6000) + 1,
			User:     uuid.NewString(),
			Password: &expectedPassword,
			Database: uuid.NewString(),
			Schema:   uuid.NewString(),
		},
		VerboseLogging: true,
	}

	t.Setenv(config.VerboseLoggingKey, strconv.FormatBool(expected.VerboseLogging))
	t.Setenv(config.PostgresHostKey, expected.PostgresDB.Host)
	t.Setenv(config.PostgresPortKey, fmt.Sprintf("%d", expected.PostgresDB.Port))
	t.Setenv(config.PostgresUserKey, expected.PostgresDB.User)
	t.Setenv(config.PostgresPasswordKey, *expected.PostgresDB.Password)
	t.Setenv(config.PostgresDatabaseKey, expected.PostgresDB.Database)
	t.Setenv(config.PostgresSchemaKey, expected.PostgresDB.Schema)

	envConfig, err := config.LoadConfig(config.NewDefaultSettings())
	require.NoError(t, err)
	assert.Equal(t, expected, envConfig)
}

func TestLoadConfig_DefaultSettings(t *testing.T) {
	expectedPassword := uuid.NewString()
	expected := config.Config{
		PostgresDB: config.PostgresDBConfig{
			Host:     uuid.NewString(),
			Port:     rand.Intn(6000) + 1,
			User:     uuid.NewString(),
			Password: &expectedPassword,
			Database: uuid.NewString(),
			Schema:   uuid.NewString(),
		},
		VerboseLogging: true,
	}
	settings := config.NewDefaultSettings()
	settings[config.VerboseLoggingKey] = strconv.FormatBool(expected.VerboseLogging)
	settings[config.PostgresHostKey] = expected.PostgresDB.Host
	settings[config.PostgresPortKey] = fmt.Sprintf("%d", expected.PostgresDB.Port)
	settings[config.PostgresUserKey] = expected.PostgresDB.User
	settings[config.PostgresPasswordKey] = *expected.PostgresDB.Password
	settings[config.PostgresDatabaseKey] = expected.PostgresDB.Database
	settings[config.PostgresSchemaKey] = expected.PostgresDB.Schema

	// Need to unset vars for CI
	unsetConfigEnvVars(t)

	settingsConfig, err := config.LoadConfig(settings)
	require.NoError(t, err)
	assert.Equal(t, expected, settingsConfig)
}

// Unsetenv unsets the environment variable 'key' in the scope of 't' and will
// reset it to its previous value (if any) when the test scoped by 't' completes.
func unsetenv(t *testing.T, key string) {
	t.Helper()
	// Setenv takes care of registering the value-reset with its Cleanup method.
	// But setting the value to "" is not the same as unsetting it, so we follow
	// with an Unsetenv.
	t.Setenv(key, "")
	require.NoError(t, os.Unsetenv(key))
}

func unsetConfigEnvVars(t *testing.T) {
	t.Helper()
	unsetenv(t, config.VerboseLoggingKey)
	unsetenv(t, config.PostgresHostKey)
	unsetenv(t, config.PostgresPortKey)
	unsetenv(t, config.PostgresUserKey)
	unsetenv(t, config.PostgresPasswordKey)
	unsetenv(t, config.PostgresDatabaseKey)
	unsetenv(t, config.PostgresSchemaKey)
}
