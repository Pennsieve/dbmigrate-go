package test

import (
	"github.com/pennsieve/dbmigrate-go/pkg/dbmigrate"
	"github.com/stretchr/testify/require"
)

func Close(t require.TestingT, migrator *dbmigrate.DatabaseMigrator) {
	Helper(t)
	srcErr, dbErr := migrator.Close()
	require.NoError(t, srcErr)
	require.NoError(t, dbErr)
}
