package server

import (
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

		songIdentifier := r.URL.Query().Get("song_identifier")
		songNameChan := make(chan string)
		go func() {
			defer close(songNameChan)

			songData, err := s.Spotify.GetSongData(user, songIdentifier)
			if err != nil {
				s.Logger.Println(err)
				songNameChan <- "Song"
				return
			}

			songNameChan <- fmt.Sprintf(`"%s"`, songData.Name)
		}()

		if err := s.Spotify.QueueSong(friend, songIdentifier); err != nil {
			s.Error(w, err, http.StatusInternalServerError, failureMessage)
			return
		}

		//songName := <-songNameChan
		//go func() {
		//	if err := s.Twilio.SongQueued(user, songName); err != nil {
		//		s.Logger.Println(err)
		//	}
		//}()

		s.Message(w, http.StatusCreated, "%s dropped to %s's queue", <-songNameChan, friend.DisplayName)
	}
}
