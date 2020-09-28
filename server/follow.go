package server

import (
	"fmt"
	"net/http"

	"github.com/broothie/queuecumber/db"
	"github.com/broothie/queuecumber/model"
	"github.com/broothie/queuecumber/spotify"
)

func (s *Server) Follow() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.Follow")
		defer s.AppRedirect(w, r)

		user, _ := model.UserFromContext(r.Context())
		friendIdentifier := r.FormValue("user_identifier")
		friendID, err := spotify.IDFromIdentifier(friendIdentifier)
		if err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, fmt.Sprintf("An error occurred. Please try again."))
			return
		}

		friend := new(model.User)
		if err := s.DB.Get(r.Context(), friendID, friend); err != nil {
			if db.IsNotFound(err) {
				friend, err := s.Spotify.GetUserByID(user, friendID)
				if err != nil {
					s.Logger.Println(err)
					s.Flash(w, r, fmt.Sprintf("An error occurred. Please try again."))
					return
				}

				s.Flash(w, r, fmt.Sprintf("%s has not signed up with queuecumber yet.", friend.DisplayName))
				return
			}

			s.Logger.Println(err)
			s.Flash(w, r, fmt.Sprintf("An error occurred. Please try again."))
			return
		}

		follow := &model.Follow{FollowerID: user.ID, FolloweeID: friendID}
		if err := s.DB.Create(r.Context(), follow); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		s.Flash(w, r, fmt.Sprintf("Now following %s", friend.DisplayName))
	}
}
