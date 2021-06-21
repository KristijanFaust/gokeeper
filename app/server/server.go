package server

import (
	"fmt"
	"log"
	"net/http"
)

const portNumber = ":8080"

func Run() {
	log.Printf("Starting GoKeeper server on port %s", portNumber)

	server := &http.Server{
		Addr: portNumber,
	}

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("Server error occurred: %s", err))
	}
}
