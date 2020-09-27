package model

type SessionToken struct {
	UserID string `firestore:"user_id"`
}
