package database

import (
	"errors"
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
	"reflect"
	"time"
)

var settings *postgresql.ConnectionURL
var Session db.Session
var ping = db.Session.Ping

func InitializeDatabaseConnection() {
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

	Session = session

	if err := ping(Session); err != nil {
		log.Panicf("Could not ping database: %s", err)
	}

	log.Printf("Successfully connected to database: %s at %s", session.Name(), settings.Host)
}

func CloseDatabaseConnection() {
	var err error
	if Session != nil {
		err = Session.Close()
	} else {
		err = errors.New("session is nil (connection to database is probably already closed)")
	}

	if err != nil {
		log.Printf("Could not close connection to database\nError: %s", err)
		return
	}

	log.Printf("Successfully terminated connection to database: %s at %s", Session.Name(), settings.Host)
	Session = nil
}
