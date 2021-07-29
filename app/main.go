package main

import (
	"context"
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/KristijanFaust/gokeeper/app/database"
	"github.com/KristijanFaust/gokeeper/app/server"
	"github.com/KristijanFaust/gokeeper/app/utility/stdout"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	stdout.PrintApplicationBanner()
	applicationConfig := config.LoadConfiguration("./config.yml")
	session := database.InitializeDatabaseConnection(applicationConfig.Datasource)
	defer database.CloseDatabaseConnection(session)

	serverDoneWaitGroup := &sync.WaitGroup{}
	serverDoneWaitGroup.Add(1)
	server := server.Run(applicationConfig, serverDoneWaitGroup, session)

	waitForQuitSignal()

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(context); err != nil {
		log.Panicf("Error during server shutdown: %s", err)
	}

	serverDoneWaitGroup.Wait()

	log.Println("Application terminated successfully")
}

func waitForQuitSignal() {
	quitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(quitSignalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	<-quitSignalChannel
}
