package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"crypto/subtle"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	databaseModel "github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/KristijanFaust/gokeeper/app/gql/generated"
	"github.com/KristijanFaust/gokeeper/app/gql/model"
	"github.com/lib/pq"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
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

func (r *mutationResolver) SignIn(ctx context.Context, input model.UserSignIn) (*model.Token, error) {
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

	expireAt := int(time.Now().Add(time.Minute * 15).Unix())
	token, err := r.authenticationService.GenerateJwt(fetchedUser.Id, int64(expireAt))
	if err != nil {
		return nil, gqlerror.Errorf(signInErrorMessage)
	}

	return &model.Token{Token: token, ExpireAt: expireAt}, nil
}

func (r *queryResolver) QueryUserByEmail(ctx context.Context, email string) (*model.User, error) {
	fetchedUser := databaseModel.User{}
	err := r.userRepository.FetchByEmail(&fetchedUser, email, graphql.CollectAllFields(ctx))
	if err != nil {
		if strings.Contains(err.Error(), "upper: no more rows in this result set") {
			return nil, gqlerror.Errorf(queryNonExistingEmailErrorMessage)
		}
		return nil, gqlerror.Errorf(userFetchErrorMessage)
	}

	user := &model.User{
		ID:       strconv.FormatUint(fetchedUser.Id, 10),
		Email:    fetchedUser.Email,
		Username: fetchedUser.Username,
	}
	return user, nil
}

func (r *queryResolver) QueryUserPasswords(ctx context.Context, userID string) ([]*model.Password, error) {
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		log.Printf("Error occurred while converting user id to uint64: %s", err)
		return nil, gqlerror.Errorf(userPasswordsFetchErrorMessage)
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
