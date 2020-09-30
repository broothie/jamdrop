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
		defer s.AppRedirect(w, r)

		shareUserIdentifier := r.FormValue("user_identifier")
		shareUserID, err := spotify.IDFromIdentifier(shareUserIdentifier)
		if err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, fmt.Sprintf("An error occurred. Please try again."))
			return
		}

		nameChan := make(chan string)
		go func() {
			defer close(nameChan)

			shareUser := new(model.User)
			if err := s.DB.Get(r.Context(), shareUserID, shareUser); err != nil {
				s.Logger.Println(err)
				nameChan <- "user"
				return
			}

			nameChan <- shareUser.DisplayName
		}()

		user, _ := model.UserFromContext(r.Context())
		if err := s.DB.AddShare(r.Context(), user, shareUserID); err != nil {
			if db.IsNotFound(err) {
				shareUser, err := s.Spotify.GetUserByID(user, shareUserID)
				if err != nil {
					s.Logger.Println(err)
					s.Flash(w, r, "An error occurred. Please try again.")
					return
				}

				s.Flash(w, r, fmt.Sprintf("%s has not connected to queuecumber yet.", shareUser.DisplayName))
				return
			}

			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		s.Flash(w, r, fmt.Sprintf("Shared queue with %s", <-nameChan))
	}
}
