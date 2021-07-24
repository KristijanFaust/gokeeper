package repository

import (
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/KristijanFaust/gokeeper/app/database"
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/KristijanFaust/gokeeper/app/utility/test/databaseutil"
	"github.com/KristijanFaust/gokeeper/app/utility/test/testcontainersutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	isDatabaseUp       bool
	isDatabaseMigrated bool
	userRepository     UserRepository
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	suite.isDatabaseUp = testcontainersutil.DockerComposeUp()
	databaseutil.GenerateTestDatasourceConfiguration()
	database.InitializeDatabaseConnection()
	suite.isDatabaseMigrated = databaseutil.RunDatabaseMigrations()
	suite.userRepository = &UserRepositoryService{}
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	testcontainersutil.DockerComposeDown()
	database.CloseDatabaseConnection()
	config.ApplicationConfig = nil
}

// InsertNewUser should successfully insert a new user in the database
func (suite *UserRepositoryTestSuite) TestInsertNewUser() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	newUser := &model.User{Email: "testUsername@test.com", Username: "testUsername", Password: []byte("testPassword")}
	newUserInsertResult, err := suite.userRepository.InsertNewUser(newUser)
	assert.Nil(suite.T(), err)

	insertedUser := model.User{}
	err = UserCollection().Find("id", newUserInsertResult).One(&insertedUser)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), insertedUser.Id, uint64(newUserInsertResult.ID().(int64)))
	assert.Equal(suite.T(), insertedUser.Email, newUser.Email)
	assert.Equal(suite.T(), insertedUser.Username, newUser.Username)
	assert.Equal(suite.T(), insertedUser.Password, newUser.Password)
}

// FetchByEmail should successfully fetch an existing user by email from the database
func (suite *UserRepositoryTestSuite) TestFetchByEmail() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	newUser := &model.User{Email: "testFetchUsername@test.com", Username: "testFetchUsername", Password: []byte("testFetchPassword")}

	newUserInsertResult, err := suite.userRepository.InsertNewUser(newUser)
	assert.Nil(suite.T(), err)

	targetUser := &model.User{}
	err = suite.userRepository.FetchByEmail(targetUser, newUser.Email)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), targetUser.Id, uint64(newUserInsertResult.ID().(int64)))
	assert.Equal(suite.T(), targetUser.Email, newUser.Email)
	assert.Equal(suite.T(), targetUser.Username, newUser.Username)
	assert.Equal(suite.T(), targetUser.Password, newUser.Password)
}

// FetchMasterPasswordByUserId should successfully fetch an existing user's master password by id from the database
func (suite *UserRepositoryTestSuite) TestFetchMasterPasswordByUserId() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	newUser := &model.User{Email: "testFetchPassword@test.com", Username: "testFetchPassword", Password: []byte("testFetchPassword")}

	newUserInsertResult, err := suite.userRepository.InsertNewUser(newUser)
	assert.Nil(suite.T(), err)

	targetUser := &model.User{}
	err = suite.userRepository.FetchMasterPasswordByUserId(targetUser, uint64(newUserInsertResult.ID().(int64)))
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), targetUser.Id, uint64(0))
	assert.Equal(suite.T(), targetUser.Email, "")
	assert.Equal(suite.T(), targetUser.Username, "")
	assert.Equal(suite.T(), targetUser.Password, newUser.Password)
}
