package server

import (
	"fmt"
	"net/http"

	"github.com/broothie/queuecumber/db"
	"github.com/broothie/queuecumber/model"
	"github.com/broothie/queuecumber/spotify"
)

func (s *Server) Share() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.Share")

		shareUserIdentifier := r.URL.Query().Get("user_identifier")
		shareUserID, err := spotify.IDFromIdentifier(shareUserIdentifier)
		if err != nil {
			s.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, _ := model.UserFromContext(r.Context())
		nameChan := make(chan string)
		go func() {
			defer close(nameChan)

			shareUser, err := s.Spotify.GetUserByID(user, shareUserID)
			if err != nil {
				s.Logger.Println(err)
				nameChan <- "user"
				return
			}

			nameChan <- shareUser.DisplayName
		}()

		if err := s.DB.AddShare(r.Context(), user, shareUserID); err != nil {
			if db.IsNotFound(err) {
				s.Error(w, fmt.Sprintf("%s has not connected to queuecumber yet.", <-nameChan), http.StatusNotFound)
				return
			}

			s.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		s.JSON(w, map[string]string{"message": fmt.Sprintf("Queue shared with %s", <-nameChan)})
	}
}
