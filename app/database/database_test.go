package database

import (
	"github.com/KristijanFaust/gokeeper/app/utility/test/databaseutil"
	"github.com/KristijanFaust/gokeeper/app/utility/test/testcontainersutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/upper/db/v4"
	"testing"
)

type DatabaseTestSuite struct {
	suite.Suite
	isDatabaseUp bool
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (suite *DatabaseTestSuite) SetupSuite() {
	suite.isDatabaseUp = testcontainersutil.DockerComposeUp()
}

func (suite *DatabaseTestSuite) TearDownSuite() {
	testcontainersutil.DockerComposeDown()
}

// InitializeDatabaseConnection should successfully connect to application database
func (suite *DatabaseTestSuite) TestInitializeDatabaseConnection() {
	if !suite.isDatabaseUp {
		suite.T().Skip("Skipping test since database container is not ready")
	}
	databaseConfig := databaseutil.GenerateTestDatasourceConfiguration()
	var session *db.Session
	assert.NotPanics(suite.T(), func() { session = InitializeDatabaseConnection(databaseConfig) }, "Database connections should initialize without panics")
	assert.NotNil(suite.T(), session, "Session should be set up now")
	CloseDatabaseConnection(session)
}

// InitializeDatabaseConnection should panic if datasource configuration is not loaded
func (suite *DatabaseTestSuite) TestInitializeDatabaseConnectionWithoutDatasourceConfiguration() {
	assert.PanicsWithValue(
		suite.T(), "Datasource configuration not loaded, cannot connect to database",
		func() { InitializeDatabaseConnection(nil) }, "Database connection setup should panic if no datasource configuration is loaded",
	)
}

// InitializeDatabaseConnection should panic if connection to database cannot be established
func (suite *DatabaseTestSuite) TestInitializeDatabaseConnectionWithDatabaseConnectionError() {
	databaseConfig := databaseutil.GenerateTestDatasourceConfiguration()
	databaseConfig.Host = "invalid"
	assert.PanicsWithValue(
		suite.T(), "Could not connect to database: dial tcp: lookup invalid: Temporary failure in name resolution",
		func() { InitializeDatabaseConnection(databaseConfig) }, "Database connection setup should panic if connection to database cannot be established",
	)
}

// InitializeDatabaseConnection should panic if ping to database is unsuccessful
func (suite *DatabaseTestSuite) TestInitializeDatabaseConnectionWithUnsuccessfulPing() {
	if !suite.isDatabaseUp {
		suite.T().Skip("Skipping test since database container is not ready")
	}
	databaseConfig := databaseutil.GenerateTestDatasourceConfiguration()
	pingDatabase = func(session db.Session) error { return db.ErrNotConnected }
	defer func() { pingDatabase = db.Session.Ping }()
	assert.PanicsWithValue(
		suite.T(), "Could not ping database: upper: not connected to a database",
		func() { InitializeDatabaseConnection(databaseConfig) }, "Database connection setup should panic if ping to database is unsuccessful",
	)
}

// CloseDatabaseConnection should successfully close database connection
func (suite *DatabaseTestSuite) TestCloseDatabaseConnection() {
	if !suite.isDatabaseUp {
		suite.T().Skip("Skipping test since database container is not ready")
	}
	databaseConfig := databaseutil.GenerateTestDatasourceConfiguration()
	session := InitializeDatabaseConnection(databaseConfig)
	assert.NotPanics(suite.T(), func() { CloseDatabaseConnection(session) }, "Database connections should close without panics")
}

// CloseDatabaseConnection should not panic on error
func (suite *DatabaseTestSuite) TestCloseDatabaseConnectionWithError() {
	if !suite.isDatabaseUp {
		suite.T().Skip("Skipping test since database container is not ready")
	}
	databaseConfig := databaseutil.GenerateTestDatasourceConfiguration()
	session := InitializeDatabaseConnection(databaseConfig)
	closeConnection = func(session db.Session) error { return db.ErrNotConnected }
	assert.NotPanics(suite.T(), func() { CloseDatabaseConnection(session) }, "Database connections should not panic on error")
	closeConnection = db.Session.Close
	CloseDatabaseConnection(session)
}
