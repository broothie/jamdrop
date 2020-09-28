package server

import (
	"net/http"

	"github.com/broothie/queuecumber/model"
)

func (s *Server) QueueSong() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.QueueSong")

		friendID := r.URL.Query().Get("user_id")
		friend := new(model.User)
		if err := s.DB.Get(r.Context(), friendID, friend); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		songIdentifier := r.URL.Query().Get("song_identifier")
		if err := s.Spotify.QueueSongForUser(friend, songIdentifier); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		s.Flash(w, r, "Song queued successfully!")
	}
}
