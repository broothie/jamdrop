package app

import (
	"log"
	"os"

	"github.com/broothie/queuecumber/config"
	"github.com/broothie/queuecumber/spotify"
)

type App struct {
	*config.Config
	Logger  *log.Logger
	Spotify *spotify.Client
}

func New(cfg *config.Config) *App {
	return &App{
		Config:  cfg,
		Logger:  log.New(os.Stdout, "[queuecumber] ", log.LstdFlags),
		Spotify: spotify.New(cfg),
	}
}
