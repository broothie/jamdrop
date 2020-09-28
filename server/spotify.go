package server

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/broothie/queuecumber/model"
)

func (s *Server) SpotifyAuthorizeRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/spotify/authorize", http.StatusTemporaryRedirect)
}

func (s *Server) SpotifyAuthorize() http.Handler {
	return http.RedirectHandler(s.Spotify.UserAuthorizeURL(), http.StatusTemporaryRedirect)
}

func (s *Server) SpotifyAuthorizeCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.SpotifyAuthorizeCallback")
		defer s.AppRedirect(w, r)

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

		userExists, err := s.DB.Exists(r.Context(), model.CollectionUsers, user.ID)
		if err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		if userExists {
			updates := []firestore.Update{
				{Path: "display_name", Value: user.DisplayName},
				{Path: "access_token", Value: user.AccessToken},
				{Path: "refresh_token", Value: user.RefreshToken},
			}

			if err := s.DB.Update(r.Context(), user, updates...); err != nil {
				s.Logger.Println(err)
				s.Flash(w, r, "An error occurred. Please try again.")
				return
			}
		} else {
			if err := s.DB.Create(r.Context(), user); err != nil {
				s.Logger.Println(err)
				s.Flash(w, r, "An error occurred. Please try again.")
				return
			}
		}

		if err := s.LogInUser(r.Context(), w, r, user); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}
	}
}
