package db

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/broothie/queuecumber/model"
	"github.com/pkg/errors"
)

func (db *DB) GetUserFollowees(ctx context.Context, user *model.User) ([]*model.User, error) {
	db.Logger.Println("db.GetUserFollowees", user.ID)

	followDocs, err := db.
		collection(model.CollectionFollows).
		Where("follower_id", "==", user.ID).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, errors.Wrapf(err, "")
	}

	userCollection := db.collection(model.CollectionUsers)
	var followeeDocRefs []*firestore.DocumentRef
	for _, doc := range followDocs {
		followeeID, err := doc.DataAt("followee_id")
		if err != nil {
			db.Logger.Println("failed to read follow data", doc.Data())
		}

		followeeDocRefs = append(followeeDocRefs, userCollection.Doc(followeeID.(string)))
	}

	followeeDocs, err := db.GetAll(ctx, followeeDocRefs)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get followee data")
	}

	var followees []*model.User
	for _, doc := range followeeDocs {
		followee := new(model.User)
		if err := doc.DataTo(followee); err != nil {
			db.Logger.Println("failed to read followee data", doc.Data())
		}

		followees = append(followees, followee)
	}

	return followees, nil
}
