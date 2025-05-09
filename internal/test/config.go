package test

import (
	"github.com/pennsieve/dbmigrate-go/pkg/config"
)

// NewTestSettings returns a config.DefaultSettings for running tests against
// the DB started by docker-compose.test.yml.
// When tests are run CI in Docker these settings will be handled by env vars
// set in the Docker container. But this function returns settings that will
// allow tests to also run locally outside of Docker where the env vars are not
// set.
func NewTestSettings(testSchema string) config.DefaultSettings {
	return config.DefaultSettings{
		config.VerboseLoggingKey:   "true",
		config.PostgresHostKey:     "localhost",
		config.PostgresPortKey:     "5432",
		config.PostgresUserKey:     "postgres",
		config.PostgresPasswordKey: "password",
		config.PostgresDatabaseKey: "postgres",
		config.PostgresSchemaKey:   testSchema,
	}
}
