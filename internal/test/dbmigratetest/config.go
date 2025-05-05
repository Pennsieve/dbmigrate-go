package dbmigratetest

import (
	"github.com/pennsieve/dbmigrate-go/internal/test/configtest"
	"github.com/pennsieve/dbmigrate-go/pkg/dbmigrate"
)

// Config returns a dbmigrate.Config suitable for use against
// the pennseivedb instance started for testing. It is preferred in tests over
// calling dbmigrate.LoadConfig() because that method
// will not create the correct configs if the tests are running locally instead
// of in the Docker test container.
func Config(pgOptions ...configtest.PostgresOption) dbmigrate.Config {
	return dbmigrate.Config{
		PostgresDB:     configtest.PostgresDBConfig(pgOptions...),
		VerboseLogging: true,
	}
}
