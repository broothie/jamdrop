package main

import (
	"log"

	"github.com/broothie/queuecumber/app"
	"github.com/broothie/queuecumber/config"
	"github.com/broothie/queuecumber/server"
)

func main() {
	cfg := config.New()
	app, err := app.New(cfg)
	if err != nil {
		log.Panic(err)
		return
	}

	server.Run(app)
}
