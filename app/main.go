package main

import (
	"context"
	"github.com/KristijanFaust/gokeeper/app/config"
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
	config.LoadConfiguration("./config.yml")

	serverDoneWaitGroup := &sync.WaitGroup{}
	serverDoneWaitGroup.Add(1)
	server := server.Run(serverDoneWaitGroup)

	waitForQuitSignal()

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(context); err != nil {
		log.Panicf("Error during server shutdown: %s", err)
	}

	serverDoneWaitGroup.Wait()

	log.Printf("Application terminateted succesfully")
}

func waitForQuitSignal() {
	quitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(quitSignalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	<-quitSignalChannel
}
