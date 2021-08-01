package authentication

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type AuthenticationMiddlewareTestSuite struct {
	suite.Suite
	server            *httptest.Server
	client            *http.Client
	token             string
	defaultSigningKey string
}

func TestAuthenticationMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationMiddlewareTestSuite))
}

func (suite *AuthenticationMiddlewareTestSuite) SetupSuite() {
	suite.client = &http.Client{}
	suite.defaultSigningKey = "signingKey"
	suite.token = generateTestJwt(suite.defaultSigningKey, 1)
	suite.server = setUpTestServerWithAuthenticationMiddleware(suite.defaultSigningKey)
}

func (suite *AuthenticationMiddlewareTestSuite) TearDownSuite() {
	suite.server.Close()
}

// AuthenticationMiddleware should successfully put user authentication data in request context
func (suite *AuthenticationMiddlewareTestSuite) TestAuthenticationMiddleware() {
	request, _ := http.NewRequest("GET", suite.server.URL+"/", nil)
	request.Header.Set("Authentication", suite.token)
	response, _ := suite.client.Do(request)
	defer response.Body.Close()

	responseBody, _ := io.ReadAll(response.Body)

	assert.Equal(suite.T(), string(responseBody), "UserId: 1")
}

// AuthenticationMiddleware should successfully process requests that don't have an authentication value in the header
func (suite *AuthenticationMiddlewareTestSuite) TestAuthenticationMiddlewareWithoutAuthenticationHeader() {
	request, _ := http.NewRequest("GET", suite.server.URL+"/", nil)
	response, _ := suite.client.Do(request)
	defer response.Body.Close()

	responseBody, _ := io.ReadAll(response.Body)

	assert.Equal(suite.T(), string(responseBody), "No authentication header in client request")
}

// AuthenticationMiddleware should not put user authentication data in request context if error occurs while decoding JWT (includes token expiration)
func (suite *AuthenticationMiddlewareTestSuite) TestAuthenticationMiddlewareWithJwtDecodeError() {
	suite.token = generateTestJwt(suite.defaultSigningKey, -1)
	defer func(token string) { token = generateTestJwt(suite.defaultSigningKey, 1) }(suite.token)
	request, _ := http.NewRequest("GET", suite.server.URL+"/", nil)
	request.Header.Set("Authentication", suite.token)
	response, _ := suite.client.Do(request)
	defer response.Body.Close()

	responseBody, _ := io.ReadAll(response.Body)

	assert.Equal(suite.T(), string(responseBody), "No authentication header in client request")
}

func setUpTestServerWithAuthenticationMiddleware(jwtSigningKey string) *httptest.Server {
	router := chi.NewRouter()
	router.Use(AuthenticationMiddleware(jwtSigningKey))

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		if userAuthenticationData, ok := request.Context().Value(userContextKey).(userAuthentication); ok {
			writer.Write([]byte("UserId: " + strconv.FormatUint(userAuthenticationData.UserId, 10)))
		} else {
			writer.Write([]byte("No authentication header in client request"))
		}
	})

	return httptest.NewServer(router)
}

func generateTestJwt(signingKey string, minutesToExpire int) string {
	authenticationService := NewJwtAuthenticationService("issuer", []byte(signingKey))
	expireAt := time.Now().Add(time.Minute * time.Duration(minutesToExpire)).Unix()
	token, _ := authenticationService.GenerateJwt(uint64(1), expireAt)

	return token
}
