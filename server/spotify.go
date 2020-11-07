package server

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"net/http"
)

func (s *Server) SpotifyAuthorizeRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/spotify/authorize", http.StatusTemporaryRedirect)
}

func (s *Server) SpotifyAuthorizeFailureRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/spotify/authorize/failure", http.StatusTemporaryRedirect)
}

func (s *Server) SpotifyAuthorize() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := s.randCode()
		session, _ := s.Sessions.Get(r, sessionName)
		session.Values["state"] = state
		if err := session.Save(r, w); err != nil {
			s.Logger.Err(err, "failed to save session")
			s.SpotifyAuthorizeRedirect(w, r)
			return
		}

		http.Redirect(w, r, s.Spotify.UserAuthorizeURL(state), http.StatusTemporaryRedirect)
	}
}

func (s *Server) SpotifyAuthorizeCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Debug("server.SpotifyAuthorizeCallback")
		defer http.Redirect(w, r, "/", http.StatusPermanentRedirect)

		urlState := r.URL.Query().Get("state")
		session, _ := s.Sessions.Get(r, sessionName)
		sessionState, ok := session.Values["state"].(string)
		if !ok {
			s.Logger.Info("failed to get state from session")
			s.SpotifyAuthorizeFailureRedirect(w, r)
			return
		}

		if urlState != sessionState {
			s.Logger.Info("states do not match")
			s.SpotifyAuthorizeFailureRedirect(w, r)
			return
		}

		code := r.URL.Query().Get("code")
		user, err := s.Spotify.UserFromAuthorizationCode(r.Context(), code)
		if err != nil {
			s.Logger.Err(err, "failed to get user from authorization code")
			s.SpotifyAuthorizeFailureRedirect(w, r)
			return
		}

		if err := s.LogInUser(r.Context(), w, r, user); err != nil {
			s.Logger.Err(err, "failed to log in user")
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

func (s *Server) randCode() string {
	const size = 8
	const alphabet = "abcdefghijlkmnopqrstuvwxyz"

	runes := make([]rune, size)
	for i := 0; i < size; i++ {
		index, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			s.Logger.Err(err, "randCode")
			runes[i] = rune(alphabet[mathrand.Intn(len(alphabet))])
			continue
		}

		runes[i] = rune(alphabet[index.Int64()])
	}

	return string(runes)
}
