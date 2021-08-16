package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/KristijanFaust/gokeeper/app/authentication"
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/KristijanFaust/gokeeper/app/database/repository"
	"github.com/KristijanFaust/gokeeper/app/gql"
	"github.com/KristijanFaust/gokeeper/app/gql/generated"
	"github.com/KristijanFaust/gokeeper/app/security"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/upper/db/v4"
	"log"
	"net/http"
	"reflect"
	"sync"
	"syscall"
)

func Run(applicationConfig *config.Config, serverDoneWaitGroup *sync.WaitGroup, session *db.Session) *http.Server {
	if reflect.ValueOf(applicationConfig).IsZero() {
		log.Panic("Application configuration not loaded, cannot start server")
	}

	hostname := applicationConfig.Server.Hostname
	portNumber := applicationConfig.Server.Port
	log.Printf("Starting GoKeeper server on http://%s:%s", hostname, portNumber)

	router := chi.NewRouter()
	router.Use(authentication.AuthenticationMiddleware(applicationConfig.Authentication.JwtSigningKey))

	graphqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{Resolvers: gql.NewResolver(
			repository.NewUserRepositoryService(session),
			repository.NewPasswordRepositoryService(session),
			&security.PasswordSecurityService{
				Argon2PasswordHasher: &security.PasswordHashService{},
				AesPasswordCryptor:   &security.PasswordCryptoService{},
			},
			authentication.NewJwtAuthenticationService(applicationConfig.Authentication),
		)},
	))

	if reflect.ValueOf(applicationConfig.Profile).IsZero() || !applicationConfig.Profile.Production {
		router.Use(cors.New(cors.Options{
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowedHeaders: []string{"Authentication", "Content-Type"},
		}).Handler)

		router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
		log.Printf("Serving GraphQL playground on http://%s:%s/playground", hostname, portNumber)
	}
	router.Handle("/query", graphqlHandler)

	server := &http.Server{
		Addr:    hostname + ":" + portNumber,
		Handler: router,
	}

	go func(server *http.Server) {
		defer serverDoneWaitGroup.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Server error occurred: %s", err)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		} else {
			log.Printf("Received shutdown signal, terminating server")
		}
	}(server)

	return server
}
