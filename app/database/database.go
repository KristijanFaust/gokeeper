package database

import (
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
	"reflect"
	"time"
)

// Variables meant for mocking
var (
	pingDatabase    = db.Session.Ping
	closeConnection = db.Session.Close
)

func InitializeDatabaseConnection(datasourceConfig *config.Datasource) *db.Session {
	if reflect.ValueOf(datasourceConfig).IsZero() {
		log.Panic("Datasource configuration not loaded, cannot connect to database")
	}

	settings := &postgresql.ConnectionURL{
		Database: datasourceConfig.Database,
		Host:     datasourceConfig.Host,
		User:     datasourceConfig.User,
		Password: datasourceConfig.Password,
	}

	session, err := postgresql.Open(settings)
	if err != nil {
		log.Panicf("Could not connect to database: %s", err)
	}

	session.SetMaxOpenConns(datasourceConfig.MaxOpenConnections)
	session.SetMaxIdleConns(datasourceConfig.MaxOpenConnections / 3)
	session.SetConnMaxLifetime(time.Duration(datasourceConfig.MaxConnectionLifetime) * time.Minute)

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
	}

	log.Printf(
		"Successfully terminated connection to database: %s at %s",
		(*session).Name(), (*session).ConnectionURL().(*postgresql.ConnectionURL).Host,
	)
}
