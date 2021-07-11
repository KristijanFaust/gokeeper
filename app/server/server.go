package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/KristijanFaust/gokeeper/app/config"
	graph2 "github.com/KristijanFaust/gokeeper/app/gql"
	generated2 "github.com/KristijanFaust/gokeeper/app/gql/generated"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"reflect"
	"sync"
	"syscall"
)

func Run(serverDoneWaitGroup *sync.WaitGroup) *http.Server {
	if config.ApplicationConfig == nil || reflect.ValueOf(config.ApplicationConfig.Server).IsZero() {
		log.Panic("Server configuration not loaded, cannot start server")
	}

	hostname := config.ApplicationConfig.Server.Hostname
	portNumber := config.ApplicationConfig.Server.Port
	log.Printf("Starting GoKeeper server on http://%s:%s", hostname, portNumber)

	router := chi.NewRouter()
	graphqlHandler := handler.NewDefaultServer(generated2.NewExecutableSchema(generated2.Config{Resolvers: &graph2.Resolver{}}))

	if reflect.ValueOf(config.ApplicationConfig.Profile).IsZero() || !config.ApplicationConfig.Profile.Production {
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
