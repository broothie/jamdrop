package server

import (
	"net/http"

	"jamdrop/job"
)

func (s *Server) EjectSessionTokens(w http.ResponseWriter, r *http.Request) {
	if err := job.New(s.App).EjectSessionTokens(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) ScanUserPlayers(w http.ResponseWriter, r *http.Request) {
	if err := job.New(s.App).ScanUserPlayers(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
