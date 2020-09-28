package model

const CollectionSessionTokens Collection = "session_tokens"

type SessionToken struct {
	Base
	UserID string `firestore:"user_id"`
}

func (*SessionToken) Collection() Collection {
	return CollectionSessionTokens
}
