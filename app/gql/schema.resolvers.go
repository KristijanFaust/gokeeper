package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"crypto/subtle"
	"log"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	databaseModel "github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/KristijanFaust/gokeeper/app/gql/generated"
	"github.com/KristijanFaust/gokeeper/app/gql/model"
	"github.com/lib/pq"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) SignUp(ctx context.Context, input model.NewUser) (*model.User, error) {
	validationErrors := manageValidationsErrors(r.validator.Struct(input), ctx)
	if validationErrors != nil {
		return nil, gqlerror.Errorf("validation error/s on user input")
	}

	passwordHash := r.passwordSecurityService.HashWithArgon2id(input.Password)

	newUser := databaseModel.User{Email: input.Email, Username: input.Username, Password: passwordHash}
	insertResult, err := r.userRepository.InsertNewUser(&newUser)
	if err != nil {
		switch errorType := err.(type) {
		case *pq.Error:
			if errorType.Code == "23505" {
				return nil, gqlerror.Errorf(existingEmailErrorMessage)
			}
		default:
			return nil, gqlerror.Errorf(userCreationErrorMessage)
		}
	}

	insertedUser := &model.User{
		ID:       strconv.FormatUint(uint64(insertResult.ID().(int64)), 10),
		Email:    input.Email,
		Username: input.Username,
	}
	return insertedUser, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, input model.UserSignIn) (*model.UserWithToken, error) {
	fetchedUser := databaseModel.User{}
	err := r.userRepository.FetchByEmail(&fetchedUser, input.Email, nil)
	if err != nil {
		if strings.Contains(err.Error(), "upper: no more rows in this result set") {
			return nil, gqlerror.Errorf(queryNonExistingEmailErrorMessage)
		}
		return nil, gqlerror.Errorf(signInErrorMessage)
	}

	if subtle.ConstantTimeCompare(r.passwordSecurityService.HashWithArgon2id(input.Password), fetchedUser.Password) == 0 {
		return nil, gqlerror.Errorf(wrongPasswordErrorMessage)
	}

	jwt, err := r.authenticationService.GenerateJwt(fetchedUser.Id)
	if err != nil {
		return nil, gqlerror.Errorf(signInErrorMessage)
	}

	user := &model.User{ID: strconv.FormatUint(fetchedUser.Id, 10), Email: fetchedUser.Email, Username: fetchedUser.Username}

	return &model.UserWithToken{User: user, Token: jwt}, nil
}

func (r *mutationResolver) CreatePassword(ctx context.Context, input model.NewPassword) (*model.Password, error) {
	validationErrors := manageValidationsErrors(r.validator.Struct(input), ctx)
	if validationErrors != nil {
		return nil, gqlerror.Errorf("validation error/s on password input")
	}

	userId, err := strconv.ParseUint(input.UserID, 10, 64)
	if err != nil {
		log.Printf("Error occurred while converting user id to uint64: %s", err)
		return nil, gqlerror.Errorf(passwordCreationErrorMessage)
	}

	userAuthentication := r.authenticationService.GetAuthenticatedUserDataFromContext(ctx)
	if userAuthentication == nil || userAuthentication.UserId != userId {
		return nil, gqlerror.Errorf(passwordAuthenticationErrorMessage)
	}

	user := databaseModel.User{}
	err = r.userRepository.FetchMasterPasswordByUserId(&user, userId)
	if err != nil {
		log.Printf("Error while fetching user master password: %s", err)
		return nil, gqlerror.Errorf(passwordCreationErrorMessage)
	}

	encryptedPassword, err := r.passwordSecurityService.EncryptWithAes(input.Password, user.Password)
	if err != nil {
		log.Printf("Error while encrypting user password: %s", err)
		return nil, gqlerror.Errorf(passwordCreationErrorMessage)
	}

	newPassword := databaseModel.Password{UserId: userId, Name: input.Name, Password: encryptedPassword}

	insertResult, err := r.passwordRepository.InsertNewPassword(&newPassword)
	if err != nil {
		return nil, gqlerror.Errorf(passwordCreationErrorMessage)
	}

	insertedPassword := &model.Password{
		ID:       strconv.FormatUint(uint64(insertResult.ID().(int64)), 10),
		UserID:   input.UserID,
		Name:     input.Name,
		Password: input.Password,
	}
	return insertedPassword, nil
}

func (r *mutationResolver) UpdatePassword(ctx context.Context, input model.UpdatePassword) (*model.Password, error) {
	validationErrors := manageValidationsErrors(r.validator.Struct(input), ctx)
	if validationErrors != nil {
		return nil, gqlerror.Errorf("validation error/s on password input")
	}

	passwordId, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		log.Printf("Error occurred while converting password id to uint64: %s", err)
		return nil, gqlerror.Errorf(passwordUpdateErrorMessage)
	}

	userAuthentication := r.authenticationService.GetAuthenticatedUserDataFromContext(ctx)
	userPassword := &databaseModel.Password{}
	err = r.passwordRepository.FetchPasswordById(userPassword, passwordId)
	if err != nil {
		log.Printf("Error occurred while fetching user password by id: %s", err)
		return nil, gqlerror.Errorf(passwordUpdateErrorMessage)
	}
	if userAuthentication == nil || userPassword.UserId != userAuthentication.UserId {
		return nil, gqlerror.Errorf(passwordAuthenticationErrorMessage)
	}

	user := databaseModel.User{}
	err = r.userRepository.FetchMasterPasswordByUserId(&user, userAuthentication.UserId)
	if err != nil {
		log.Printf("Error while fetching user master password: %s", err)
		return nil, gqlerror.Errorf(passwordUpdateErrorMessage)
	}

	encryptedPassword, err := r.passwordSecurityService.EncryptWithAes(input.Password, user.Password)
	if err != nil {
		log.Printf("Error while encrypting user password: %s", err)
		return nil, gqlerror.Errorf(passwordUpdateErrorMessage)
	}

	err = r.passwordRepository.UpdatePasswordById(input.Name, encryptedPassword, passwordId)
	if err != nil {
		log.Printf("Error while updating user password: %s", err)
		return nil, gqlerror.Errorf(passwordUpdateErrorMessage)
	}

	return &model.Password{ID: input.ID, UserID: input.UserID, Name: input.Name, Password: input.Password}, nil
}

func (r *queryResolver) QueryUserPasswords(ctx context.Context, userID string) ([]*model.Password, error) {
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		log.Printf("Error occurred while converting user id to uint64: %s", err)
		return nil, gqlerror.Errorf(userPasswordsFetchErrorMessage)
	}

	userAuthentication := r.authenticationService.GetAuthenticatedUserDataFromContext(ctx)
	if userAuthentication == nil || userAuthentication.UserId != userId {
		return nil, gqlerror.Errorf(userPasswordsAuthenticationErrorMessage)
	}

	var passwords []*model.Password
	fetchedPasswords := databaseModel.Passwords{}

	user := databaseModel.User{}
	err = r.userRepository.FetchMasterPasswordByUserId(&user, userId)
	if err != nil {
		log.Printf("Error while fetching user master password: %s", err)
		return nil, gqlerror.Errorf(userPasswordsFetchErrorMessage)
	}

	err = r.passwordRepository.FetchAllByUserId(&fetchedPasswords, userId, graphql.CollectAllFields(ctx))
	if err != nil {
		return nil, gqlerror.Errorf(userPasswordsFetchErrorMessage)
	}

	for _, password := range fetchedPasswords {
		decryptedPassword, err := r.passwordSecurityService.DecryptWithAes(password.Password, user.Password)
		if err != nil {
			log.Printf("Error while decrypting user password: %s", err)
			return nil, gqlerror.Errorf(userPasswordsFetchErrorMessage)
		}
		passwords = append(
			passwords,
			&model.Password{
				ID:       strconv.FormatUint(password.Id, 10),
				UserID:   strconv.FormatUint(password.UserId, 10),
				Name:     password.Name,
				Password: decryptedPassword,
			},
		)
	}
	return passwords, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
