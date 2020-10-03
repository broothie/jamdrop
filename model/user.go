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
	PhoneNumber          string               `firestore:"phone_number"`
	LastPlaying          time.Time            `firestore:"last_playing"`
	LastPing             time.Time            `firestore:"last_ping"`
	QueuedSongEvents     []QueuedSongEvent    `firestore:"queued_song_events"`
}

type UserShare struct {
	Exists  bool `firestore:"exists"`
	Enabled bool `firestore:"enabled"`
}

type Image struct {
	URL string `json:"url" firestore:"url"`
}

type QueuedSongEvent struct {
	SongName string `firestore:"song_name" json:"song_name"`
	UserName string `firestore:"user_name" json:"user_name"`
}

func (*User) Collection() Collection {
	return CollectionUsers
}

func (u *User) IsPlaying() bool {
	return time.Now().Add(-5 * time.Minute).Before(u.LastPlaying)
}

func (u *User) IsActive() bool {
	return time.Now().Add(-time.Minute).Before(u.LastPing)
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

func (u *User) SetShare(otherUserID string, share UserShare) {
	u.EnsureShares()
	u.Shares[otherUserID] = share
}

func (u *User) GetShareFor(otherUserID string) UserShare {
	u.EnsureShares()
	return u.Shares[otherUserID]
}

func (u *User) HasQueueSharedWith(other *User) bool {
	u.EnsureShares()
	return u.GetShareFor(other.ID).Exists
}

func (u *User) HasQueueShareFrom(other *User) bool {
	return other.HasQueueSharedWith(u)
}

func (u *User) ShareReciprocated(other *User) bool {
	return u.HasQueueSharedWith(other) && u.HasQueueShareFrom(other)
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
