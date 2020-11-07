package spotify

import (
	"jamdrop/config"
	"jamdrop/db"
	"jamdrop/logger"
)

type Client struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	Logger       *logger.Logger
	DB           *db.DB
}

func New(cfg *config.Config, db *db.DB, logger *logger.Logger) *Client {
	return &Client{
		ClientID:     cfg.SpotifyClientID,
		ClientSecret: cfg.SpotifyClientSecret,
		BaseURL:      cfg.BaseURL(),
		Logger:       logger,
		DB:           db,
	}
}
