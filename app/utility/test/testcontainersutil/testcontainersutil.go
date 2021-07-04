package testcontainersutil

import (
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var dockerComposeFilePaths []string
var dockerComposeExecutionIdentifier string

func DockerComposeUp() bool {
	setRelativePathToDockerComposeFile()
	response, err := pingDatabaseContainer()
	if err != nil && strings.Contains(err.Error(), "EOF") || err == nil && response.StatusCode == http.StatusOK {
		log.Println("Docker container already running, skipping test containers setup")
		return true
	}
	log.Println("Starting test containers")
	dockerComposeExecutionIdentifier = strings.ToLower(uuid.New().String())
	dockerCompose := testcontainers.NewLocalDockerCompose(dockerComposeFilePaths, dockerComposeExecutionIdentifier)
	dockerComposeExecutionError := dockerCompose.
		WithCommand([]string{"up", "-d"}).
		Invoke()
	err = dockerComposeExecutionError.Error
	if err != nil {
		log.Printf("Could not run docker-compose files: %v - %v", dockerComposeFilePaths, err)
		return false
	}

	return isDatabaseReady()
}

func DockerComposeDown() {
	defer func() { dockerComposeExecutionIdentifier = "" }()
	if dockerComposeExecutionIdentifier == "" {
		log.Println("Docker compose execution identifier has no value, skipping container termination")
		return
	}
	dockerCompose := testcontainers.NewLocalDockerCompose(dockerComposeFilePaths, dockerComposeExecutionIdentifier)
	dockerComposeExecutionError := dockerCompose.Down()
	err := dockerComposeExecutionError.Error
	if err != nil {
		log.Printf("Could not terminate docker-compose file: %v - %v\nconsider terminating them manually", dockerComposeFilePaths, err)
		return
	}
	log.Println("Test docker containers terminated")
}

func pingDatabaseContainer() (*http.Response, error) {
	response, err := http.Get("http://localhost:50000")
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

// For some reason I can't get the wait strategy to work with testcontainers-go
// so this function should wait for the database to be ready
func isDatabaseReady() bool {
	timeout := time.After(30 * time.Second)
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-timeout:
			log.Println("Could not communicate with database after an expected period, terminating containers")
			DockerComposeDown()
			return false
		case <-tick:
			response, err := pingDatabaseContainer()
			if err != nil && strings.Contains(err.Error(), "EOF") || err == nil && response.StatusCode == http.StatusOK {
				log.Println("Test database is ready to use")
				return true
			}
		}
	}
}

func setRelativePathToDockerComposeFile() {
	_, filename, _, _ := runtime.Caller(0)
	dockerComposeFilePaths = []string{filepath.Dir(filename) + "/../../../../support/docker/docker-compose-tests.yml"}
}
