# dbmigrate-go
Database Migrations for Pennsieve Go-based Services

This library is a wrapper around [golang-migrate](https://github.com/golang-migrate/migrate).

See [migrate_test.go](pkg/dbmigrate/migrate_test.go), [collections-service](https://github.com/Pennsieve/collections-service), 
and [github-service](https://github.com/Pennsieve/github-service) for examples of use.

See [config.go](pkg/config/config.go) and [postgres.go](pkg/config/postgres.go) for environment variables used to
configure the database connection for migrations.

You will also need to create a [Migration Source](https://github.com/golang-migrate/migrate?tab=readme-ov-file#migration-sources)
to read migration files. The examples above all use `io/fs` but other migration source types are available.

See [Migration Files](https://github.com/golang-migrate/migrate?tab=readme-ov-file#migration-files) for naming and writing migration files.