package gql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const (
	userCreationErrorMessage                = "could not create a new user"
	passwordCreationErrorMessage            = "could not create a new password"
	passwordUpdateErrorMessage              = "could not update password"
	passwordAuthenticationErrorMessage      = "unauthorized password input"
	userPasswordsFetchErrorMessage          = "could not fetch user's passwords"
	userPasswordsAuthenticationErrorMessage = "unauthorized passwords fetch"
	signInErrorMessage                      = "could not sign in"
	existingEmailErrorMessage               = "the e-mail address is already taken"
	queryNonExistingEmailErrorMessage       = "user doesn't exist"
	wrongPasswordErrorMessage               = "wrong password"
)

func manageValidationsErrors(validationErrors error, ctx context.Context) error {
	if validationErrors != nil {
		for _, err := range validationErrors.(validator.ValidationErrors) {
			graphql.AddError(ctx, gqlerror.Errorf("field '%s' with value '%s' violates constraint: %s", err.Field(), err.Value(), err.Tag()))
		}
	}

	return validationErrors
}
