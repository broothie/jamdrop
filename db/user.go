package db

import (
	"context"
	"fmt"

	"jamdrop/model"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

func (db *DB) AddShare(ctx context.Context, user *model.User, shareUserID string) error {
	db.Logger.Println("db.AddShare", user.ID, shareUserID)

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
	db.Logger.Println("db.GetUserShares", user.ID)

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
			db.Logger.Println("failed to read followee data", doc.Data())
		}

		shares = append(shares, share)
	}

	return shares, nil
}

func (db *DB) GetUserSharers(ctx context.Context, user *model.User) ([]*model.User, error) {
	db.Logger.Println("db.GetUserSharers", user.ID)

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
			db.Logger.Printf("failed to read sharer data; user_id: %s, error: %v\n", user.ID, err)
		}

		sharers = append(sharers, sharer)
	}

	return sharers, nil
}
