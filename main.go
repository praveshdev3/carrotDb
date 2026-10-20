package main

import (
	"flag"
	"log"

	"github.com/praveshdev3/carrotDb/config"
	"github.com/praveshdev3/carrotDb/server"
)

func setupflags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Host to listen on")
	flag.StringVar(&config.Port, "port", "5379", "Port to listen on")
	flag.Parse()
}

func main() {
	setupflags()
	log.Println("Configuring the carrot 🥕")
	server.RunSyncTCPServer()
}
