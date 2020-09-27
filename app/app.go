package app

import (
	"log"
	"os"

	"github.com/broothie/queuecumber/config"
	"github.com/broothie/queuecumber/db"
	"github.com/broothie/queuecumber/spotify"
	"github.com/pkg/errors"
)

type App struct {
	*config.Config
	Logger  *log.Logger
	Spotify *spotify.Spotify
	DB      *db.DB
}

func New(cfg *config.Config) (*App, error) {
	logger := log.New(os.Stdout, "[queuecumber] ", log.LstdFlags)

	db, err := db.New(cfg, logger)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &App{
		Config:  cfg,
		Logger:  logger,
		Spotify: spotify.New(cfg, logger),
		DB:      db,
	}, nil
}
