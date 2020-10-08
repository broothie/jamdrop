package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) Routes() http.Handler {
	root := mux.NewRouter()

	// Fileserver
	root.
		Methods(http.MethodGet).
		PathPrefix("/public").
		Handler(http.StripPrefix("/public", http.FileServer(http.Dir("dist"))))

	// Build info
	root.
		Methods(http.MethodGet).
		Path("/build_info").
		HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { s.DumpJSON(w, http.StatusOK, s.App.BuildInfo) })

	// Spotify auth endpoints
	spotify := root.PathPrefix("/spotify").Subrouter()
	spotify.
		Methods(http.MethodGet).
		Path("/authorize").
		Handler(s.SpotifyAuthorize())

	spotify.
		Methods(http.MethodGet).
		Path("/authorize/callback").
		Handler(s.SpotifyAuthorizeCallback())

	spotify.
		Methods(http.MethodGet).
		Path("/authorize/failure").
		HandlerFunc(s.SpotifyAuthorizeFailure)

	// API
	api := root.PathPrefix("/api").Subrouter()
	api.Use(s.RequireLoggedIn)

	api.
		Methods(http.MethodPost).
		Path("/share").
		Handler(s.Share())

	// Users
	users := api.PathPrefix("/users").Subrouter()
	users.
		Methods(http.MethodGet).
		Path("/me").
		Handler(s.GetUser())

	users.
		Methods(http.MethodGet).
		Path("/me/sharers").
		Handler(s.GetUserSharers())

	users.
		Methods(http.MethodGet).
		Path("/me/shares").
		Handler(s.GetUserShares())

	users.
		Methods(http.MethodPatch).
		Path("/me").
		Handler(s.UserUpdate())

	users.
		Methods(http.MethodGet).
		Path("/me/ping").
		Handler(s.PingUser())

	users.
		Methods(http.MethodPost).
		Path("/{user_id}/queue").
		Handler(s.QueueSong())

	users.
		Methods(http.MethodPatch).
		Path("/{user_id}/enabled").
		Handler(s.SetShareEnabled())

	// Job endpoints
	if s.App.Config.Internal {
		jobs := root.PathPrefix("/jobs").Subrouter()
		jobs.
			Methods(http.MethodGet).
			Path("/eject_session_tokens").
			HandlerFunc(s.EjectSessionTokens)

		jobs.
			Methods(http.MethodGet).
			Path("/scan_user_players").
			HandlerFunc(s.ScanUserPlayers)
	}

	// Root
	root.
		Methods(http.MethodGet).
		Path("/").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "dist/index.html") })

	printRoutes(root)
	return root
}

func printRoutes(router *mux.Router) {
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, _ := route.GetMethods()
		if len(methods) == 0 {
			return nil
		}

		path, _ := route.GetPathTemplate()
		fmt.Printf("%v %s\n", methods, path)
		return nil
	})
}
