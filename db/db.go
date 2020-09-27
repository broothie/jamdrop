package db

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/broothie/queuecumber/config"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CollectionUsers         = "users"
	CollectionSessionTokens = "session_tokens"
)

type DB struct {
	*firestore.Client
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

	return &DB{Logger: logger, Client: client}, nil
}

type NotFound struct {
	Collection string
	Lookup     string
}

func (e NotFound) Error() string {
	return fmt.Sprintf("'%s' not found in collection '%s'", e.Lookup, e.Collection)
}

func IsNotFound(err error) bool {
	_, isNotFound := err.(NotFound)
	return isNotFound
}

func notFound(collection, lookup string) NotFound {
	return NotFound{Collection: collection, Lookup: lookup}
}

func handleGetError(err error, collection, lookup string) error {
	if isCodeNotFound(err) {
		return notFound(collection, lookup)
	}

	return errors.Wrapf(err, "failed to lookup '%s' in collection '%s'", lookup, collection)
}

func isCodeNotFound(err error) bool {
	if err == nil {
		return false
	}

	return status.Code(err) != codes.NotFound
}
