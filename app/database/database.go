package database

import (
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
	"reflect"
	"time"
)

var settings *postgresql.ConnectionURL

// Variables meant for mocking
var (
	pingDatabase    = db.Session.Ping
	closeConnection = db.Session.Close
)

func InitializeDatabaseConnection() *db.Session {
	if config.ApplicationConfig == nil || reflect.ValueOf(config.ApplicationConfig.Datasource).IsZero() {
		log.Panic("Datasource configuration not loaded, cannot connect to database")
	}

	settings = &postgresql.ConnectionURL{
		Database: config.ApplicationConfig.Datasource.Database,
		Host:     config.ApplicationConfig.Datasource.Host,
		User:     config.ApplicationConfig.Datasource.User,
		Password: config.ApplicationConfig.Datasource.Password,
	}

	session, err := postgresql.Open(settings)
	if err != nil {
		log.Panicf("Could not connect to database: %s", err)
	}

	session.SetMaxOpenConns(config.ApplicationConfig.Datasource.MaxOpenConnections)
	session.SetMaxIdleConns(config.ApplicationConfig.Datasource.MaxOpenConnections / 3)
	session.SetConnMaxLifetime(time.Duration(config.ApplicationConfig.Datasource.MaxConnectionLifetime) * time.Minute)

	if err := pingDatabase(session); err != nil {
		log.Panicf("Could not ping database: %s", err)
	}

	log.Printf("Successfully connected to database: %s at %s", session.Name(), settings.Host)
	return &session
}

func CloseDatabaseConnection(session *db.Session) {
	err := closeConnection(*session)
	if err != nil {
		log.Printf("Could not close connection to database\nError: %s", err)
		return
	}

	log.Printf("Successfully terminated connection to database: %s at %s", (*session).Name(), settings.Host)
}
