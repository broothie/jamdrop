package db

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/broothie/queuecumber/config"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type DB struct {
	*firestore.Client
	Config *config.Config
	Logger *log.Logger
}

func New(cfg *config.Config, logger *log.Logger) (*DB, error) {
	var options []option.ClientOption
	if !cfg.IsProduction() {
		options = append(options, option.WithCredentialsFile("queuecumber.json"))
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "queuecumber", options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create firestore client")
	}

	return &DB{Config: cfg, Logger: logger, Client: client}, nil
}
