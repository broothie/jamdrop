package spotify

import (
	"jamdrop/config"
	"jamdrop/db"
	"log"
)

type Client struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	Logger       *log.Logger
	DB           *db.DB
}

func New(cfg *config.Config, db *db.DB, logger *log.Logger) *Client {
	return &Client{
		ClientID:     cfg.SpotifyClientID,
		ClientSecret: cfg.SpotifyClientSecret,
		BaseURL:      cfg.BaseURL(),
		Logger:       logger,
		DB:           db,
	}
}
