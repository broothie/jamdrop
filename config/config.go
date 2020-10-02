package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gorilla/securecookie"
)

var (
	BuildTime string
)

type BuildInfo struct {
	BuildTime string `json:"build_time"`
}

type Environment string

const (
	EnvTest        Environment = "test"
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
)

type Config struct {
	BuildInfo

	Environment Environment `json:"environment"`
	Internal    bool        `json:"internal"`
	SecretKey   string      `json:"-"`

	// HTTP
	IsNgrok  bool   `json:"is_ngrok"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`

	// Spotify
	SpotifyClientID     string `json:"spotify_client_id"`
	SpotifyClientSecret string `json:"-"`

	// Twilio
	TwilioAccountSID string `json:"twilio_account_sid"`
	TwilioAuthToken  string `json:"-"`
}

func New() *Config {
	c := new(Config)
	c.BuildTime = BuildTime

	c.Environment = Environment(env("APP_ENV", "development"))
	c.Internal = env("INTERNAL", "false") == "true"
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

	// Twilio
	c.TwilioAccountSID = os.Getenv("TWILIO_ACCOUNT_SID")
	c.TwilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")

	return c
}

func (c *Config) String() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("failed to marshal config: %v\n", err)
	}

	return fmt.Sprintf("config: %s\n", bytes)
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
