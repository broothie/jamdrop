package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/securecookie"
	_ "github.com/joho/godotenv/autoload"
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
	Host string `json:"host"`
	Port string `json:"port"`

	// Spotify
	SpotifyClientID     string `json:"spotify_client_id"`
	SpotifyClientSecret string `json:"-"`
	ScanWorkers         int    `json:"scan_workers"`
}

func New() *Config {
	c := new(Config)
	c.BuildTime = BuildTime

	c.Environment = Environment(env("APP_ENV", "development"))
	c.Internal = env("INTERNAL", "false") == "true"
	c.SecretKey = env("SECRET_KEY", string(securecookie.GenerateRandomKey(32)))

	// HTTP
	c.Host = env("HOST", "localhost")
	c.Port = env("PORT", "3000")

	// Spotify
	c.SpotifyClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	c.SpotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	var err error
	if c.ScanWorkers, err = strconv.Atoi(env("SCAN_WORKERS", "3")); err != nil {
		panic(err)
	}

	return c
}

func (c *Config) String() string {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("failed to marshal config: %v\n", err)
	}

	return fmt.Sprintf("config: %s\n", bytes)
}

func (c *Config) Protocol() string {
	if c.IsProduction() || c.IsNgrok() {
		return "https"
	}

	return "http"
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

func (c *Config) IsNgrok() bool {
	return strings.Contains(c.Host, "ngrok")
}

func (c *Config) BaseURL() string {
	return fmt.Sprintf("%s://%s%s", c.Protocol(), c.Host, c.portString())
}

func (c *Config) portString() string {
	if c.IsProduction() || c.IsNgrok() {
		return ""
	}

	return fmt.Sprintf(":%s", c.Port)
}

func (c *Config) devProd(dev, prod interface{}) interface{} {
	if c.IsProduction() {
		return prod
	}

	return dev
}

func (c *Config) devNgrok(dev, ngrok interface{}) interface{} {
	if c.IsNgrok() {
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
