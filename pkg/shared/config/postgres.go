package config

type PostgresDBConfig struct {
	Host     string
	Port     int
	User     string
	Password *string
	Database string
	Schema   string
}

func LoadPostgresDBConfig(defaultSettings DefaultSettings) PostgresDBConfig {
	return NewPostgresDBConfigBuilder(defaultSettings).Build()
}

type PostgresDBConfigBuilder struct {
	d DefaultSettings
	c *PostgresDBConfig
}

func NewPostgresDBConfigBuilder(defaultSettings DefaultSettings) *PostgresDBConfigBuilder {
	var present bool
	if _, present = defaultSettings["POSTGRES_HOST"]; !present {
		defaultSettings["POSTGRES_HOST"] = "localhost"
	}
	if _, present = defaultSettings["POSTGRES_PORT"]; !present {
		defaultSettings["POSTGRES_PORT"] = "5432"
	}
	if _, present = defaultSettings["POSTGRES_USER"]; !present {
		defaultSettings["POSTGRES_USER"] = ""
	}
	if _, present = defaultSettings["POSTGRES_PASSWORD"]; !present {
		defaultSettings["POSTGRES_PASSWORD"] = ""
	}
	if _, present = defaultSettings["POSTGRES_DATABASE"]; !present {
		defaultSettings["POSTGRES_DATABASE"] = "postgres"
	}
	if _, present = defaultSettings["POSTGRES_SCHEMA"]; !present {
		defaultSettings["POSTGRES_SCHEMA"] = ""
	}

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

func (b *PostgresDBConfigBuilder) Build() PostgresDBConfig {
	if len(b.c.Host) == 0 {
		b.c.Host = GetEnvOrDefault("POSTGRES_HOST", b.d["POSTGRES_HOST"])
	}
	if b.c.Port == 0 {
		b.c.Port = Atoi(GetEnvOrDefault("POSTGRES_PORT", b.d["POSTGRES_PORT"]))
	}
	if len(b.c.User) == 0 {
		b.c.User = GetEnvOrDefault("POSTGRES_USER", b.d["POSTGRES_USER"])
	}
	if b.c.Password == nil {
		password := GetEnvOrDefault("POSTGRES_PASSWORD", b.d["POSTGRES_PASSWORD"])
		if password != "" {
			b.c.Password = &password
		} else {
			b.c.Password = nil
		}
	}
	if len(b.c.Database) == 0 {
		b.c.Database = GetEnvOrDefault("POSTGRES_DATABASE", b.d["POSTGRES_DATABASE"])
	}
	if len(b.c.Schema) == 0 {
		b.c.Schema = GetEnvOrDefault("POSTGRES_SCHEMA", b.d["POSTGRES_SCHEMA"])
	}
	return *b.c
}
