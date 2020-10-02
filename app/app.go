package app

import (
	"log"
	"os"

	"jamdrop/config"
	"jamdrop/db"
	"jamdrop/spotify"
	"jamdrop/twilio"

	"github.com/pkg/errors"
)

type App struct {
	*config.Config
	Logger  *log.Logger
	Spotify *spotify.Spotify
	DB      *db.DB
	Twilio  *twilio.Twilio
}

func New(cfg *config.Config) (*App, error) {
	logFlags := log.Lshortfile
	if cfg.IsDevelopment() {
		logFlags |= log.LstdFlags
	}

	logger := log.New(os.Stdout, "[jamdrop] ", logFlags)
	db, err := db.New(cfg, logger)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &App{
		Config:  cfg,
		Logger:  logger,
		Spotify: spotify.New(cfg, db, logger),
		DB:      db,
		Twilio:  twilio.New(cfg, logger),
	}, nil
}
