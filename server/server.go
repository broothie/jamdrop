package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"jamdrop/app"
	"jamdrop/db"
	"jamdrop/job"
	"jamdrop/logger"
	"jamdrop/requestid"
	"jamdrop/spotify"
	"jamdrop/twilio"

	"github.com/gorilla/sessions"
)

type Server struct {
	App      *app.App
	Logger   *logger.Logger
	Spotify  *spotify.Client
	DB       *db.DB
	Sessions *sessions.CookieStore
	Twilio   *twilio.Twilio
}

func New(app *app.App) *Server {
	return &Server{
		App:      app,
		Logger:   app.Logger,
		Spotify:  app.Spotify,
		DB:       app.DB,
		Sessions: sessions.NewCookieStore([]byte(app.Config.SecretKey)),
		Twilio:   app.Twilio,
	}
}

func Run(app *app.App) {
	New(app).Run()
}

func (s *Server) Run() {
	if s.App.Config.IsDevelopment() {
		go func() {
			s.Logger.Info("starting user playing job ticker")
			ticker := time.NewTicker(time.Minute)

			for {
				if err := job.New(s.App).ScanUserPlayers(context.Background()); err != nil {
					s.Logger.Err(err, "failed to scan user players")
				}
				<-ticker.C
			}
		}()
	}

	s.Logger.Info(fmt.Sprintf("serving @ %s", s.App.Config.BaseURL()))
	s.Logger.Err(http.ListenAndServe(fmt.Sprintf(":%s", s.App.Config.Port), s.Handler()), "server panicked")
}

func (s *Server) Handler() http.Handler {
	return requestid.Middleware(ApplyLoggerMiddleware(s.Routes(), s.Logger))
}

func (s *Server) Error(w http.ResponseWriter, err error, code int, format string, a ...interface{}) {
	s.Logger.Err(err, "server.Error")

	errorMessage := fmt.Sprintf(format, a...)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": errorMessage}); err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}
}

func (s *Server) Message(w http.ResponseWriter, code int, format string, a ...interface{}) {
	s.DumpJSON(w, code, map[string]string{"message": fmt.Sprintf(format, a...)})
}

func (s *Server) DumpJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.Error(w, err, http.StatusInternalServerError, err.Error())
	}
}

func (s *Server) ParseJSON(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		s.Error(w, err, http.StatusBadRequest, "")
		return false
	}

	return true
}
