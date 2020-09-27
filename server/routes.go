package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) Routes() http.Handler {
	root := mux.NewRouter()
	root.NotFoundHandler = http.RedirectHandler("/app", http.StatusPermanentRedirect)
	root.
		Methods(http.MethodGet).
		Path("/app").
		Handler(s.Index())

	users := root.PathPrefix("/users").Subrouter()
	users.
		Methods(http.MethodGet).
		Path("/friends").
		Handler(s.AddFriend())

	users.
		Methods(http.MethodGet).
		Path("/friends/queue").
		Handler(s.QueueSong())

	spotify := root.PathPrefix("/spotify").Subrouter()
	spotify.
		Methods(http.MethodGet).
		Path("/authorize").
		Handler(s.SpotifyAuthorize())

	spotify.
		Methods(http.MethodGet).
		Path("/authorize/callback").
		HandlerFunc(s.SpotifyAuthorizeCallback())

	fmt.Println("routes")
	root.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, _ := route.GetMethods()
		if len(methods) == 0 {
			return nil
		}

		path, _ := route.GetPathTemplate()
		fmt.Printf("%v %s\n", methods, path)
		return nil
	})

	return root
}
