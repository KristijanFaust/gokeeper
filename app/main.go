package main

import (
	"github.com/KristijanFaust/gokeeper/app/server"
	"github.com/KristijanFaust/gokeeper/app/utility/stdout"
)

func main() {
	stdout.PrintApplicationBanner()
	server.Run()
}
