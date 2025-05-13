package config

// PostgresHostKey is the env var for the host the migrator will connect to
const PostgresHostKey = "POSTGRES_HOST"

// PostgresPortKey is the env var for the port the migrator will connect to
const PostgresPortKey = "POSTGRES_PORT"

// PostgresUserKey is the env var for the user used to connect with
const PostgresUserKey = "POSTGRES_USER"

// PostgresPasswordKey is the env var for the password used to connect with.
// If this is not set or is set to "", the migrator will attempt to connect
// with an RDS auth token.
const PostgresPasswordKey = "POSTGRES_PASSWORD"

// PostgresDatabaseKey is the env var for the database name the migrator will run in
const PostgresDatabaseKey = "POSTGRES_DATABASE"

// PostgresSchemaKey is the env var for the schema name the migrator will create if necessary and run in
const PostgresSchemaKey = "POSTGRES_SCHEMA"

type PostgresDBConfig struct {
	Host     string
	Port     int
	User     string
	Password *string
	Database string
	Schema   string
}

func LoadPostgresDBConfig(defaultSettings DefaultSettings) (PostgresDBConfig, error) {
	return NewPostgresDBConfigBuilder(defaultSettings).Build()
}

type PostgresDBConfigBuilder struct {
	d DefaultSettings
	c *PostgresDBConfig
}

func NewPostgresDBConfigBuilder(defaultSettings DefaultSettings) *PostgresDBConfigBuilder {
	return &PostgresDBConfigBuilder{
		d: defaultSettings,
		c: &PostgresDBConfig{},
	}
}

func (b *PostgresDBConfigBuilder) WithPostgresUser(postgresUser string) *PostgresDBConfigBuilder {
	b.c.User = postgresUser
	return b
}

func (b *PostgresDBConfigBuilder) WithPostgresPassword(postgresPassword string) *PostgresDBConfigBuilder {
	b.c.Password = &postgresPassword
	return b
}

func (b *PostgresDBConfigBuilder) WithHost(host string) *PostgresDBConfigBuilder {
	b.c.Host = host
	return b
}

func (b *PostgresDBConfigBuilder) WithPort(port int) *PostgresDBConfigBuilder {
	b.c.Port = port
	return b
}

func (b *PostgresDBConfigBuilder) WithSchema(schema string) *PostgresDBConfigBuilder {
	b.c.Schema = schema
	return b
}

func (b *PostgresDBConfigBuilder) Build() (PostgresDBConfig, error) {
	if len(b.c.Host) == 0 {
		b.c.Host = getEnvOrDefault(PostgresHostKey, b.d.getWithFallback(PostgresHostKey, "localhost"))
	}
	if b.c.Port == 0 {
		port, err := getEnvIntOrDefault(PostgresPortKey, b.d.getWithFallback(PostgresPortKey, "5432"))
		if err != nil {
			return PostgresDBConfig{}, err
		}
		b.c.Port = port
	}
	if len(b.c.User) == 0 {
		b.c.User = getEnvOrDefault(PostgresUserKey, b.d.get(PostgresUserKey))
	}
	if b.c.Password == nil {
		password := getEnvOrDefault(PostgresPasswordKey, b.d.get(PostgresPasswordKey))
		if password != "" {
			b.c.Password = &password
		} else {
			b.c.Password = nil
		}
	}
	if len(b.c.Database) == 0 {
		b.c.Database = getEnvOrDefault(PostgresDatabaseKey, b.d.getWithFallback(PostgresDatabaseKey, "postgres"))
	}
	if len(b.c.Schema) == 0 {
		b.c.Schema = getEnvOrDefault(PostgresSchemaKey, b.d.get(PostgresSchemaKey))
	}
	return *b.c, nil
}
