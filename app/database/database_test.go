package database

import (
	"github.com/KristijanFaust/gokeeper/app/config"
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

func (suite *DatabaseTestSuite) TearDownTest() {
	config.ApplicationConfig = nil
}

// InitializeDatabaseConnection should successfully connect to application database
func (suite *DatabaseTestSuite) TestInitializeDatabaseConnection() {
	if !suite.isDatabaseUp {
		suite.T().Skip("Skipping test since database container is not ready")
	}
	databaseutil.GenerateTestDatasourceConfiguration()
	assert.NotPanics(suite.T(), InitializeDatabaseConnection, "Database connections should initialize without panics")
	assert.NotNil(suite.T(), Session, "Session should be set up now")
	CloseDatabaseConnection()
}

// InitializeDatabaseConnection should panic if datasource configuration is not loaded
func (suite *DatabaseTestSuite) TestInitializeDatabaseConnectionWithoutDatasourceConfiguration() {
	assert.PanicsWithValue(
		suite.T(), "Datasource configuration not loaded, cannot connect to database",
		InitializeDatabaseConnection, "Database connection setup should panic if no datasource configuration is loaded",
	)
}

// InitializeDatabaseConnection should panic if connection to database cannot be established
func (suite *DatabaseTestSuite) TestInitializeDatabaseConnectionWithDatabaseConnectionError() {
	databaseutil.GenerateTestDatasourceConfiguration()
	config.ApplicationConfig.Datasource.Host = "invalid"
	assert.PanicsWithValue(
		suite.T(), "Could not connect to database: dial tcp: lookup invalid: Temporary failure in name resolution",
		InitializeDatabaseConnection, "Database connection setup should panic if connection to database cannot be established",
	)
}

// InitializeDatabaseConnection should panic if ping to database is unsuccessful
func (suite *DatabaseTestSuite) TestInitializeDatabaseConnectionWithUnsuccessfulPing() {
	if !suite.isDatabaseUp {
		suite.T().Skip("Skipping test since database container is not ready")
	}
	databaseutil.GenerateTestDatasourceConfiguration()
	ping = func(session db.Session) error { return db.ErrNotConnected }
	defer func() { ping = db.Session.Ping }()
	assert.PanicsWithValue(
		suite.T(), "Could not ping database: upper: not connected to a database",
		InitializeDatabaseConnection, "Database connection setup should panic if ping to database is unsuccessful",
	)
	CloseDatabaseConnection()
}

// CloseDatabaseConnection should successfully disconnect from application database
func (suite *DatabaseTestSuite) TestCloseDatabaseConncetion() {
	if !suite.isDatabaseUp {
		suite.T().Skip("Skipping test since database container is not ready")
	}
	databaseutil.GenerateTestDatasourceConfiguration()
	InitializeDatabaseConnection()
	assert.NotPanics(suite.T(), CloseDatabaseConnection, "Database connections should close without panics")
	assert.Nil(suite.T(), Session, "Session should be invalidated now")
}

// CloseDatabaseConnection should not panic on error
func (suite *DatabaseTestSuite) TestCloseDatabaseConncetionWithError() {
	if !suite.isDatabaseUp {
		suite.T().Skip("Skipping test since database container is not ready")
	}
	databaseutil.GenerateTestDatasourceConfiguration()
	InitializeDatabaseConnection()
	CloseDatabaseConnection()
	assert.NotPanics(suite.T(), CloseDatabaseConnection, "Database connections should not panic on error")
}
