package main

import (
	"github.com/broothie/queuecumber/app"
	"github.com/broothie/queuecumber/config"
	"github.com/broothie/queuecumber/server"
)

func main() {
	cfg := config.New()
	app := app.New(cfg)
	server.Run(app)
}
