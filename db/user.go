package db

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/broothie/queuecumber/model"
	"github.com/pkg/errors"
)

func (db *DB) UpsertUser(ctx context.Context, user *model.User, opts ...firestore.SetOption) error {
	db.Logger.Println("db.UpsertUser", user.ID)

	_, err := db.Collection(CollectionUsers).Doc(user.ID).Set(ctx, user, opts...)
	return errors.Wrapf(err, "failed to upsert user '%s'", user.ID)
}

func (db *DB) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	db.Logger.Println("db.GetUserByID", userID)

	doc, err := db.Collection(CollectionUsers).Doc(userID).Get(ctx)
	if err != nil {
		return nil, handleGetError(err, CollectionUsers, userID)
	}

	user := new(model.User)
	if err := doc.DataTo(user); err != nil {
		return nil, errors.Wrapf(err, "failed to unserialize data for user '%s'", userID)
	}

	return user, nil
}
