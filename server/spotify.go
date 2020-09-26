package server

import "net/http"

func (s *Server) SpotifyAuthorizeCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
