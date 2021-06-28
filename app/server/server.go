package server

import (
	"github.com/KristijanFaust/gokeeper/app/config"
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

	portNumber := config.ApplicationConfig.Server.Port
	log.Printf("Starting GoKeeper server on port %s", portNumber)

	server := &http.Server{
		Addr: ":" + portNumber,
	}

	go func() {
		defer serverDoneWaitGroup.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Server error occurred: %s", err)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		} else {
			log.Printf("Received shutdown signal, terminating server")
		}
	}()

	return server
}
