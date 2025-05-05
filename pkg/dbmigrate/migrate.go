package dbmigrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"net"
	"net/url"
	"strings"
)

type DatabaseMigrator struct {
	wrapped *migrate.Migrate
}

func NewRDSProxyDatabaseMigrator(ctx context.Context, migrateConfig Config, migrationsSource source.Driver, awsConfig aws.Config) (*DatabaseMigrator, error) {
	authenticationToken, err := auth.BuildAuthToken(
		ctx,
		fmt.Sprintf("%s:%d", migrateConfig.PostgresDB.Host, migrateConfig.PostgresDB.Port),
		awsConfig.Region,
		migrateConfig.PostgresDB.User,
		awsConfig.Credentials,
	)
	if err != nil {
		return nil, fmt.Errorf("error building auth token for Migrator: %w", err)
	}
	return newDatabaseMigrator(
		ctx,
		migrateConfig.PostgresDB.User,
		authenticationToken,
		migrateConfig.PostgresDB.Host,
		migrateConfig.PostgresDB.Port,
		migrateConfig.PostgresDB.Database,
		migrateConfig.PostgresDB.Schema,
		migrationsSource,
		migrateConfig.VerboseLogging)
}

func NewLocalMigrator(ctx context.Context, migrateConfig Config, migrationsSource source.Driver) (*DatabaseMigrator, error) {
	if migrateConfig.PostgresDB.Password == nil {
		return nil, fmt.Errorf("password cannot be nil for local Migrator")
	}
	return newDatabaseMigrator(
		ctx,
		migrateConfig.PostgresDB.User,
		*migrateConfig.PostgresDB.Password,
		migrateConfig.PostgresDB.Host,
		migrateConfig.PostgresDB.Port,
		migrateConfig.PostgresDB.Database,
		migrateConfig.PostgresDB.Schema,
		migrationsSource,
		migrateConfig.VerboseLogging)

}

// Up will run any un-applied migrations
func (m *DatabaseMigrator) Up() error {
	if err := m.wrapped.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.wrapped.Log.Printf("no changes")
			return nil
		}
		return err
	}
	return nil
}

func (m *DatabaseMigrator) Down() error {
	if err := m.wrapped.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.wrapped.Log.Printf("no changes")
			return nil
		}
		return err
	}
	return nil
}

// Drop will drop all tables in the schema.
// Used for testing
func (m *DatabaseMigrator) Drop() error {
	return m.wrapped.Drop()
}

func (m *DatabaseMigrator) Close() (source error, database error) {
	return m.wrapped.Close()
}

func (m *DatabaseMigrator) CloseAndLogError() {
	sourceErr, databaseErr := m.Close()
	if sourceErr != nil {
		m.wrapped.Log.Printf("warning: source error closing DatabaseMigrator: %v", sourceErr)
	}
	if databaseErr != nil {
		m.wrapped.Log.Printf("warning: database error closing DatabaseMigrator: %v", databaseErr)

	}
}

func newDatabaseMigrator(ctx context.Context, username, password, host string,
	port int,
	databaseName string,
	schemaName string,
	migrationsSource source.Driver,
	verboseLogging bool) (*DatabaseMigrator, error) {

	// Migrate needs two things, a database.Driver to access Postgres, and a source.Driver to read the
	// migration files.

	// Create database.Driver and create schema (which Migrate won't do on its own)
	db, err := sql.Open("pgx",
		datasourceName(username,
			password,
			host,
			port,
			databaseName,
			schemaName),
	)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	// WithInstance will try to ensure that golang-migrate's migration table exists, so we need
	// to create the schema before it is called.
	createSchemaQuery := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %q", schemaName)
	if _, err := db.ExecContext(ctx, createSchemaQuery); err != nil {
		return nil, closeOnError(fmt.Errorf("error creating schema %q: %w", schemaName, err), db)
	}
	driver, err := pgx.WithInstance(db, &pgx.Config{SchemaName: schemaName})
	if err != nil {
		return nil, closeOnError(fmt.Errorf("error creating migration database.Driver: %w", err), db)
	}

	// Now we can create the Migrate instance
	m, err := migrate.NewWithInstance(
		"migration iofs",
		migrationsSource,
		"postgres",
		driver)
	if err != nil {
		return nil, closeOnError(fmt.Errorf("error creating Migrate instance: %w", err), driver, migrationsSource)
	}
	// we use this logger too in a couple of places, so need it non-nil
	m.Log = NewLogger(verboseLogging)
	return &DatabaseMigrator{wrapped: m}, nil
}

func datasourceName(username, password, host string, port int, databaseName string, schemaName string) string {
	datasource := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(username, password),
		Host:     net.JoinHostPort(host, fmt.Sprintf("%d", port)),
		Path:     databaseName,
		RawQuery: fmt.Sprintf("search_path=%s", schemaName),
	}
	return datasource.String()
}

func closeOnError(originalErr error, closers ...io.Closer) error {
	var closeErrs []string
	for _, closer := range closers {
		if closeErr := closer.Close(); closeErr != nil {
			closeErrs = append(closeErrs,
				fmt.Sprintf("in addition an error occured when closing %T: %v",
					closer,
					closeErr))
		}
	}
	if len(closeErrs) == 0 {
		return originalErr
	}
	return fmt.Errorf("%w; %s", originalErr, strings.Join(closeErrs, "; "))
}
