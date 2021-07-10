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

type PasswordTestSuite struct {
	suite.Suite
	isDatabaseUp       bool
	isDatabaseMigrated bool
}

func TestPasswordSuite(t *testing.T) {
	suite.Run(t, new(PasswordTestSuite))
}

func (suite *PasswordTestSuite) SetupSuite() {
	suite.isDatabaseUp = testcontainersutil.DockerComposeUp()
	databaseutil.GenerateTestDatasourceConfiguration()
	database.InitializeDatabaseConnection()
	suite.isDatabaseMigrated = databaseutil.RunDatabaseMigrations()
}

func (suite *PasswordTestSuite) TearDownSuite() {
	testcontainersutil.DockerComposeDown()
	database.CloseDatabaseConnection()
	config.ApplicationConfig = nil
}

// InsertNewPassword should successfully insert a new user password in the database
func (suite *PasswordTestSuite) TestInsertNewUser() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	var user = User{Email: "testInsertPassword@test.com", Username: "testInsertPassword", Password: "testInsertPassword"}
	userId, err := user.InsertNewUser()
	assert.Nil(suite.T(), err)

	var newUserPassword = Password{UserId: uint64(userId.ID().(int64)), Name: "SomeApplication", Password: "password"}
	passwordId, err := newUserPassword.InsertNewPassword()
	assert.Nil(suite.T(), err)

	insertedUserPassword := Password{}
	err = PasswordCollection().Find("id", passwordId).One(&insertedUserPassword)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), insertedUserPassword.Id, uint64(passwordId.ID().(int64)), "The two should be the same")
	assert.Equal(suite.T(), insertedUserPassword.UserId, newUserPassword.UserId, "The two should be the same")
	assert.Equal(suite.T(), insertedUserPassword.Name, newUserPassword.Name, "The two should be the same")
	assert.Equal(suite.T(), insertedUserPassword.Password, newUserPassword.Password, "The two should be the same")
}

// FetchAllByUserId should successfully fetch all user's password from the database
func (suite *PasswordTestSuite) TestFetchAllByUserId() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	var testUser = User{Email: "testFetchPasswords@test.com", Username: "testFetchPasswords", Password: "testFetchPassword"}
	testUserId, err := testUser.InsertNewUser()
	assert.Nil(suite.T(), err)

	var testUserPassword1 = Password{UserId: uint64(testUserId.ID().(int64)), Name: "SomeApplication1", Password: "password1"}
	passwordId1, err := testUserPassword1.InsertNewPassword()
	assert.Nil(suite.T(), err)

	var testUserPassword2 = Password{UserId: uint64(testUserId.ID().(int64)), Name: "SomeApplication2", Password: "password2"}
	passwordId2, err := testUserPassword2.InsertNewPassword()
	assert.Nil(suite.T(), err)

	var additionalUser = User{Email: "additionalUser@test.com", Username: "additionalUser", Password: "additionalUser"}
	additionalUserId, err := additionalUser.InsertNewUser()
	assert.Nil(suite.T(), err)

	var additionalUserPassword = Password{UserId: uint64(additionalUserId.ID().(int64)), Name: "SomeApplication", Password: "password"}
	_, err = additionalUserPassword.InsertNewPassword()
	assert.Nil(suite.T(), err)

	testUserPasswords := Passwords{}
	testUserPasswords.FetchAllByUserId(uint64(testUserId.ID().(int64)))

	assert.Equal(suite.T(), len(testUserPasswords), 2, "Should fetch exactly two passwords")

	assert.Equal(suite.T(), testUserPasswords[0].Id, uint64(passwordId1.ID().(int64)), "The two should be the same")
	assert.Equal(suite.T(), testUserPasswords[1].Id, uint64(passwordId2.ID().(int64)), "The two should be the same")
	assert.Equal(suite.T(), testUserPasswords[0].UserId, testUserPasswords[1].UserId, "The passwords should belong to the same user")
	assert.Equal(suite.T(), testUserPasswords[0].Name, testUserPassword1.Name, "The two should be the same")
	assert.Equal(suite.T(), testUserPasswords[1].Name, testUserPassword2.Name, "The two should be the same")
	assert.Equal(suite.T(), testUserPasswords[0].Password, testUserPassword1.Password, "The two should be the same")
	assert.Equal(suite.T(), testUserPasswords[1].Password, testUserPassword2.Password, "The two should be the same")
}
