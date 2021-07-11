package gql

import (
	"context"
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/KristijanFaust/gokeeper/app/database"
	"github.com/KristijanFaust/gokeeper/app/gql/model"
	"github.com/KristijanFaust/gokeeper/app/utility/test/databaseutil"
	"github.com/KristijanFaust/gokeeper/app/utility/test/testcontainersutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"testing"
)

type SchemaResolverTestSuite struct {
	suite.Suite
	isDatabaseUp       bool
	isDatabaseMigrated bool
	mutationResolver   mutationResolver
	queryResolver      queryResolver
}

func TestPasswordSuite(t *testing.T) {
	suite.Run(t, new(SchemaResolverTestSuite))
}

func (suite *SchemaResolverTestSuite) SetupSuite() {
	suite.isDatabaseUp = testcontainersutil.DockerComposeUp()
	databaseutil.GenerateTestDatasourceConfiguration()
	database.InitializeDatabaseConnection()
	suite.isDatabaseMigrated = databaseutil.RunDatabaseMigrations()
}

func (suite *SchemaResolverTestSuite) TearDownSuite() {
	testcontainersutil.DockerComposeDown()
	database.CloseDatabaseConnection()
	config.ApplicationConfig = nil
}

// CreateUser should successfully create a new user
func (suite *SchemaResolverTestSuite) TestCreateUser() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	input := model.NewUser{Email: "testuser@email.com", Username: "testUsername", Password: "password"}
	user, err := suite.mutationResolver.CreateUser(context.Background(), input)

	assert.Nil(suite.T(), err, "User should be created without errors")

	assert.NotNil(suite.T(), user.ID, "User should have an Id after creation")
	assert.Equal(suite.T(), user.Email, input.Email, "The two should be the same")
	assert.Equal(suite.T(), user.Username, input.Username, "The two should be the same")
}

// CreateUser should detect existing user emails
func (suite *SchemaResolverTestSuite) TestCreateUserWithExistingEmail() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	input := model.NewUser{Email: "testexistinguser@email.com", Username: "testUsername", Password: "password"}
	_, err := suite.mutationResolver.CreateUser(context.Background(), input)
	assert.Nil(suite.T(), err, "User should be created without errors")

	user, err := suite.mutationResolver.CreateUser(context.Background(), input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("the e-mail address is already taken"),
		"Should return expected error when user email already exists",
	)
	assert.Nil(suite.T(), user, "Should not return any user data")
}

// TODO (when repository interfaces are implemented) - CreateUser should return error on unsuccessful user creation

// CreatePassword should successfully create a new user password
func (suite *SchemaResolverTestSuite) TestCreatePassword() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	userInput := model.NewUser{Email: "testuserpass@email.com", Username: "testUsername", Password: "password"}
	user, err := suite.mutationResolver.CreateUser(context.Background(), userInput)
	assert.Nil(suite.T(), err, "User should be created without errors")

	passwordInput := model.NewPassword{UserID: user.ID, Name: "testDomain", Password: "password"}
	password, err := suite.mutationResolver.CreatePassword(context.Background(), passwordInput)
	assert.Nil(suite.T(), err, "Password should be created without errors")

	assert.NotNil(suite.T(), password.ID, "Password should have an Id after creation")
	assert.Equal(suite.T(), password.UserID, user.ID, "The two should be the same")
	assert.Equal(suite.T(), password.Name, passwordInput.Name, "The two should be the same")
	assert.Equal(suite.T(), password.Password, passwordInput.Password, "The two should be the same")
}

// CreatePassword should return expected error when userId is of an unexpected value
func (suite *SchemaResolverTestSuite) TestCreatePasswordWithUnexpectedUserIdValue() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	passwordInput := model.NewPassword{UserID: "unexpectedIdValue", Name: "testDomain", Password: "password"}
	password, err := suite.mutationResolver.CreatePassword(context.Background(), passwordInput)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not create a new password"),
		"Should return expected error when user id is of an unexpected value",
	)
	assert.Nil(suite.T(), password, "Should not return any password data")
}

// CreatePassword should return expected error when user doesn't exist
func (suite *SchemaResolverTestSuite) TestCreatePasswordWithNonexistentUser() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	passwordInput := model.NewPassword{UserID: "500", Name: "testDomain", Password: "password"}
	password, err := suite.mutationResolver.CreatePassword(context.Background(), passwordInput)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not create a new password"),
		"Should return expected error when user doesn't exist",
	)
	assert.Nil(suite.T(), password, "Should not return any password data")
}

// QueryUserByEmail should successfully query for a specific user by email
func (suite *SchemaResolverTestSuite) TestQueryUserByEmail() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	userInput := model.NewUser{Email: "testqueryuser@email.com", Username: "testUsername", Password: "password"}
	user, err := suite.mutationResolver.CreateUser(context.Background(), userInput)
	assert.Nil(suite.T(), err, "Should create user without errors")

	queriedUser, err := suite.queryResolver.QueryUserByEmail(context.Background(), user.Email)
	assert.Nil(suite.T(), err, "Should fetch user without errors")

	assert.Equal(suite.T(), queriedUser.ID, user.ID, "The two should be the same")
	assert.Equal(suite.T(), queriedUser.Email, user.Email, "The two should be the same")
	assert.Equal(suite.T(), queriedUser.Username, user.Username, "The two should be the same")
}

// QueryUserByEmail should return expected error when user doesn't exist
func (suite *SchemaResolverTestSuite) TestQueryUserByEmailWithNonexistentUser() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	user, err := suite.queryResolver.QueryUserByEmail(context.Background(), "nonexistentmail@mail.com")
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("user doesn't exist"),
		"Should return expected error when user email already exists",
	)
	assert.Nil(suite.T(), user, "Should not return any user data")
}

// TODO (when repository interfaces are implemented) - QueryUserByEmail should return error on unsuccessful user fetch

// QueryUserPasswords should successfully query for all user's passwords
func (suite *SchemaResolverTestSuite) TestQueryUserPasswords() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	userInput := model.NewUser{Email: "testquerypasswords@email.com", Username: "testUsername", Password: "password"}
	user, err := suite.mutationResolver.CreateUser(context.Background(), userInput)
	assert.Nil(suite.T(), err, "Should create user without errors")

	additionalUserInput := model.NewUser{Email: "testquerypasswords2@email.com", Username: "testUsername", Password: "password"}
	additionalUser, err := suite.mutationResolver.CreateUser(context.Background(), additionalUserInput)
	assert.Nil(suite.T(), err, "Should create user without errors")

	passwordInput := model.NewPassword{UserID: user.ID, Name: "testDomain1", Password: "password1"}
	password1, err := suite.mutationResolver.CreatePassword(context.Background(), passwordInput)
	assert.Nil(suite.T(), err, "Should create password without errors")
	passwordInput = model.NewPassword{UserID: user.ID, Name: "testDomain2", Password: "password2"}
	password2, err := suite.mutationResolver.CreatePassword(context.Background(), passwordInput)
	assert.Nil(suite.T(), err, "Should create password without errors")

	passwordInput = model.NewPassword{UserID: additionalUser.ID, Name: "testDomain3", Password: "password3"}
	_, err = suite.mutationResolver.CreatePassword(context.Background(), passwordInput)
	assert.Nil(suite.T(), err, "Should create password without errors")

	passwords, err := suite.queryResolver.QueryUserPasswords(context.Background(), user.ID)
	assert.Nil(suite.T(), err, "Should fetch passwords without errors")

	assert.Equal(suite.T(), len(passwords), 2, "Query should fetch exactly two passwords")

	assert.Equal(suite.T(), passwords[0].Name, password1.Name, "The two should be the same")
	assert.Equal(suite.T(), passwords[1].Name, password2.Name, "The two should be the same")
	assert.Equal(suite.T(), passwords[0].Password, password1.Password, "The two should be the same")
	assert.Equal(suite.T(), passwords[1].Password, password2.Password, "The two should be the same")
	for _, password := range passwords {
		assert.Equal(suite.T(), password.UserID, user.ID, "The two should be the same")
	}
}

// QueryUserPasswords should should return empty slice when user has got no passwords
func (suite *SchemaResolverTestSuite) TestQueryUserPasswordsWithoutUserPasswords() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	userInput := model.NewUser{Email: "testuserwithoutpasswords@email.com", Username: "testUsername", Password: "password"}
	user, err := suite.mutationResolver.CreateUser(context.Background(), userInput)
	assert.Nil(suite.T(), err, "Should create user without errors")

	passwords, err := suite.queryResolver.QueryUserPasswords(context.Background(), user.ID)

	assert.Nil(suite.T(), err, "Should fetch passwords without errors")
	assert.Nil(suite.T(), passwords, "Should return nil passwords slice")
}

// QueryUserPasswords should should return expected error when userId is of an unexpected value
func (suite *SchemaResolverTestSuite) TestQueryUserPasswordsWithUnexpectedUserIdValue() {
	if !suite.isDatabaseUp || !suite.isDatabaseMigrated {
		suite.T().Skip("Skipping test since database container is not ready")
	}

	passwords, err := suite.queryResolver.QueryUserPasswords(context.Background(), "invalidUserId")
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not fetch user's passwords"),
		"Should return expected error when user id is of an unexpected value",
	)
	assert.Nil(suite.T(), passwords, "Should not return any passwords data")
}

// TODO (when repository interfaces are implemented) - QueryUserByEmail should return error on unsuccessful user fetch

// Mutation should return expected mutation resolver
func (suite *SchemaResolverTestSuite) TestMutation() {
	mutationResolver := suite.mutationResolver.Mutation()
	assert.Equal(suite.T(), mutationResolver, &suite.mutationResolver)
}

// Query should return expected query resolver
func (suite *SchemaResolverTestSuite) TestQuery() {
	queryResolver := suite.queryResolver.Query()
	assert.Equal(suite.T(), queryResolver, &suite.queryResolver)
}
