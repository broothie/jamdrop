package db

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/broothie/queuecumber/model"
	"github.com/pkg/errors"
)

type Model interface {
	GetID() string
	SetID(string)

	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)

	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time)

	Collection() model.Collection
}

func (db *DB) Create(ctx context.Context, m Model) error {
	now := time.Now()
	collection := db.EnvCollectionForModel(m)
	db.Logger.Println("db.Create", collection)

	m.SetCreatedAt(now)
	m.SetUpdatedAt(now)

	if m.GetID() == "" {
		result, _, err := db.CollectionForModel(m).Add(ctx, m)
		if err != nil {
			return errors.Wrapf(err, "failed to create record; collection: %s", collection)
		}

		m.SetID(result.ID)
	} else {
		if _, err := db.CollectionForModel(m).Doc(m.GetID()).Create(ctx, m); err != nil {
			return errors.Wrapf(err, "failed to create record; collection: %s, id: %s", collection, m.GetID())
		}
	}

	return nil
}

func (db *DB) Touch(ctx context.Context, m Model) error {
	collection := db.EnvCollectionForModel(m)
	db.Logger.Println("db.Touch", collection, m.GetID())
	return db.Update(ctx, m)
}

func (db *DB) Update(ctx context.Context, m Model, updates ...firestore.Update) error {
	collection := db.EnvCollectionForModel(m)
	db.Logger.Println("db.Update", collection, m.GetID())

	updatedAt := time.Now()
	updates = append(updates, firestore.Update{Path: "updated_at", Value: updatedAt})
	if _, err := db.CollectionForModel(m).Doc(m.GetID()).Update(ctx, updates); err != nil {
		return errors.Wrapf(err, "failed to update record; collection: %s, id: %s", collection, m.GetID())
	}

	m.SetUpdatedAt(updatedAt)
	return nil
}

func (db *DB) Get(ctx context.Context, id string, m Model) error {
	collection := db.EnvCollectionForModel(m)
	db.Logger.Println("db.Get", collection, id)

	doc, err := db.CollectionForModel(m).Doc(id).Get(ctx)
	if err != nil {
		if isCodeNotFound(err) {
			return db.notFound(m.Collection(), id)
		}

		return errors.Wrapf(err, "failed to get record; collection: %s, id: %s", db.EnvCollectionForModel(m), id)
	}

	if err := doc.DataTo(m); err != nil {
		return errors.Wrapf(err, "failed to read record data; collection: %s, id: %s", db.EnvCollectionForModel(m), id)
	}

	m.SetID(doc.Ref.ID)
	return nil
}

func (db *DB) Exists(ctx context.Context, collection model.Collection, id string) (bool, error) {
	envCollection := db.EnvCollectionForName(collection)
	db.Logger.Println("db.Exists", envCollection, id)

	_, err := db.Collection(envCollection).Doc(id).Get(ctx)
	if err != nil {
		if isCodeNotFound(err) {
			return false, nil
		}

		return false, errors.Wrapf(err, "failed to check record exists; collection: %s, id: %s", collection, id)
	}

	return true, nil
}

func (db *DB) CollectionForModel(m Model) *firestore.CollectionRef {
	return db.CollectionForName(m.Collection())
}

func (db *DB) CollectionForName(collection model.Collection) *firestore.CollectionRef {
	return db.Collection(db.EnvCollectionForName(collection))
}

func (db *DB) EnvCollectionForModel(m Model) string {
	return db.EnvCollectionForName(m.Collection())
}

func (db *DB) EnvCollectionForName(collection model.Collection) string {
	if db.Config.IsProduction() {
		return fmt.Sprintf("production.%s", collection)
	} else if db.Config.IsTest() {
		return fmt.Sprintf("test.%s", collection)
	} else {
		return fmt.Sprintf("development.%s", collection)
	}
}
