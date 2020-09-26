package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Environment string

	// HTTP
	IsNgrok  bool
	Protocol string
	Host     string
	Port     int

	// Spotify
	SpotifyClientID     string
	SpotifyClientSecret string
}

func New() *Config {
	c := &Config{}
	c.Environment = env("APP_ENV", "development")

	// HTTP
	c.IsNgrok = env("NGROK", "false") == "true"
	c.Protocol = c.devProd(c.devNgrok("http", "https"), "https")
	c.Host = env("HOST", "localhost")

	var err error
	if c.Port, err = strconv.Atoi(env("PORT", "3000")); err != nil {
		panic(err)
	}

	// Spotify
	c.SpotifyClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	c.SpotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")

	return c
}

func (c *Config) BaseURL() string {
	return fmt.Sprintf("%s://%s%s", c.Protocol, c.Host, c.devProd(c.devNgrok(fmt.Sprintf(":%d", c.Port), ""), ""))
}

func (c *Config) devProd(dev, prod string) string {
	if c.Environment == "production" {
		return prod
	}

	return dev
}

func (c *Config) devNgrok(dev, ngrok string) string {
	if c.IsNgrok {
		return ngrok
	}

	return dev
}

func env(name, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		value = fallback
	}

	return value
}
