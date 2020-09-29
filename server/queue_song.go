package server

import (
	"fmt"
	"net/http"

	"github.com/broothie/queuecumber/model"
)

func (s *Server) QueueSong() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.QueueSong")

		user, _ := model.UserFromContext(r.Context())
		friendID := r.URL.Query().Get("user_id")
		friend := new(model.User)
		if err := s.DB.Get(r.Context(), friendID, friend); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
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

		if err := s.Spotify.QueueSongForUser(friend, songIdentifier); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		s.Flash(w, r, fmt.Sprintf("%s added to %s's queue!", <-songNameChan, friend.DisplayName))
	}
}
