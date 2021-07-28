package gql

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/KristijanFaust/gokeeper/app/gql/generated"
	"github.com/KristijanFaust/gokeeper/app/gql/model"
	"github.com/KristijanFaust/gokeeper/app/utility/test/mockutil"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"testing"
)

type schemaResolverTestSuite struct {
	suite.Suite
	resolver              Resolver
	mutationResolver      generated.MutationResolver
	queryResolver         generated.QueryResolver
	graphqlRequestContext context.Context
}

func TestSchemaResolverSuite(t *testing.T) {
	suite.Run(t, new(schemaResolverTestSuite))
}

func (suite *schemaResolverTestSuite) SetupSuite() {
	suite.graphqlRequestContext = graphql.WithFieldContext(
		context.WithValue(
			graphql.WithOperationContext(context.Background(), &graphql.OperationContext{}),
			"operation_context", []string{},
		),
		&graphql.FieldContext{},
	)
}

func (suite *schemaResolverTestSuite) SetupTest() {
	injectDefaultMockedResolverServices(suite)
}

// SignUp should successfully create a new user
func (suite *schemaResolverTestSuite) TestSignUp() {
	input := model.NewUser{Email: mockutil.DefaultEmail, Username: mockutil.DefaultUsername, Password: mockutil.DefaultPassword}

	user, err := suite.mutationResolver.SignUp(context.Background(), input)
	assert.Nil(suite.T(), err, "User should be created without errors")

	assert.Equal(suite.T(), user.ID, mockutil.DefaultIdAsString)
	assert.Equal(suite.T(), user.Email, input.Email)
	assert.Equal(suite.T(), user.Username, input.Username)
}

// SignUp should return error on failed input validation
func (suite *schemaResolverTestSuite) TestSignUpValidation() {
	input := model.NewUser{Email: "invalidEmail", Username: "", Password: ""}
	ctx := graphql.WithResponseContext(context.Background(), graphql.DefaultErrorPresenter, graphql.DefaultRecover)

	user, err := suite.mutationResolver.SignUp(ctx, input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("validation error/s on user input"),
		"Should return expected error when input validation for new user fails",
	)
	assert.Nil(suite.T(), user, "Should not return any user data")
}

// SignUp should detect existing user emails
func (suite *schemaResolverTestSuite) TestSignUpWithExistingEmail() {
	userRepositoryServiceMock := new(mockutil.UserRepositoryServiceMock)
	userRepositoryServiceMock.On("InsertNewUser", mock.Anything).Return(nil, &pq.Error{Code: "23505"}).Times(1)
	suite.resolver.userRepository = userRepositoryServiceMock
	input := model.NewUser{Email: mockutil.DefaultEmail, Username: mockutil.DefaultUsername, Password: mockutil.DefaultPassword}

	user, err := suite.mutationResolver.SignUp(context.Background(), input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("the e-mail address is already taken"),
		"Should return expected error when user email already exists",
	)
	assert.Nil(suite.T(), user, "Should not return any user data")
}

// SignUp should return expected error on unsuccessful user creation
func (suite *schemaResolverTestSuite) TestSignUpWithInsertError() {
	userRepositoryServiceMock := new(mockutil.UserRepositoryServiceMock)
	userRepositoryServiceMock.On("InsertNewUser", mock.Anything).Return(nil, errors.New(mockutil.MockedGenericErrorMessage)).Times(1)
	suite.resolver.userRepository = userRepositoryServiceMock
	input := model.NewUser{Email: mockutil.DefaultEmail, Username: mockutil.DefaultUsername, Password: mockutil.DefaultPassword}

	user, err := suite.mutationResolver.SignUp(context.Background(), input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not create a new user"),
		"Should return expected error when insert to database fails",
	)
	assert.Nil(suite.T(), user, "Should not return any user data")
}

// CreatePassword should successfully create a new user password
func (suite *schemaResolverTestSuite) TestCreatePassword() {
	input := model.NewPassword{UserID: mockutil.DefaultIdAsString, Name: mockutil.DefaultPasswordName, Password: mockutil.DefaultPassword}

	password, err := suite.mutationResolver.CreatePassword(context.Background(), input)
	assert.Nil(suite.T(), err, "Password should be created without errors")

	assert.Equal(suite.T(), password.ID, mockutil.DefaultIdAsString)
	assert.Equal(suite.T(), password.UserID, mockutil.DefaultIdAsString)
	assert.Equal(suite.T(), password.Name, input.Name)
	assert.Equal(suite.T(), password.Password, input.Password)
}

// CreatePassword should return error on failed input validation
func (suite *schemaResolverTestSuite) TestCreatePasswordValidation() {
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
func (suite *schemaResolverTestSuite) TestCreatePasswordWithUnexpectedUserIdValue() {
	input := model.NewPassword{UserID: "invalidId", Name: mockutil.DefaultPasswordName, Password: mockutil.DefaultPassword}

	password, err := suite.mutationResolver.CreatePassword(context.Background(), input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not create a new password"),
		"Should return expected error when user id is of an unexpected value",
	)
	assert.Nil(suite.T(), password, "Should not return any password data")
}

// CreatePassword should return expected error when user's master password fetch fails
func (suite *schemaResolverTestSuite) TestCreatePasswordWithMasterPasswordFetchError() {
	userRepositoryServiceMock := new(mockutil.UserRepositoryServiceMock)
	userRepositoryServiceMock.On("FetchMasterPasswordByUserId", mock.Anything, mock.Anything).Return(
		errors.New(mockutil.MockedGenericErrorMessage),
	).Times(1)
	suite.resolver.userRepository = userRepositoryServiceMock
	input := model.NewPassword{UserID: mockutil.DefaultIdAsString, Name: mockutil.DefaultPasswordName, Password: mockutil.DefaultPassword}

	password, err := suite.mutationResolver.CreatePassword(context.Background(), input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not create a new password"),
		"Should return expected error when user's master password fetch fails",
	)
	assert.Nil(suite.T(), password, "Should not return any password data")
}

// CreatePassword should return expected error when insert to database fails
func (suite *schemaResolverTestSuite) TestCreatePasswordWithInsertError() {
	passwordRepositoryServiceMock := new(mockutil.PasswordRepositoryServiceMock)
	passwordRepositoryServiceMock.On("InsertNewPassword", mock.Anything).Return(
		nil, errors.New(mockutil.MockedGenericErrorMessage),
	).Times(1)
	suite.resolver.passwordRepository = passwordRepositoryServiceMock
	input := model.NewPassword{UserID: mockutil.DefaultIdAsString, Name: mockutil.DefaultPasswordName, Password: mockutil.DefaultPassword}

	password, err := suite.mutationResolver.CreatePassword(context.Background(), input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not create a new password"),
		"Should return expected error when insert to database fails",
	)
	assert.Nil(suite.T(), password, "Should not return any password data")
}

// CreatePassword should return expected error on unsuccessful password encryption
func (suite *schemaResolverTestSuite) TestCreatePasswordWithEncryptionError() {
	passwordSecurityServiceMock := new(mockutil.PasswordSecurityServiceMock)
	passwordSecurityServiceMock.On("EncryptWithAes", mock.Anything, mock.Anything).Return(
		nil, errors.New(mockutil.MockedGenericErrorMessage),
	).Times(1)
	suite.resolver.passwordSecurityService = passwordSecurityServiceMock
	input := model.NewPassword{UserID: mockutil.DefaultIdAsString, Name: mockutil.DefaultPasswordName, Password: mockutil.DefaultPassword}

	password, err := suite.mutationResolver.CreatePassword(context.Background(), input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not create a new password"),
		"Should return expected error when password encryption fails",
	)
	assert.Nil(suite.T(), password, "Should not return any password data")
}

// SignIn should successfully sign in a user
func (suite *schemaResolverTestSuite) TestSignIn() {
	input := model.UserSignIn{Email: mockutil.DefaultEmail, Password: mockutil.DefaultPassword}

	userWithToken, err := suite.mutationResolver.SignIn(context.Background(), input)
	assert.Nil(suite.T(), err, "User should sign in without any errors")

	assert.Equal(suite.T(), userWithToken.Token.Token, mockutil.MockedJwtToken)
	assert.NotNil(suite.T(), userWithToken.Token.ExpireAt)

	assert.Equal(suite.T(), userWithToken.User.ID, mockutil.DefaultIdAsString)
	assert.Equal(suite.T(), userWithToken.User.Email, mockutil.DefaultEmail)
	assert.Equal(suite.T(), userWithToken.User.Username, mockutil.DefaultUsername)
}

// SignIn should return expected error when a non existing user is trying to sign in
func (suite *schemaResolverTestSuite) TestSignInWithNonExistingUser() {
	userRepositoryServiceMock := new(mockutil.UserRepositoryServiceMock)
	userRepositoryServiceMock.On("FetchByEmail", mock.Anything, mock.Anything, []string(nil)).Return(
		errors.New("upper: no more rows in this result set"),
	).Times(1)
	suite.resolver.userRepository = userRepositoryServiceMock
	input := model.UserSignIn{Email: mockutil.DefaultEmail, Password: mockutil.DefaultPassword}

	token, err := suite.mutationResolver.SignIn(context.Background(), input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("user doesn't exist"),
		"Should return expected error when a non existing user is signing in",
	)
	assert.Nil(suite.T(), token, "Token should not be generated")
}

// SignIn should return an error when fetching user by email fails
func (suite *schemaResolverTestSuite) TestSignInWithFetchUserByEmailError() {
	userRepositoryServiceMock := new(mockutil.UserRepositoryServiceMock)
	userRepositoryServiceMock.On("FetchByEmail", mock.Anything, mock.Anything, []string(nil)).Return(
		errors.New(mockutil.MockedGenericErrorMessage),
	).Times(1)
	suite.resolver.userRepository = userRepositoryServiceMock
	input := model.UserSignIn{Email: mockutil.DefaultEmail, Password: mockutil.DefaultPassword}

	token, err := suite.mutationResolver.SignIn(context.Background(), input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not sign in"),
		"Should return expected error when fetch user by email fails",
	)
	assert.Nil(suite.T(), token, "Token should not be generated")
}

// SignIn should return expected error when user gives wrong password
func (suite *schemaResolverTestSuite) TestSignInWithWrongPassword() {
	passwordSecurityServiceMock := new(mockutil.PasswordSecurityServiceMock)
	passwordSecurityServiceMock.On("HashWithArgon2id", mock.Anything).Return([]byte("WrongPassword")).Times(1)
	suite.resolver.passwordSecurityService = passwordSecurityServiceMock
	input := model.UserSignIn{Email: mockutil.DefaultEmail, Password: mockutil.DefaultPassword}

	token, err := suite.mutationResolver.SignIn(context.Background(), input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("wrong password"),
		"Should return expected error when user enters wrong password",
	)
	assert.Nil(suite.T(), token, "Token should not be generated")
}

// SignIn should return expected error when jwt generation fails
func (suite *schemaResolverTestSuite) TestSignInWithGenerateJwtError() {
	jwtAuthenticationServiceMock := new(mockutil.JwtAuthenticationServiceMock)
	jwtAuthenticationServiceMock.On("GenerateJwt", mock.Anything, mock.Anything).Return(
		"", errors.New(mockutil.MockedGenericErrorMessage),
	).Times(1)
	suite.resolver.authenticationService = jwtAuthenticationServiceMock
	input := model.UserSignIn{Email: mockutil.DefaultEmail, Password: mockutil.DefaultPassword}

	token, err := suite.mutationResolver.SignIn(context.Background(), input)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not sign in"),
		"Should return expected error when jwt generation fails",
	)
	assert.Nil(suite.T(), token, "Token should not be generated")
}

// QueryUserPasswords should successfully query for all user's passwords
func (suite *schemaResolverTestSuite) TestQueryUserPasswords() {
	passwordSecurityServiceMock := new(mockutil.PasswordSecurityServiceMock)
	passwordSecurityServiceMock.On("DecryptWithAes", mock.Anything, mock.Anything).Return("DecryptedPasswordMock", nil).Times(2)
	suite.resolver.passwordSecurityService = passwordSecurityServiceMock

	passwords, err := suite.queryResolver.QueryUserPasswords(suite.graphqlRequestContext, mockutil.DefaultIdAsString)
	assert.Nil(suite.T(), err, "Should fetch passwords without errors")

	assert.Equal(suite.T(), len(passwords), 2, "Query should fetch exactly two passwords")

	assert.Equal(suite.T(), passwords[0].Name, "Domain1")
	assert.Equal(suite.T(), passwords[1].Name, "Domain2")
	assert.Equal(suite.T(), passwords[0].Password, mockutil.MockedDecryptedPassword)
	assert.Equal(suite.T(), passwords[1].Password, mockutil.MockedDecryptedPassword)
	for _, password := range passwords {
		assert.Equal(suite.T(), password.UserID, mockutil.DefaultIdAsString)
	}
}

// QueryUserPasswords should should return empty slice when user has got no passwords
func (suite *schemaResolverTestSuite) TestQueryUserPasswordsWithoutUserPasswords() {
	passwords, err := suite.queryResolver.QueryUserPasswords(suite.graphqlRequestContext, "2")
	assert.Nil(suite.T(), err, "Should fetch passwords without errors")
	assert.Nil(suite.T(), passwords, "Should return nil passwords slice")
}

// QueryUserPasswords should should return expected error when userId is of an unexpected value
func (suite *schemaResolverTestSuite) TestQueryUserPasswordsWithUnexpectedUserIdValue() {
	passwords, err := suite.queryResolver.QueryUserPasswords(context.Background(), "invalidUserId")
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not fetch user's passwords"),
		"Should return expected error when user id is of an unexpected value",
	)
	assert.Nil(suite.T(), passwords, "Should not return any passwords data")
}

// QueryUserPasswords should return expected error on unsuccessful user's master password fetch
func (suite *schemaResolverTestSuite) TestQueryUserPasswordsWithMasterPasswordFetchError() {
	userRepositoryServiceMock := new(mockutil.UserRepositoryServiceMock)
	userRepositoryServiceMock.On("FetchMasterPasswordByUserId", mock.Anything, mock.Anything).Return(
		errors.New(mockutil.MockedGenericErrorMessage),
	).Times(1)
	suite.resolver.userRepository = userRepositoryServiceMock

	passwords, err := suite.queryResolver.QueryUserPasswords(context.Background(), mockutil.DefaultIdAsString)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not fetch user's passwords"),
		"Should return expected error when user master password fetch fails",
	)
	assert.Nil(suite.T(), passwords, "Should not return any passwords data")
}

// QueryUserPasswords should return expected error on unsuccessful user's passwords fetch
func (suite *schemaResolverTestSuite) TestQueryUserPasswordsWithFetchError() {
	passwordRepositoryServiceMock := new(mockutil.PasswordRepositoryServiceMock)
	passwordRepositoryServiceMock.On("FetchAllByUserId", mock.Anything, mock.Anything, mock.Anything).Return(
		errors.New(mockutil.MockedGenericErrorMessage),
	).Times(1)
	suite.resolver.passwordRepository = passwordRepositoryServiceMock

	passwords, err := suite.queryResolver.QueryUserPasswords(suite.graphqlRequestContext, mockutil.DefaultIdAsString)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not fetch user's passwords"),
		"Should return expected error when user's passwords fetch fails",
	)
	assert.Nil(suite.T(), passwords, "Should not return any password data")
}

// QueryUserPasswords should return expected error on unsuccessful user's passwords decryption
func (suite *schemaResolverTestSuite) TestQueryUserPasswordsWithDecryptionError() {
	passwordSecurityServiceMock := new(mockutil.PasswordSecurityServiceMock)
	passwordSecurityServiceMock.On("DecryptWithAes", mock.Anything, mock.Anything).Return(
		"", errors.New(mockutil.MockedGenericErrorMessage),
	).Times(1)
	suite.resolver.passwordSecurityService = passwordSecurityServiceMock

	passwords, err := suite.queryResolver.QueryUserPasswords(suite.graphqlRequestContext, mockutil.DefaultIdAsString)
	assert.Equal(
		suite.T(), err, gqlerror.Errorf("could not fetch user's passwords"),
		"Should return expected error when user's password decryption fails",
	)
	assert.Nil(suite.T(), passwords, "Should not return any user data")
}

func injectDefaultMockedResolverServices(suite *schemaResolverTestSuite) {
	resolver := NewResolver(
		mockutil.DefaultUserRepositoryServiceMock(),
		mockutil.DefaultPasswordRepositoryServiceMock(),
		mockutil.DefaultPasswordSecurityServiceMock(),
		mockutil.DefaultJwtAuthenticationServiceMock(),
	)
	suite.resolver = *resolver

	suite.mutationResolver = suite.resolver.Mutation()
	suite.queryResolver = suite.resolver.Query()
}
