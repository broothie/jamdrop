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
		friendID := mux.Vars(r)["user_id"]
		if !user.GetShareFor(friendID).Enabled {
			s.Error(w, fmt.Errorf("shares not enabled by %s for %s", user.ID, friendID), http.StatusUnauthorized, failureMessage)
			return
		}

		friend := new(model.User)
		if err := s.DB.Get(r.Context(), friendID, friend); err != nil {
			s.Error(w, err, http.StatusInternalServerError, failureMessage)
			return
		}

		if !friend.IsPlaying() || !friend.IsActive() {
			message := fmt.Sprintf("%s is not currently active", user.DisplayName)
			s.Error(w, errors.New(message), http.StatusUnauthorized, message)
			return
		}

		// TODO: Possible move this song name stuff into the Spotify service?
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

		if err := s.Spotify.QueueSong(friend, songIdentifier); err != nil {
			s.Error(w, err, http.StatusInternalServerError, failureMessage)
			return
		}

		// TODO: Possible move this song event stuff into the Spotify service?
		songName := <-songNameChan
		go func() {
			event := model.QueuedSongEvent{SongName: songName, UserName: user.DisplayName}
			if err := s.DB.AddSongQueuedEvent(context.Background(), friend, event); err != nil {
				s.Logger.Println("failed to add song queued event", err)
			}
		}()

		s.Message(w, http.StatusCreated, `"%s" dropped to %s's queue`, songName, friend.DisplayName)
	}
}
