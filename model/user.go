package model

import "context"

type userContextKey struct{}

var userContextK userContextKey

type User struct {
	ID           string `json:"id" firestore:"id"`
	AccessToken  string `json:"access_token" firestore:"access_token"`
	RefreshToken string `json:"refresh_token" firestore:"refresh_token"`

	DisplayName string             `json:"display_name" firestore:"display_name"`
	Friends     map[string]*Friend `json:"friends" firestore:"friends"`
}

type Friend struct {
	ID          string `json:"id" firestore:"id"`
	DisplayName string `json:"display_name" firestore:"display_name"`
}

func (u *User) AddFriend(friendUser *User) {
	if u.Friends == nil {
		u.Friends = make(map[string]*Friend)
	}

	u.Friends[friendUser.ID] = friendUser.ToFriend()
}

func (u *User) ToFriend() *Friend {
	return &Friend{ID: u.ID, DisplayName: u.DisplayName}
}

func (u *User) Context(parent context.Context) context.Context {
	return context.WithValue(parent, userContextK, u)
}

func UserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userContextK).(*User)
	return user, ok
}
