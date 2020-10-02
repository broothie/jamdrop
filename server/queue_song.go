package server

import (
	"fmt"
	"net/http"

	"jamdrop/model"

	"github.com/gorilla/mux"
)

func (s *Server) QueueSong() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.QueueSong")

		user, _ := model.UserFromContext(r.Context())
		friendID := mux.Vars(r)["user_id"]
		friend := new(model.User)
		if err := s.DB.Get(r.Context(), friendID, friend); err != nil {
			s.Error(w, err.Error(), http.StatusInternalServerError)
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
			s.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//songName := <-songNameChan
		//go func() {
		//	if err := s.Twilio.SongQueued(user, songName); err != nil {
		//		s.Logger.Println(err)
		//	}
		//}()

		s.JSON(w, map[string]string{"message": fmt.Sprintf("%s dropped to %s's queue", <-songNameChan, friend.DisplayName)})
	}
}
