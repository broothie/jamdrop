package server

import (
	"net/http"

	"github.com/broothie/queuecumber/model"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func (s *Server) GetUser() http.HandlerFunc {
	type Payload struct {
		User    User   `json:"user"`
		Shares  []User `json:"shares"`
		Sharers []User `json:"sharers"`
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
			Shares:  publicUsers(shares),
			Sharers: publicUsers(sharers),
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

		s.JSON(w, publicUsers(sharers))
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

		s.JSON(w, publicUsers(shares))
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
