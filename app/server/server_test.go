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
	applicationConfig := config.LoadConfiguration("../../config.yml")

	var server *http.Server
	serverDoneWaitGroup := &sync.WaitGroup{}
	serverDoneWaitGroup.Add(1)
	assert.NotPanics(t, func() { server = Run(applicationConfig, serverDoneWaitGroup, nil) }, "Server should run without panics")

	err := server.Shutdown(context.TODO())
	assert.Nil(t, err, "Server should shutdown without any errors")
	serverDoneWaitGroup.Wait()
}

// Run should panic if no configuration is loaded
func TestRunWithoutConfiguration(t *testing.T) {
	assert.PanicsWithValue(
		t, "Application configuration not loaded, cannot start server",
		func() { Run(nil, nil, nil) }, "Server boot should panic if no configuration is loaded",
	)
}

// Run should terminate gracefully on server error
func TestRunWithServerBootError(t *testing.T) {
	applicationConfig := config.LoadConfiguration("../../config.yml")
	applicationConfig.Server.Port = "invalid"

	serverDoneWaitGroup := &sync.WaitGroup{}
	serverDoneWaitGroup.Add(1)

	signal.Ignore(syscall.SIGINT)
	assert.NotPanics(t, func() { Run(applicationConfig, serverDoneWaitGroup, nil) }, "Server should try to boot without panics")

	serverDoneWaitGroup.Wait()
}
