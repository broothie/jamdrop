package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"jamdrop/model"

	"github.com/gorilla/mux"
)

func (s *Server) QueueSong() http.HandlerFunc {
	const failureMessage = "There was a problem queueing the requested song"

	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.QueueSong")

		user, _ := model.UserFromContext(r.Context())
		targetUserID := mux.Vars(r)["user_id"]
		targetUser := new(model.User)
		if err := s.DB.Get(r.Context(), targetUserID, targetUser); err != nil {
			s.Error(w, err, http.StatusInternalServerError, failureMessage)
			return
		}

		if !user.CanDropTo(targetUser) {
			message := fmt.Sprintf("%s is not currently active", targetUser.DisplayName)
			s.Error(w, errors.New(message), http.StatusUnauthorized, message)
			return
		}

		// TODO: Possibly move this song name stuff into the Spotify service?
		songIdentifier := r.URL.Query().Get("song_identifier")
		songNameChan := make(chan string)
		go func() {
			defer close(songNameChan)

			songData, err := s.Spotify.GetSongData(user, songIdentifier)
			if err != nil {
				s.Logger.Println("failed to get song data", err)
				songNameChan <- "Song"
				return
			}

			songNameChan <- songData.Name
		}()

		if err := s.Spotify.QueueSong(targetUser, songIdentifier); err != nil {
			s.Error(w, err, http.StatusInternalServerError, failureMessage)
			return
		}

		// TODO: Possibly move this song event stuff into the Spotify service?
		songName := <-songNameChan
		go func() {
			if targetUser.IsActive() {
				event := model.QueuedSongEvent{SongName: songName, UserName: user.DisplayName}
				if err := s.DB.AddSongQueuedEvent(context.Background(), targetUser, event); err != nil {
					s.Logger.Println("failed to add song queued event", err)
				}
			} else if targetUser.StayActive {
				if err := s.Twilio.SongQueued(user, targetUser, songName); err != nil {
					s.Logger.Printf("failed to send sms to user; user_id: %s; %v\n", user.ID, err)
				}
			}
		}()

		s.Message(w, http.StatusCreated, `"%s" dropped to %s's queue`, songName, targetUser.DisplayName)
	}
}
