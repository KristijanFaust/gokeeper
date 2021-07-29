package databaseutil

import (
	"fmt"
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
)

func GenerateTestDatasourceConfiguration() *config.Datasource {
	datasourceConfig := new(config.Datasource)
	datasourceConfig.Host = "localhost:50000"
	datasourceConfig.User = "gokeeperapp-test"
	datasourceConfig.Password = "password-test"
	datasourceConfig.Database = "gokeeper-test"
	datasourceConfig.MaxOpenConnections = 1
	datasourceConfig.MaxConnectionLifetime = 1

	return datasourceConfig
}

func RunDatabaseMigrations(datasourceConfig *config.Datasource) bool {
	if reflect.ValueOf(datasourceConfig).IsZero() {
		log.Panic("Datasource configuration not loaded, cannot run test database migrations")
	}

	migrationFilesPath := "file://" + getRelativePathToDatabaseMigrationFiles()
	databaseUri := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		datasourceConfig.User,
		datasourceConfig.Password,
		datasourceConfig.Host,
		datasourceConfig.Database,
	)

	migration, err := migrate.New(migrationFilesPath, databaseUri)
	if err != nil {
		log.Printf(
			"An error occured while trying to setup test migrations from: %s, with database uri: %s\nError: %s",
			migrationFilesPath, databaseUri, err,
		)
		return false
	}

	err = migration.Steps(1)
	if err != nil {
		log.Printf("An error occured during test migration execution: %s", err)
		return false
	}

	return true
}

func getRelativePathToDatabaseMigrationFiles() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename) + "/../../../../support/database/postgres/migration"
}
