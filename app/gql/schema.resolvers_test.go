package gql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/KristijanFaust/gokeeper/app/database"
	"github.com/KristijanFaust/gokeeper/app/database/repository"
	"github.com/KristijanFaust/gokeeper/app/gql/generated"
	"github.com/KristijanFaust/gokeeper/app/gql/model"
	"github.com/KristijanFaust/gokeeper/app/security"
	"github.com/KristijanFaust/gokeeper/app/utility/test/databaseutil"
	"github.com/KristijanFaust/gokeeper/app/utility/test/mock"
	"github.com/KristijanFaust/gokeeper/app/utility/test/testcontainersutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/upper/db/v4"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"testing"
)

type SchemaResolverTestSuite struct {
	suite.Suite
	session            *db.Session
	isDatabaseUp       bool
	isDatabaseMigrated bool
	mutationResolver   generated.MutationResolver
	queryResolver      generated.QueryResolver
}

func TestPasswordSuite(t *testing.T) {
	suite.Run(t, new(SchemaResolverTestSuite))
}

func (suite *SchemaResolverTestSuite) SetupSuite() {
	suite.isDatabaseUp = testcontainersutil.DockerComposeUp()
	databaseutil.GenerateTestDatasourceConfiguration()
	suite.session = database.InitializeDatabaseConnection()
	suite.isDatabaseMigrated = databaseutil.RunDatabaseMigrations()
	injectRuntimeResolverServices(suite)
}

func (suite *SchemaResolverTestSuite) TearDownSuite() {
	testcontainersutil.DockerComposeDown()
	database.CloseDatabaseConnection(suite.session)
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
	assert.Equal(suite.T(), user.Email, input.Email)
	assert.Equal(suite.T(), user.Username, input.Username)
}

// CreateUser should return error on failed input validation
func (suite *SchemaResolverTestSuite) TestCreateUserValidation() {
	input := model.NewUser{Email: "invalidEmail", Username: "", Password: ""}
	ctx := graphql.WithResponseContext(context.Background(), graphql.DefaultErrorPresenter, graphql.DefaultRecover)
	user, err := suite.mutationResolver.CreateUser(ctx, input)

	assert.Equal(
		suite.T(), err, gqlerror.Errorf("validation error/s on user input"),
		"Should return expected error when input validation for new user fails",
	)
	assert.Nil(suite.T(), user, "Should not return any user data")
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

// CreateUser should return expected error on unsuccessful user creation
func (suite *SchemaResolverTestSuite) TestCreateUserWithError() {
	injectMockedResolverServices(suite, true, false, false, false, false, false, false)
	defer injectRuntimeResolverServices(suite)

	input := model.NewUser{Email: "usercreationerror@email.com", Username: "testUsername", Password: "password"}
	user, err := suite.mutationResolver.CreateUser(context.Background(), input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not create a new user"),
		"Should return expected error when user email already exists",
	)
	assert.Nil(suite.T(), user, "Should not return any user data")
}

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
	assert.Equal(suite.T(), password.UserID, user.ID)
	assert.Equal(suite.T(), password.Name, passwordInput.Name)
	assert.Equal(suite.T(), password.Password, passwordInput.Password)
}

// CreatePassword should return error on failed input validation
func (suite *SchemaResolverTestSuite) TestCreatePasswordValidation() {
	input := model.NewPassword{UserID: "", Name: "", Password: ""}
	ctx := graphql.WithResponseContext(context.Background(), graphql.DefaultErrorPresenter, graphql.DefaultRecover)
	password, err := suite.mutationResolver.CreatePassword(ctx, input)

	assert.Equal(
		suite.T(), err, gqlerror.Errorf("validation error/s on password input"),
		"Should return expected error when input validation for new user fails",
	)
	assert.Nil(suite.T(), password, "Should not return any user data")
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

// CreatePassword should return expected error when insert to database fails
func (suite *SchemaResolverTestSuite) TestCreatePasswordWithInsertError() {
	injectMockedResolverServices(suite, false, false, false, true, false, false, false)
	defer injectRuntimeResolverServices(suite)

	passwordInput := model.NewPassword{UserID: "1", Name: "testDomain", Password: "password"}
	password, err := suite.mutationResolver.CreatePassword(context.Background(), passwordInput)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not create a new password"),
		"Should return expected error when user doesn't exist",
	)
	assert.Nil(suite.T(), password, "Should not return any password data")
}

// CreatePassword should return expected error on unsuccessful password encryption
func (suite *SchemaResolverTestSuite) TestCreatePasswordWithEncryptionError() {
	injectMockedResolverServices(suite, false, false, false, false, false, true, false)
	defer injectRuntimeResolverServices(suite)

	passwordInput := model.NewPassword{UserID: "1", Name: "testDomain", Password: "password"}
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

	assert.Equal(suite.T(), queriedUser.ID, user.ID)
	assert.Equal(suite.T(), queriedUser.Email, user.Email)
	assert.Equal(suite.T(), queriedUser.Username, user.Username)
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

// QueryUserByEmail should return expected error on unsuccessful user fetch
func (suite *SchemaResolverTestSuite) TestQueryUserByEmailWithError() {
	injectMockedResolverServices(suite, false, true, false, false, false, false, false)
	defer injectRuntimeResolverServices(suite)

	user, err := suite.queryResolver.QueryUserByEmail(context.Background(), "testemail@mail.com")
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not fetch user"),
		"Should return expected error when user email already exists",
	)
	assert.Nil(suite.T(), user, "Should not return any user data")
}

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

	assert.Equal(suite.T(), passwords[0].Name, password1.Name)
	assert.Equal(suite.T(), passwords[1].Name, password2.Name)
	assert.Equal(suite.T(), passwords[0].Password, password1.Password)
	assert.Equal(suite.T(), passwords[1].Password, password2.Password)
	for _, password := range passwords {
		assert.Equal(suite.T(), password.UserID, user.ID)
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

// QueryUserPasswords should return expected error on unsuccessful user's master password fetch
func (suite *SchemaResolverTestSuite) TestQueryUserPasswordsWithMasterPasswordFetchError() {
	injectMockedResolverServices(suite, false, false, true, false, false, false, false)
	defer injectRuntimeResolverServices(suite)

	passwords, err := suite.queryResolver.QueryUserPasswords(context.Background(), "1")
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not fetch user's passwords"),
		"Should return expected error when user email already exists",
	)
	assert.Nil(suite.T(), passwords, "Should not return any user data")
}

// QueryUserPasswords should return expected error on unsuccessful user's passwords fetch
func (suite *SchemaResolverTestSuite) TestQueryUserPasswordsWithFetchError() {
	injectMockedResolverServices(suite, false, false, false, false, true, false, false)
	defer injectRuntimeResolverServices(suite)

	passwords, err := suite.queryResolver.QueryUserPasswords(context.Background(), "1")
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not fetch user's passwords"),
		"Should return expected error when user email already exists",
	)
	assert.Nil(suite.T(), passwords, "Should not return any user data")
}

// QueryUserPasswords should return expected error on unsuccessful user's passwords decryption
func (suite *SchemaResolverTestSuite) TestQueryUserPasswordsWithDecryptionError() {
	injectMockedResolverServices(suite, false, false, false, false, false, false, true)
	defer injectRuntimeResolverServices(suite)

	passwords, err := suite.queryResolver.QueryUserPasswords(context.Background(), "1")
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not fetch user's passwords"),
		"Should return expected error when user email already exists",
	)
	assert.Nil(suite.T(), passwords, "Should not return any user data")
}

func injectMockedResolverServices(
	suite *SchemaResolverTestSuite,
	insertNewUserError bool,
	fetchByEmailError bool,
	fetchMasterPasswordByUserIdError bool,
	insertPasswordError bool,
	fetchAllByUserIdError bool,
	encryptionError bool,
	decryptionError bool,
) {
	resolver := NewResolver(
		&mock.UserRepositoryServiceMock{
			InsertNewUserError:               insertNewUserError,
			FetchByEmailError:                fetchByEmailError,
			FetchMasterPasswordByUserIdError: fetchMasterPasswordByUserIdError,
		},
		&mock.PasswordRepositoryServiceMock{InsertPasswordError: insertPasswordError, FetchAllByUserIdError: fetchAllByUserIdError},
		&mock.PasswordSecurityServiceMock{EncryptionError: encryptionError, DecryptionError: decryptionError},
	)
	suite.mutationResolver = resolver.Mutation()
	suite.queryResolver = resolver.Query()
}

func injectRuntimeResolverServices(suite *SchemaResolverTestSuite) {
	resolver := NewResolver(
		repository.NewUserRepositoryService(suite.session),
		repository.NewPasswordRepositoryService(suite.session),
		&security.PasswordSecurityService{
			Argon2PasswordHasher: &security.PasswordHashService{},
			AesPasswordCryptor:   &security.PasswordCryptoService{},
		},
	)
	suite.mutationResolver = resolver.Mutation()
	suite.queryResolver = resolver.Query()
}
