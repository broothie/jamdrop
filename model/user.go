package model

import (
	"context"
	"time"
)

const CollectionUsers Collection = "users"

type User struct {
	Base
	AccessToken          string    `json:"access_token" firestore:"access_token"`
	RefreshToken         string    `json:"refresh_token" firestore:"refresh_token"`
	ExpiresIn            int       `json:"expires_in" firestore:"-"`
	AccessTokenExpiresAt time.Time `firestore:"access_token_expires_at"`

	DisplayName string  `json:"display_name" firestore:"display_name"`
	Images      []Image `json:"images" firestore:"images"`
}

type Image struct {
	URL string `json:"url" firestore:"url"`
}

func (*User) Collection() Collection {
	return CollectionUsers
}

func (u *User) UpdateAccessTokenExpiration() {
	u.AccessTokenExpiresAt = time.Now().Add(time.Duration(u.ExpiresIn) * time.Second)
}

func (u *User) AccessTokenIsFresh() bool {
	return time.Now().Before(u.AccessTokenExpiresAt)
}

type userContextKey struct{}

var userContextK userContextKey

func (u *User) Context(parent context.Context) context.Context {
	return context.WithValue(parent, userContextK, u)
}

func UserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userContextK).(*User)
	return user, ok
}
