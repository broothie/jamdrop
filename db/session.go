package db

import (
	"context"

	"github.com/broothie/queuecumber/model"
	"github.com/pkg/errors"
)

func (db *DB) GetUserBySessionToken(ctx context.Context, token string) (*model.User, error) {
	doc, err := db.collection(model.CollectionSessionTokens).Doc(token).Get(ctx)
	if err != nil {
		if isCodeNotFound(err) {
			return nil, db.notFound(model.CollectionSessionTokens, token)
		}

		return nil, errors.Wrapf(err, "failed to find session_token '%s'", token)
	}

	userID, err := doc.DataAt("user_id")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve userID from session token '%s'", token)
	}

	user := new(model.User)
	if err := db.Get(ctx, userID.(string), user); err != nil {
		if IsNotFound(err) {
			return nil, err
		}

		return nil, errors.Wrapf(err, "failed to get user '%v'", userID)
	}

	return user, nil
}

func (db *DB) CreateUserSession(ctx context.Context, userID string) (string, error) {
	sessionToken := &model.SessionToken{UserID: userID}
	if err := db.Create(ctx, sessionToken); err != nil {
		return "", err
	}

	return sessionToken.ID, nil
}
