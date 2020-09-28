package db

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/broothie/queuecumber/model"
	"github.com/pkg/errors"
)

func (db *DB) GetUserFollowers(ctx context.Context, user *model.User) ([]*model.User, error) {
	db.Logger.Println("db.GetUserFollowers", user.ID)

	followDocs, err := db.
		collection(model.CollectionFollows).
		Where("followee_id", "==", user.ID).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get follows; user_id: %s", user.ID)
	}

	userCollection := db.collection(model.CollectionUsers)
	var followerDocRefs []*firestore.DocumentRef
	for _, doc := range followDocs {
		followeeID, err := doc.DataAt("follower_id")
		if err != nil {
			db.Logger.Println("failed to read follow data", doc.Data())
		}

		followerDocRefs = append(followerDocRefs, userCollection.Doc(followeeID.(string)))
	}

	followerDocs, err := db.GetAll(ctx, followerDocRefs)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get followee data")
	}

	var followers []*model.User
	for _, doc := range followerDocs {
		followee := new(model.User)
		if err := doc.DataTo(followee); err != nil {
			db.Logger.Println("failed to read followee data", doc.Data())
		}

		followers = append(followers, followee)
	}

	return followers, nil
}
