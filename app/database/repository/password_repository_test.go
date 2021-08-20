package repository

import (
	"github.com/KristijanFaust/gokeeper/app/database"
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/KristijanFaust/gokeeper/app/utility/test/databaseutil"
	"github.com/KristijanFaust/gokeeper/app/utility/test/testcontainersutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/upper/db/v4"
	"testing"
)

type PasswordTestSuite struct {
	suite.Suite
	session            *db.Session
	isDatabaseUp       bool
	isDatabaseMigrated bool
	userRepository     UserRepository
	passwordRepository PasswordRepository
}

func TestPasswordSuite(t *testing.T) {
	suite.Run(t, new(PasswordTestSuite))
}

func (suite *PasswordTestSuite) SetupSuite() {
	suite.isDatabaseUp = testcontainersutil.DockerComposeUp()
	databaseConfiguration := databaseutil.GenerateTestDatasourceConfiguration()
	suite.session = database.InitializeDatabaseConnection(databaseConfiguration)
	suite.isDatabaseMigrated = databaseutil.RunDatabaseMigrations(databaseConfiguration)
	suite.userRepository = NewUserRepositoryService(suite.session)
	suite.passwordRepository = NewPasswordRepositoryService(suite.session)
}

func (suite *PasswordTestSuite) TearDownSuite() {
	testcontainersutil.DockerComposeDown()
	database.CloseDatabaseConnection(suite.session)
}

// InsertNewPassword should successfully insert a new user password in the database
func (suite *PasswordTestSuite) TestInsertNewPassword() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	user := &model.User{Email: "testInsertPassword@test.com", Username: "testInsertPassword", Password: []byte("testInsertPassword")}
	userId, err := suite.userRepository.InsertNewUser(user)
	assert.Nil(suite.T(), err)

	var newUserPassword = &model.Password{UserId: uint64(userId.ID().(int64)), Name: "SomeApplication", Password: []byte("password")}
	passwordId, err := suite.passwordRepository.InsertNewPassword(newUserPassword)
	assert.Nil(suite.T(), err)

	insertedUserPassword := model.Password{}
	err = (*suite.session).Collection("password").Find("id", passwordId).One(&insertedUserPassword)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), insertedUserPassword.Id, uint64(passwordId.ID().(int64)))
	assert.Equal(suite.T(), insertedUserPassword.UserId, newUserPassword.UserId)
	assert.Equal(suite.T(), insertedUserPassword.Name, newUserPassword.Name)
	assert.Equal(suite.T(), insertedUserPassword.Password, newUserPassword.Password)
}

// UpdatePasswordById should successfully update a password record
func (suite *PasswordTestSuite) TestUpdatePasswordById() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	user := &model.User{Email: "testUpdatePassword@test.com", Username: "testInsertPassword", Password: []byte("testInsertPassword")}
	userId, err := suite.userRepository.InsertNewUser(user)

	newUserPassword := &model.Password{UserId: uint64(userId.ID().(int64)), Name: "SomeApplication", Password: []byte("password")}
	passwordId, err := suite.passwordRepository.InsertNewPassword(newUserPassword)

	suite.passwordRepository.UpdatePasswordById("UpdatedName", []byte("updatedPassword"), uint64(passwordId.ID().(int64)))

	updatedUserPassword := model.Password{}
	err = (*suite.session).Collection("password").Find("id", passwordId).One(&updatedUserPassword)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), updatedUserPassword.Id, uint64(passwordId.ID().(int64)))
	assert.Equal(suite.T(), updatedUserPassword.UserId, newUserPassword.UserId)
	assert.Equal(suite.T(), updatedUserPassword.Name, "UpdatedName")
	assert.Equal(suite.T(), updatedUserPassword.Password, []byte("updatedPassword"))
}

// FetchPasswordById should fetch user password by id
func (suite *PasswordTestSuite) TestFetchPasswordById() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	testUser := &model.User{Email: "testFetchPasswordById@test.com", Username: "testFetchPassword", Password: []byte("testFetchPassword")}
	testUserId, err := suite.userRepository.InsertNewUser(testUser)

	userPassword := &model.Password{UserId: uint64(testUserId.ID().(int64)), Name: "SomeApplication", Password: []byte("password")}
	passwordId, err := suite.passwordRepository.InsertNewPassword(userPassword)

	fetchedUserPassword := &model.Password{}
	err = suite.passwordRepository.FetchPasswordById(fetchedUserPassword, uint64(passwordId.ID().(int64)))
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), fetchedUserPassword.Id, uint64(passwordId.ID().(int64)))
	assert.Equal(suite.T(), fetchedUserPassword.UserId, userPassword.UserId)
	assert.Equal(suite.T(), fetchedUserPassword.Name, userPassword.Name)
	assert.Equal(suite.T(), fetchedUserPassword.Password, userPassword.Password)
}

// FetchAllByUserId should successfully fetch all user's password from the database
func (suite *PasswordTestSuite) TestFetchAllByUserId() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	testUser := &model.User{Email: "testFetchPasswords@test.com", Username: "testFetchPasswords", Password: []byte("testFetchPassword")}
	testUserId, err := suite.userRepository.InsertNewUser(testUser)
	assert.Nil(suite.T(), err)

	testUserPassword1 := &model.Password{UserId: uint64(testUserId.ID().(int64)), Name: "SomeApplication1", Password: []byte("password1")}
	passwordId1, err := suite.passwordRepository.InsertNewPassword(testUserPassword1)
	assert.Nil(suite.T(), err)

	testUserPassword2 := &model.Password{UserId: uint64(testUserId.ID().(int64)), Name: "SomeApplication2", Password: []byte("password2")}
	passwordId2, err := suite.passwordRepository.InsertNewPassword(testUserPassword2)
	assert.Nil(suite.T(), err)

	additionalUser := &model.User{Email: "additionalUser@test.com", Username: "additionalUser", Password: []byte("additionalUser")}
	additionalUserId, err := suite.userRepository.InsertNewUser(additionalUser)
	assert.Nil(suite.T(), err)

	additionalUserPassword := &model.Password{UserId: uint64(additionalUserId.ID().(int64)), Name: "SomeApplication", Password: []byte("password")}
	_, err = suite.passwordRepository.InsertNewPassword(additionalUserPassword)
	assert.Nil(suite.T(), err)

	testUserPasswords := model.Passwords{}
	err = suite.passwordRepository.FetchAllByUserId(&testUserPasswords, uint64(testUserId.ID().(int64)), nil)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), len(testUserPasswords), 2, "Should fetch exactly two passwords")

	assert.Equal(suite.T(), testUserPasswords[0].Id, uint64(passwordId1.ID().(int64)))
	assert.Equal(suite.T(), testUserPasswords[1].Id, uint64(passwordId2.ID().(int64)))
	assert.Equal(suite.T(), testUserPasswords[0].UserId, testUserPasswords[1].UserId, "The passwords should belong to the same user")
	assert.Equal(suite.T(), testUserPasswords[0].Name, testUserPassword1.Name)
	assert.Equal(suite.T(), testUserPasswords[1].Name, testUserPassword2.Name)
	assert.Equal(suite.T(), testUserPasswords[0].Password, testUserPassword1.Password)
	assert.Equal(suite.T(), testUserPasswords[1].Password, testUserPassword2.Password)
}

// FetchAllByUserId should only fetch requested columns from the database
func (suite *PasswordTestSuite) TestFetchAllByUserIdWithSpecificFields() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	testUser := &model.User{Email: "testSpecificFieldsFetch@test.com", Username: "testFetchPasswords", Password: []byte("testFetchPassword")}
	testUserId, err := suite.userRepository.InsertNewUser(testUser)
	assert.Nil(suite.T(), err)

	testUserPassword := &model.Password{UserId: uint64(testUserId.ID().(int64)), Name: "SomeApplication", Password: []byte("password")}
	passwordId, err := suite.passwordRepository.InsertNewPassword(testUserPassword)
	assert.Nil(suite.T(), err)

	testUserPasswords := model.Passwords{}
	err = suite.passwordRepository.FetchAllByUserId(&testUserPasswords, uint64(testUserId.ID().(int64)), []string{"id", "password"})
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), testUserPasswords[0].Id, uint64(passwordId.ID().(int64)))
	assert.Equal(suite.T(), testUserPasswords[0].UserId, uint64(0))
	assert.Equal(suite.T(), testUserPasswords[0].Name, "")
	assert.Equal(suite.T(), testUserPasswords[0].Password, testUserPassword.Password)
}
