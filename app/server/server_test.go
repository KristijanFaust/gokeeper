package server

import (
	"context"
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"testing"
)

// Run should boot, run and shutdown a server successfully without errors
func TestRun(t *testing.T) {
	config.ApplicationConfig = new(config.Config)
	config.ApplicationConfig.Server.Port = "8080"
	defer func() { config.ApplicationConfig = nil }()

	var server *http.Server
	serverDoneWaitGroup := &sync.WaitGroup{}
	serverDoneWaitGroup.Add(1)
	assert.NotPanics(t, func() { server = Run(serverDoneWaitGroup) }, "Server should run without panics")

	server.Shutdown(context.TODO())
	serverDoneWaitGroup.Wait()
}

// Run should panic if no configuration is loaded
func TestRunWithoutConfiguration(t *testing.T) {
	assert.PanicsWithValue(
		t, "Server configuration not loaded, cannot start server",
		func() { Run(nil) }, "Server boot should panic if no configuration is loaded",
	)
}

// Run should terminate gracefully on server error
func TestRunWithServerBootError(t *testing.T) {
	config.ApplicationConfig = new(config.Config)
	config.ApplicationConfig.Server.Port = "invalid"
	defer func() { config.ApplicationConfig = nil }()

	serverDoneWaitGroup := &sync.WaitGroup{}
	serverDoneWaitGroup.Add(1)

	signal.Ignore(syscall.SIGINT)
	assert.NotPanics(t, func() { Run(serverDoneWaitGroup) }, "Server should try to boot without panics")

	serverDoneWaitGroup.Wait()
}
