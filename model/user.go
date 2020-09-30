package model

import (
	"context"
	"time"
)

const CollectionUsers Collection = "users"

type User struct {
	Base
	AccessToken          string               `firestore:"access_token" json:"access_token"`
	RefreshToken         string               `firestore:"refresh_token" json:"refresh_token"`
	ExpiresIn            int                  `firestore:"-" json:"expires_in"`
	AccessTokenExpiresAt time.Time            `firestore:"access_token_expires_at"`
	DisplayName          string               `firestore:"display_name" json:"display_name"`
	Images               []Image              `firestore:"images" json:"images"`
	Shares               map[string]UserShare `firestore:"shares"` // Users this user has shared their queue with
}

type UserShare struct {
	Enabled bool `firestore:"enabled"`
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

func (u *User) EnsureShares() {
	if u.Shares == nil {
		u.Shares = make(map[string]UserShare)
	}
}

func (u *User) HasQueueSharedWith(other *User) bool {
	u.EnsureShares()
	_, isSharedWith := u.Shares[other.ID]
	return isSharedWith
}

func (u *User) HasQueueShareFrom(other *User) bool {
	return other.HasQueueSharedWith(u)
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
