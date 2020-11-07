package server

import (
	"net/http"

	"jamdrop/db"
	"jamdrop/model"
	"jamdrop/spotify"
)

func (s *Server) Share() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Debug("server.Share")

		shareUserIdentifier := r.URL.Query().Get("user_identifier")
		shareUserID, err := spotify.IDFromIdentifier(shareUserIdentifier)
		if err != nil {
			s.Error(w, err, http.StatusInternalServerError, "Invalid user identifier")
			return
		}

		user := model.UserFromContext(r.Context())
		nameChan := make(chan string)
		go func() {
			defer close(nameChan)

			shareUser, err := s.Spotify.GetUserByID(user, shareUserID)
			if err != nil {
				s.Logger.Err(err, "failed to get share user by id")
				nameChan <- "user"
				return
			}

			nameChan <- shareUser.DisplayName
		}()

		if err := s.DB.AddShare(r.Context(), user, shareUserID); err != nil {
			if db.IsNotFound(err) {
				s.Error(w, err, http.StatusNotFound, "%s has not connected to jamdrop yet.", <-nameChan)
				return
			}

			s.Error(w, err, http.StatusInternalServerError, "Failed to share queue with %s", <-nameChan)
			return
		}

		s.Message(w, http.StatusCreated, "Queue shared with %s", <-nameChan)
	}
}
