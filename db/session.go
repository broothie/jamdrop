package db

import (
	"context"

	"github.com/broothie/queuecumber/model"
	"github.com/pkg/errors"
)

func (db *DB) GetUserBySessionToken(ctx context.Context, token string) (*model.User, error) {
	sessionToken := new(model.SessionToken)
	if err := db.Get(ctx, token, sessionToken); err != nil {
		if IsNotFound(err) {
			return nil, err
		}

		return nil, errors.Wrapf(err, "failed to find session_token '%s'", token)
	}

	if err := sessionToken.CheckExpired(); err != nil {
		return nil, err
	}

	userID := sessionToken.UserID
	user := new(model.User)
	if err := db.Get(ctx, userID, user); err != nil {
		if IsNotFound(err) {
			return nil, err
		}

		return nil, errors.Wrapf(err, "failed to get user '%s'", userID)
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
