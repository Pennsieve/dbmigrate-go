package config_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/dbmigrate-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
)

func TestLoadConfig_EmptyDefaultSettings(t *testing.T) {
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

	settingsConfig, err := config.LoadConfig(settings)
	require.NoError(t, err)
	assert.Equal(t, expected, settingsConfig)
}
