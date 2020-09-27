package db

import (
	"context"

	"github.com/broothie/queuecumber/model"
	"github.com/pkg/errors"
)

func (db *DB) GetUserBySessionToken(ctx context.Context, token string) (*model.User, error) {
	db.Logger.Println("db.GetUserBySessionToken", token)

	sessionTokenDoc, err := db.Collection(CollectionSessionTokens).Doc(token).Get(ctx)
	if err != nil {
		return nil, handleGetError(err, CollectionSessionTokens, token)
	}

	userID, err := sessionTokenDoc.DataAt("user_id")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve userID from session token '%s'", token)
	}

	userDoc, err := db.Collection(CollectionUsers).Doc(userID.(string)).Get(ctx)
	if err != nil {
		return nil, handleGetError(err, CollectionUsers, userID.(string))
	}

	user := new(model.User)
	if err := userDoc.DataTo(user); err != nil {
		return nil, errors.Wrapf(err, "failed to unserialize data for user '%s'", userID)
	}

	return user, nil
}

func (db *DB) CreateUserSession(ctx context.Context, userID string) (string, error) {
	db.Logger.Println("db.CreateUserSession", userID)

	doc, _, err := db.Collection(CollectionSessionTokens).Add(ctx, model.SessionToken{UserID: userID})
	if err != nil {
		return "", errors.Wrapf(err, "failed to create session token for user '%s'", userID)
	}

	return doc.ID, nil
}
