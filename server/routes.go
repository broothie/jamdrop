package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) Routes() http.Handler {
	root := mux.NewRouter()

	root.
		Methods(http.MethodGet).
		Path("/app").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "public/index.html")
		})

	root.
		HandleFunc("/spotify/authorize", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, s.Spotify.UserAuthorizeURL(), http.StatusTemporaryRedirect)
		})

	return root
}
