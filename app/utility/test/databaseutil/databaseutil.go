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

func GenerateTestDatasourceConfiguration() {
	config.ApplicationConfig = new(config.Config)
	config.ApplicationConfig.Datasource.Host = "localhost:50000"
	config.ApplicationConfig.Datasource.User = "gokeeperapp-test"
	config.ApplicationConfig.Datasource.Password = "password-test"
	config.ApplicationConfig.Datasource.Database = "gokeeper-test"
	config.ApplicationConfig.Datasource.MaxOpenConnections = 1
	config.ApplicationConfig.Datasource.MaxConnectionLifetime = 1
}

func RunDatabaseMigrations() bool {
	if config.ApplicationConfig == nil || reflect.ValueOf(config.ApplicationConfig.Datasource).IsZero() {
		log.Panic("Datasource configuration not loaded, cannot run test database migrations")
	}

	migrationFilesPath := "file://" + getRelativePathToDatabaseMigrationFiles()
	databaseUri := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		config.ApplicationConfig.Datasource.User,
		config.ApplicationConfig.Datasource.Password,
		config.ApplicationConfig.Datasource.Host,
		config.ApplicationConfig.Datasource.Database,
	)

	migrate, err := migrate.New(migrationFilesPath, databaseUri)
	if err != nil {
		log.Printf(
			"An error occured while trying to setup test migrations from: %s, with database uri: %s\nError: %s",
			migrationFilesPath, databaseUri, err,
		)
		return false
	}

	err = migrate.Steps(1)
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
