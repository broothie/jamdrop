package server

import (
	"fmt"
	"net/http"
)

func (s *Server) SpotifyAuthorizeRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/spotify/authorize", http.StatusTemporaryRedirect)
}

func (s *Server) SpotifyAuthorizeFailureRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/spotify/authorize", http.StatusTemporaryRedirect)
}

func (s *Server) SpotifyAuthorize() http.Handler {
	return http.RedirectHandler(s.Spotify.UserAuthorizeURL(), http.StatusTemporaryRedirect)
}

func (s *Server) SpotifyAuthorizeCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.SpotifyAuthorizeCallback")
		defer http.Redirect(w, r, "/", http.StatusPermanentRedirect)

		code := r.URL.Query().Get("code")
		user, err := s.Spotify.UserFromAuthorizationCode(r.Context(), code)
		if err != nil {
			s.Logger.Println(err)
			s.SpotifyAuthorizeFailureRedirect(w, r)
			return
		}

		if err := s.LogInUser(r.Context(), w, r, user); err != nil {
			s.Logger.Println(err)
			s.SpotifyAuthorizeFailureRedirect(w, r)
		}
	}
}

func (s *Server) SpotifyAuthorizeFailure(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := fmt.Fprint(w, "failed to authorize with spotify"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
