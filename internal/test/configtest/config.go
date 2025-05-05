package configtest

import (
	config2 "github.com/pennsieve/dbmigrate-go/pkg/shared/config"
)

// PostgresDBConfig returns a config.PostgresDBConfig suitable for use against
// the pennseivedb instance started for testing. It is preferred in tests over
// calling config.LoadPostgresDBConfig() because that method
// will not create the correct configs if the tests are running locally instead
// of in the Docker test container.
func PostgresDBConfig(options ...PostgresOption) config2.PostgresDBConfig {
	defaultSettings := config2.NewDefaultSettings()
	builder := config2.NewPostgresDBConfigBuilder(defaultSettings).
		WithPostgresUser("postgres").
		WithPostgresPassword("password")
	for _, option := range options {
		builder = option(builder)
	}
	return builder.Build()
}

type PostgresOption func(builder *config2.PostgresDBConfigBuilder) *config2.PostgresDBConfigBuilder

func WithPort(port int) PostgresOption {
	return func(builder *config2.PostgresDBConfigBuilder) *config2.PostgresDBConfigBuilder {
		return builder.WithPort(port)
	}
}

func WithHost(host string) PostgresOption {
	return func(builder *config2.PostgresDBConfigBuilder) *config2.PostgresDBConfigBuilder {
		return builder.WithHost(host)
	}
}

func WithSchema(schema string) PostgresOption {
	return func(builder *config2.PostgresDBConfigBuilder) *config2.PostgresDBConfigBuilder {
		return builder.WithSchema(schema)
	}
}
