package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"jamdrop/config"
	"jamdrop/model"

	"cloud.google.com/go/firestore"
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
		options = append(options, option.WithCredentialsFile("gcloud-key.json"))
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "jamdrop", options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create firestore client")
	}

	return &DB{Config: cfg, Logger: logger, Client: client}, nil
}

type Model interface {
	Collectioner

	GetID() string
	SetID(string)

	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)

	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time)
}

func (db *DB) Create(ctx context.Context, m Model) error {
	now := time.Now()
	collection := db.fullCollectionName(m)
	db.Logger.Println("db.Create", collection)

	m.SetCreatedAt(now)
	m.SetUpdatedAt(now)

	if m.GetID() == "" {
		result, _, err := db.Collection(m).Add(ctx, m)
		if err != nil {
			return errors.Wrapf(err, "failed to create record; collection: %s", collection)
		}

		m.SetID(result.ID)
	} else {
		if _, err := db.Collection(m).Doc(m.GetID()).Create(ctx, m); err != nil {
			return errors.Wrapf(err, "failed to create record; collection: %s, id: %s", collection, m.GetID())
		}
	}

	return nil
}

func (db *DB) Touch(ctx context.Context, m Model) error {
	db.Logger.Println("db.Touch", db.fullCollectionName(m), m.GetID())
	return db.Update(ctx, m)
}

func (db *DB) Update(ctx context.Context, m Model, updates ...firestore.Update) error {
	collection := db.fullCollectionName(m)
	db.Logger.Println("db.Update", collection, m.GetID())

	updatedAt := time.Now()
	updates = append(updates, firestore.Update{Path: "updated_at", Value: updatedAt})
	if _, err := db.Collection(m).Doc(m.GetID()).Update(ctx, updates); err != nil {
		return errors.Wrapf(err, "failed to update record; collection: %s, id: %s", collection, m.GetID())
	}

	return db.Get(ctx, m.GetID(), m)
}

func (db *DB) Get(ctx context.Context, id string, m Model) error {
	collection := db.fullCollectionName(m)
	db.Logger.Println("db.Get", collection, id)

	doc, err := db.Collection(m).Doc(id).Get(ctx)
	if err != nil {
		if isCodeNotFound(err) {
			return db.notFound(m.Collection(), id)
		}

		return errors.Wrapf(err, "failed to get record; collection: %s, id: %s", collection, id)
	}

	if err := doc.DataTo(m); err != nil {
		return errors.Wrapf(err, "failed to read record data; collection: %s, id: %s", collection, id)
	}

	m.SetID(doc.Ref.ID)
	return nil
}

func (db *DB) Exists(ctx context.Context, collection model.Collection, id string) (bool, error) {
	db.Logger.Println("db.Exists", db.fullCollectionName(collection), id)

	_, err := db.Collection(collection).Doc(id).Get(ctx)
	if err != nil {
		if isCodeNotFound(err) {
			return false, nil
		}

		return false, errors.Wrapf(err, "failed to check record exists; collection: %s, id: %s", collection, id)
	}

	return true, nil
}

type Collectioner interface {
	Collection() model.Collection
}

func (db *DB) Collection(collectioner Collectioner) *firestore.CollectionRef {
	return db.Client.Collection(db.fullCollectionName(collectioner))
}

func (db *DB) fullCollectionName(collectioner Collectioner) string {
	if db.Config.IsProduction() {
		return fmt.Sprintf("production.%s", collectioner.Collection())
	} else if db.Config.IsTest() {
		return fmt.Sprintf("test.%s", collectioner.Collection())
	} else {
		return fmt.Sprintf("development.%s", collectioner.Collection())
	}
}
