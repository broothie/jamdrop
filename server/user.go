package server

import (
	"context"
	"errors"
	"fmt"
	"jamdrop/logger"
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
	PhoneNumber      string                  `json:"phone_number"`
	IsPlaying        bool                    `json:"is_playing"`
	IsActive         bool                    `json:"is_active"`
	StayActive       bool                    `json:"stay_active"`
	SongQueuedEvents []model.QueuedSongEvent `json:"song_queued_events"`
}

type SharedUser struct {
	User
	ShareReciprocated bool `json:"share_reciprocated"`
	Enabled           bool `json:"enabled"`
	Droppable         bool `json:"droppable"`
}

func (s *Server) GetUser() http.HandlerFunc {
	type Payload struct {
		User    User         `json:"user"`
		Shares  []SharedUser `json:"shares"`
		Sharers []SharedUser `json:"sharers"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("server.GetUser")
		user := model.UserFromContext(r.Context())

		var sharers []*model.User
		sharersDone := make(chan struct{})
		go func() {
			defer close(sharersDone)

			users, err := s.DB.GetUserSharers(r.Context(), user)
			if err != nil {
				s.Logger.Err(err, "failed to get user sharers")
				return
			}

			sharers = users
		}()

		shares, err := s.DB.GetUserShares(r.Context(), user)
		if err != nil {
			s.Logger.Err(err, "failed to get user shares")
		}

		<-sharersDone
		s.DumpJSON(w, http.StatusOK, Payload{
			User:    publicUser(user),
			Shares:  sharedUsers(user, shares),
			Sharers: sharedUsers(user, sharers),
		})
	}
}

func (s *Server) GetUserSharers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("server.GetUserSharers")
		user := model.UserFromContext(r.Context())

		sharers, err := s.DB.GetUserSharers(r.Context(), user)
		if err != nil {
			s.Error(w, err, http.StatusInternalServerError, "Failed to get queue sharers")
			return
		}

		s.DumpJSON(w, http.StatusOK, sharedUsers(user, sharers))
	}
}

func (s *Server) GetUserShares() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("server.GetUserShares")

		user := model.UserFromContext(r.Context())
		shares, err := s.DB.GetUserShares(r.Context(), user)
		if err != nil {
			s.Error(w, err, http.StatusInternalServerError, "Failed to get queue shares")
			return
		}

		s.DumpJSON(w, http.StatusOK, sharedUsers(user, shares))
	}
}

func (s *Server) UserUpdate() http.HandlerFunc {
	fields := []string{"stay_active", "phone_number"}

	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("server.UserUpdate")

		var requestUpdates map[string]interface{}
		if !s.ParseJSON(w, r, &requestUpdates) {
			return
		}

		var dbUpdates []firestore.Update
		for _, field := range fields {
			value, entryExists := requestUpdates[field]
			if !entryExists {
				continue
			}

			dbUpdates = append(dbUpdates, firestore.Update{Path: field, Value: value})
		}

		if len(dbUpdates) == 0 {
			s.Error(w, errors.New("no valid updates provided"), http.StatusBadRequest, "No valid updates provided")
			return
		}

		user := model.UserFromContext(r.Context())
		if err := s.DB.Update(r.Context(), user, dbUpdates...); err != nil {
			s.Error(w, err, http.StatusInternalServerError, "Failed to update user")
			return
		}
	}
}

func (s *Server) SetShareEnabled() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("server.SetShareEnabled")

		user := model.UserFromContext(r.Context())
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
		s.Logger.Info("server.PingUser")

		user := model.UserFromContext(r.Context())
		updates := []firestore.Update{{Path: "last_ping", Value: time.Now()}}
		if _, err := s.DB.Collection(model.CollectionUsers).Doc(user.ID).Update(r.Context(), updates); err != nil {
			s.Error(w, err, http.StatusInternalServerError, "failed to ping user")
		}

		go func() {
			user.QueuedSongEvents = []model.QueuedSongEvent{}
			update := firestore.Update{Path: "queued_song_events", Value: user.QueuedSongEvents}
			if err := s.DB.Update(context.Background(), user, update); err != nil {
				s.Logger.Err(err, "failed to clear queued song events", logger.Field("user_id", user.ID))
			}
		}()

		s.DumpJSON(w, http.StatusOK, publicUser(user))
	}
}

func publicUser(user *model.User) User {
	imageURL := fmt.Sprintf("https://robohash.org/%s", user.ID)
	if len(user.Images) > 0 {
		imageURL = user.Images[0].URL
	}

	return User{
		ID:               user.ID,
		Name:             user.DisplayName,
		ImageURL:         imageURL,
		PhoneNumber:      user.PhoneNumber,
		IsActive:         user.IsActive(),
		StayActive:       user.StayActive,
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
			Droppable:         currentUser.CanDropTo(user),
		}
	}

	return sharedUsers
}
