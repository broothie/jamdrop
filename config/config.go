package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gorilla/securecookie"
)

type Environment string

const (
	EnvTest        Environment = "test"
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
)

type Config struct {
	Environment Environment
	SecretKey   string

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
	c.Environment = Environment(env("APP_ENV", "development"))
	c.SecretKey = env("SECRET_KEY", string(securecookie.GenerateRandomKey(32)))

	// HTTP
	c.IsNgrok = env("NGROK", "false") == "true"
	c.Protocol = c.devProd(c.devNgrok("http", "https"), "https").(string)
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

func (c *Config) IsTest() bool {
	return c.Environment == EnvTest
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == EnvDevelopment
}

func (c *Config) IsProduction() bool {
	return c.Environment == EnvProduction
}

func (c *Config) BaseURL() string {
	return fmt.Sprintf("%s://%s%s", c.Protocol, c.Host, c.devProd(c.devNgrok(fmt.Sprintf(":%d", c.Port), ""), ""))
}

func (c *Config) devProd(dev, prod interface{}) interface{} {
	if c.IsProduction() {
		return prod
	}

	return dev
}

func (c *Config) devNgrok(dev, ngrok interface{}) interface{} {
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
