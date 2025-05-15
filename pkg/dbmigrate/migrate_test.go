package dbmigrate_test

import (
	"context"
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pennsieve/dbmigrate-go/internal/test"
	"github.com/pennsieve/dbmigrate-go/pkg/config"
	"github.com/pennsieve/dbmigrate-go/pkg/dbmigrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//go:embed testdata/migrations/*.sql
var migrationsFS embed.FS

const schema = "test_schema"

func TestDatabaseMigrator(t *testing.T) {
	tests := []struct {
		scenario string
		tstFunc  func(t *testing.T, migrator *dbmigrate.DatabaseMigrator, verificationConn *pgx.Conn)
	}{
		{"test Up", testUp},
		{"test Migrate", testMigrate},
		{"Up and Down run without error", testUpAndDown},
	}

	ctx := context.Background()

	testSettings := test.NewTestSettings(schema)
	migrateConfig, err := config.LoadConfig(testSettings)
	require.NoError(t, err)

	migrationsSource, err := iofs.New(migrationsFS, "testdata/migrations")
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			// make a migrator for each test and pass it into the function so that
			// we can take care of cleaning it up here
			migrator, err := dbmigrate.NewLocalMigrator(ctx, migrateConfig, migrationsSource)
			require.NoError(t, err)

			// also pass in a plain pgx.Conn to let the test function run any verifications on the migrated schema
			verificationConn, err := test.NewPostgresDBFromConfig(t, migrateConfig.PostgresDB).Connect(ctx, migrateConfig.PostgresDB.Database)
			require.NoError(t, err)

			t.Cleanup(func() {
				require.NoError(t, migrator.Drop())
				test.Close(t, migrator)
				test.CloseConnection(ctx, t, verificationConn)
			})

			tt.tstFunc(t, migrator, verificationConn)
		})
	}
}

func testUp(t *testing.T, migrator *dbmigrate.DatabaseMigrator, verificationConn *pgx.Conn) {

	require.NoError(t, migrator.Up())

	tableIdentifier := pgx.Identifier{schema, "test_table"}.Sanitize()

	insertQuery := fmt.Sprintf(`INSERT INTO %s (name, description, node_id) VALUES (@name, @description, @node_id) RETURNING id, created_at, updated_at`,
		tableIdentifier)
	ctx := context.Background()
	var id int64
	var createdAt, updatedAt time.Time
	require.NoError(t,
		verificationConn.QueryRow(ctx,
			insertQuery,
			pgx.NamedArgs{
				"name":        uuid.NewString(),
				"description": uuid.NewString(),
				"node_id":     uuid.NewString()}).
			Scan(&id, &createdAt, &updatedAt),
	)
	assert.False(t, createdAt.IsZero())
	assert.False(t, updatedAt.IsZero())

	verificationQuery := fmt.Sprintf(`UPDATE %s SET description = @description WHERE id = @id RETURNING updated_at`,
		tableIdentifier)
	var updatedUpdatedAt time.Time
	require.NoError(t,
		verificationConn.QueryRow(ctx,
			verificationQuery,
			pgx.NamedArgs{
				"description": uuid.NewString(),
				"id":          id,
			}).
			Scan(&updatedUpdatedAt),
	)
	assert.False(t, updatedUpdatedAt.IsZero())
	assert.False(t, updatedAt.Equal(updatedUpdatedAt))
}

func testMigrate(t *testing.T, migrator *dbmigrate.DatabaseMigrator, verificationConn *pgx.Conn) {
	ctx := context.Background()

	// Only run the first migration
	require.NoError(t, migrator.Migrate(20250319124829))

	expectedFunctionName := fmt.Sprintf("%s.update_updated_at_column", schema)

	var functionName string
	require.NoError(t, verificationConn.QueryRow(ctx, fmt.Sprintf(`SELECT to_regproc('%s')`, expectedFunctionName)).Scan(&functionName))
	assert.NotNil(t, functionName)
	assert.Equal(t, expectedFunctionName, functionName)

	expectedTableName := fmt.Sprintf("%s.test_table", schema)
	tableExistsQuery := fmt.Sprintf(`SELECT to_regclass('%s')`, expectedTableName)

	var tableName *string
	require.NoError(t, verificationConn.QueryRow(ctx, tableExistsQuery).Scan(&tableName))
	assert.Nil(t, tableName)

	// now run all the remaining migrations
	require.NoError(t, migrator.Up())

	require.NoError(t, verificationConn.QueryRow(ctx, tableExistsQuery).Scan(&tableName))
	require.NotNil(t, tableName)
	require.Equal(t, expectedTableName, *tableName)

}

// We don't really use the Down() method for real. Test is here so that
// if we do write 'down' files something checks that they at least run
// without error.
func testUpAndDown(t *testing.T, migrator *dbmigrate.DatabaseMigrator, _ *pgx.Conn) {

	require.NoError(t, migrator.Up())

	require.NoError(t, migrator.Down())
}
