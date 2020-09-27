package server

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/broothie/queuecumber/model"
)

func (s *Server) SpotifyAuthorize() http.Handler {
	return http.RedirectHandler(s.Spotify.UserAuthorizeURL(), http.StatusTemporaryRedirect)
}

func (s *Server) SpotifyAuthorizeCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer http.Redirect(w, r, "/app", http.StatusPermanentRedirect)
		s.Logger.Println("server.SpotifyAuthorizeCallback")

		code := r.URL.Query().Get("code")
		user := new(model.User)
		if err := s.Spotify.SetUserTokens(code, user); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		if err := s.Spotify.SetUserData(user.AccessToken, user); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		fields := firestore.Merge(
			firestore.FieldPath{"id"},
			firestore.FieldPath{"display_name"},
			firestore.FieldPath{"access_token"},
			firestore.FieldPath{"refresh_token"},
		)

		if err := s.DB.UpsertUser(r.Context(), user, fields); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		if err := s.LogInUser(r.Context(), w, r, user); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}
	}
}
