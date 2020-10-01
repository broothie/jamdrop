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
