package server

import (
	"net/http"

	"jamdrop/model"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type SharedUser struct {
	User
	ShareReciprocated bool `json:"share_reciprocated"`
}

func (s *Server) GetUser() http.HandlerFunc {
	type Payload struct {
		User    User         `json:"user"`
		Shares  []SharedUser `json:"shares"`
		Sharers []SharedUser `json:"sharers"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := model.UserFromContext(r.Context())

		var sharers []*model.User
		sharersDone := make(chan struct{})
		go func() {
			defer close(sharersDone)

			users, err := s.DB.GetUserSharers(r.Context(), user)
			if err != nil {
				s.Logger.Println(err)
				return
			}

			sharers = users
		}()

		shares, err := s.DB.GetUserShares(r.Context(), user)
		if err != nil {
			s.Logger.Println(err)
		}

		<-sharersDone
		s.JSON(w, Payload{
			User:    publicUser(user),
			Shares:  sharedUsers(user, shares),
			Sharers: sharedUsers(user, sharers),
		})
	}
}

func (s *Server) GetUserSharers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := model.UserFromContext(r.Context())

		sharers, err := s.DB.GetUserSharers(r.Context(), user)
		if err != nil {
			s.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.JSON(w, sharedUsers(user, sharers))
	}
}

func (s *Server) GetUserShares() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := model.UserFromContext(r.Context())

		shares, err := s.DB.GetUserShares(r.Context(), user)
		if err != nil {
			s.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.JSON(w, sharedUsers(user, shares))
	}
}

func publicUser(user *model.User) User {
	var imageURL string
	if len(user.Images) > 0 {
		imageURL = user.Images[0].URL
	}

	return User{
		ID:       user.ID,
		Name:     user.DisplayName,
		ImageURL: imageURL,
	}
}

func publicUsers(users []*model.User) []User {
	publicUsers := make([]User, len(users))
	for i, user := range users {
		publicUsers[i] = publicUser(user)
	}

	return publicUsers
}

func sharedUsers(currentUser *model.User, users []*model.User) []SharedUser {
	sharedUsers := make([]SharedUser, len(users))
	for i, user := range users {
		sharedUsers[i] = SharedUser{
			User:              publicUser(user),
			ShareReciprocated: currentUser.ShareReciprocated(user),
		}
	}

	return sharedUsers
}
