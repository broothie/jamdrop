package model

import (
	"fmt"
	"time"
)

const CollectionSessionTokens Collection = "session_tokens"

type SessionToken struct {
	Base
	UserID string `firestore:"user_id"`
}

func (*SessionToken) Collection() Collection {
	return CollectionSessionTokens
}

func (st *SessionToken) IsExpired() bool {
	return st.UpdatedAt.Before(time.Now().Add(-30 * 24 * time.Hour))
}

func (st *SessionToken) CheckExpired() error {
	if st.IsExpired() {
		return ExpiredSessionTokenError{Token: st.ID, UserID: st.UserID}
	}

	return nil
}

type ExpiredSessionTokenError struct {
	Token  string
	UserID string
}

func (e ExpiredSessionTokenError) Error() string {
	return fmt.Sprintf("session_token is expired; session_token: %s, user_id: %s", e.Token, e.UserID)
}

func IsExpiredSessionTokenError(err error) bool {
	_, isExpiredSessionTokenError := err.(ExpiredSessionTokenError)
	return isExpiredSessionTokenError
}
