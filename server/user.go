package server

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"

	"jamdrop/model"
)

type User struct {
	ID               string                  `json:"id"`
	Name             string                  `json:"name"`
	ImageURL         string                  `json:"image_url"`
	IsPlaying        bool                    `json:"is_playing"`
	IsActive         bool                    `json:"is_active"`
	SongQueuedEvents []model.QueuedSongEvent `json:"song_queued_events"`
}

type SharedUser struct {
	User
	ShareReciprocated bool `json:"share_reciprocated"`
	Enabled           bool `json:"enabled"`
}

func (s *Server) GetUser() http.HandlerFunc {
	type Payload struct {
		User    User         `json:"user"`
		Shares  []SharedUser `json:"shares"`
		Sharers []SharedUser `json:"sharers"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.GetUser")
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
		s.JSON(w, http.StatusOK, Payload{
			User:    publicUser(user),
			Shares:  sharedUsers(user, shares),
			Sharers: sharedUsers(user, sharers),
		})
	}
}

func (s *Server) GetUserSharers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.GetUserSharers")
		user, _ := model.UserFromContext(r.Context())

		sharers, err := s.DB.GetUserSharers(r.Context(), user)
		if err != nil {
			s.Error(w, err, http.StatusInternalServerError, "Failed to get queue sharers")
			return
		}

		s.JSON(w, http.StatusOK, sharedUsers(user, sharers))
	}
}

func (s *Server) GetUserShares() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.GetUserShares")
		user, _ := model.UserFromContext(r.Context())

		shares, err := s.DB.GetUserShares(r.Context(), user)
		if err != nil {
			s.Error(w, err, http.StatusInternalServerError, "Failed to get queue shares")
			return
		}

		s.JSON(w, http.StatusOK, sharedUsers(user, shares))
	}
}

func (s *Server) SetShareEnabled() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.SetShareEnabled")

		user, _ := model.UserFromContext(r.Context())
		shareID := mux.Vars(r)["user_id"]
		enabled := r.URL.Query().Get("enabled") == "true"

		update := firestore.Update{FieldPath: firestore.FieldPath{"shares", shareID, "enabled"}, Value: enabled}
		if err := s.DB.Update(r.Context(), user, update); err != nil {
			s.Error(w, err, http.StatusInternalServerError, "Failed to update queue share setting")
			return
		}

		enabledString := "enabled"
		if !enabled {
			enabledString = "disabled"
		}

		s.Message(w, http.StatusOK, "Sharing with %s %s", shareID, enabledString)
	}
}

func (s *Server) PingUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.PingUser")

		user, _ := model.UserFromContext(r.Context())
		updates := []firestore.Update{{Path: "last_ping", Value: time.Now()}}
		if _, err := s.DB.Collection(model.CollectionUsers).Doc(user.ID).Update(r.Context(), updates); err != nil {
			s.Error(w, err, http.StatusInternalServerError, "failed to ping user")
		}

		go func() {
			user.QueuedSongEvents = []model.QueuedSongEvent{}
			update := firestore.Update{Path: "queued_song_events", Value: user.QueuedSongEvents}
			if err := s.DB.Update(context.Background(), user, update); err != nil {
				s.Logger.Printf("failed to clear queued song events; user_id: %s; %v\n", user.ID, err)
			}
		}()

		s.JSON(w, http.StatusOK, publicUser(user))
	}
}

func publicUser(user *model.User) User {
	var imageURL string
	if len(user.Images) > 0 {
		imageURL = user.Images[0].URL
	}

	return User{
		ID:               user.ID,
		Name:             user.DisplayName,
		ImageURL:         imageURL,
		IsPlaying:        user.IsPlaying,
		IsActive:         user.IsActive(),
		SongQueuedEvents: user.QueuedSongEvents,
	}
}

func sharedUsers(currentUser *model.User, users []*model.User) []SharedUser {
	sharedUsers := make([]SharedUser, len(users))
	for i, user := range users {
		sharedUsers[i] = SharedUser{
			User:              publicUser(user),
			ShareReciprocated: currentUser.ShareReciprocated(user),
			Enabled:           currentUser.GetShareFor(user.ID).Enabled,
		}
	}

	return sharedUsers
}
