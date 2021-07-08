package model

import (
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/KristijanFaust/gokeeper/app/database"
	"github.com/KristijanFaust/gokeeper/app/utility/test/databaseutil"
	"github.com/KristijanFaust/gokeeper/app/utility/test/testcontainersutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserTestSuite struct {
	suite.Suite
	isDatabaseUp       bool
	isDatabaseMigrated bool
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) SetupSuite() {
	suite.isDatabaseUp = testcontainersutil.DockerComposeUp()
	databaseutil.GenerateTestDatasourceConfiguration()
	database.InitializeDatabaseConnection()
	suite.isDatabaseMigrated = databaseutil.RunDatabaseMigrations()
}

func (suite *UserTestSuite) TearDownSuite() {
	testcontainersutil.DockerComposeDown()
	database.CloseDatabaseConnection()
	config.ApplicationConfig = nil
}

// InsertNewUser should successfully insert a new user in the database
func (suite *UserTestSuite) TestInsertNewUser() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	var newUser = User{Email: "testUsername@test.com", Username: "testUsername", Password: "testPassword"}
	newUserInsertResult, err := newUser.InsertNewUser()
	assert.Nil(suite.T(), err)

	insertedUser := User{}
	err = UserCollection().Find("id", newUserInsertResult).One(&insertedUser)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), insertedUser.Id, uint64(newUserInsertResult.ID().(int64)), "The two should be the same")
	assert.Equal(suite.T(), insertedUser.Email, newUser.Email, "The two should be the same")
	assert.Equal(suite.T(), insertedUser.Username, newUser.Username, "The two should be the same")
	assert.Equal(suite.T(), insertedUser.Password, newUser.Password, "The two should be the same")
}

// FindByEmail should successfully fetch an existing user by email from the database
func (suite *UserTestSuite) TestFetchByEmail() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	var newUser = User{Email: "testFetchUsername@test.com", Username: "testFetchUsername", Password: "testFetchPassword"}

	newUserInsertResult, err := newUser.InsertNewUser()
	assert.Nil(suite.T(), err)

	targetUser := User{}
	err = targetUser.FetchByEmail(newUser.Email)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), targetUser.Id, uint64(newUserInsertResult.ID().(int64)), "The two should be the same")
	assert.Equal(suite.T(), targetUser.Email, newUser.Email, "The two should be the same")
	assert.Equal(suite.T(), targetUser.Username, newUser.Username, "The two should be the same")
	assert.Equal(suite.T(), targetUser.Password, newUser.Password, "The two should be the same")
}
