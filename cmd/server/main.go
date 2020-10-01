package main

import (
	"fmt"
	"log"

	"jamdrop/app"
	"jamdrop/config"
	"jamdrop/server"
)

func main() {
	cfg := config.New()
	fmt.Print(cfg)

	app, err := app.New(cfg)
	if err != nil {
		log.Panic(err)
		return
	}

	server.Run(app)
}
