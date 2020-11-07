package app

import (
	"log"

	"jamdrop/config"
	"jamdrop/db"
	"jamdrop/logger"
	"jamdrop/spotify"
	"jamdrop/twilio"

	"github.com/pkg/errors"
)

type App struct {
	*config.Config
	Logger  *logger.Logger
	Spotify *spotify.Client
	DB      *db.DB
	Twilio  *twilio.Twilio
}

func New(cfg *config.Config) (*App, error) {
	logFlags := log.Lshortfile
	if cfg.IsDevelopment() {
		logFlags |= log.LstdFlags
	}

	logLevel := logger.Info
	if cfg.IsDevelopment() {
		logLevel = logger.Debug
	}

	logger := logger.New(logger.ConfigureLevel(logLevel))
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
