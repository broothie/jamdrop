package model

const CollectionFollows Collection = "follows"

type Follow struct {
	Base
	FollowerID string `firestore:"follower_id"`
	FolloweeID string `firestore:"followee_id"`
}

func (*Follow) Collection() Collection {
	return CollectionFollows
}
