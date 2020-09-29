package server

import (
	"net/http"
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
		user, err := s.Spotify.UserFromAuthorizationCode(r.Context(), code)
		if err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		if err := s.LogInUser(r.Context(), w, r, user); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
		}
	}
}
