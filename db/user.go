package db

import (
	"context"
	"fmt"

	"jamdrop/logger"
	"jamdrop/model"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

func (db *DB) AddShare(ctx context.Context, user *model.User, shareUserID string) error {
	db.Logger.Info("db.AddShare", logger.Fields{"user_id": user.ID, "share_user_id": shareUserID})

	shareUserExists, err := db.Exists(ctx, model.CollectionUsers, shareUserID)
	if err != nil {
		return err
	}

	if !shareUserExists {
		return db.notFound(model.CollectionUsers, shareUserID)
	}

	user.EnsureShares()
	user.Shares[shareUserID] = model.UserShare{Exists: true, Enabled: true}
	option := firestore.Merge(firestore.FieldPath{"shares"})
	if _, err := db.Collection(user).Doc(user.ID).Set(ctx, user, option); err != nil {
		return errors.Wrapf(err, "failed to add share; user_id: %s, share_user_id: %s", user.ID, shareUserID)
	}

	return nil
}

// Get users this user has shared their queue with
func (db *DB) GetUserShares(ctx context.Context, user *model.User) ([]*model.User, error) {
	db.Logger.Info("db.GetUserShares", logger.Fields{"user_id": user.ID})

	userCollection := db.Collection(model.CollectionUsers)
	shareDocRefs := make([]*firestore.DocumentRef, len(user.Shares))
	counter := 0
	for id := range user.Shares {
		shareDocRefs[counter] = userCollection.Doc(id)
		counter++
	}

	shareDocs, err := db.GetAll(ctx, shareDocRefs)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get followee data")
	}

	var shares []*model.User
	for _, doc := range shareDocs {
		share := new(model.User)
		if err := doc.DataTo(share); err != nil {
			db.Logger.Err(err, "failed to read followee data", logger.Fields{"data": doc.Data()})
		}

		shares = append(shares, share)
	}

	return shares, nil
}

func (db *DB) GetUserSharers(ctx context.Context, user *model.User) ([]*model.User, error) {
	db.Logger.Info("db.GetUserSharers", logger.Field("user_id", user.ID))

	docs, err := db.
		Collection(model.CollectionUsers).
		Where(fmt.Sprintf("shares.%s.exists", user.ID), "==", true).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, errors.Wrapf(err, "")
	}

	var sharers []*model.User
	for _, doc := range docs {
		sharer := new(model.User)
		if err := doc.DataTo(sharer); err != nil {
			db.Logger.Err(err, "failed to read sharer data", logger.Field("user_id", user.ID))
		}

		sharers = append(sharers, sharer)
	}

	return sharers, nil
}

func (db *DB) AddSongQueuedEvent(ctx context.Context, user *model.User, event model.QueuedSongEvent) error {
	db.Logger.Info("db.AddSongQueuedEvent", logger.Field("user_id", user.ID))

	user.QueuedSongEvents = append(user.QueuedSongEvents, event)
	if err := db.Update(ctx, user, firestore.Update{Path: "queued_song_events", Value: user.QueuedSongEvents}); err != nil {
		return errors.Wrapf(err, "failed to add queued song event; user_id: %s, event: %+v", user.ID, event)
	}

	return nil
}
