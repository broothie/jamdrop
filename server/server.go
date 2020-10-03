package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"jamdrop/app"
	"jamdrop/db"
	"jamdrop/spotify"

	"github.com/gorilla/sessions"
)

type Server struct {
	App      *app.App
	Logger   *log.Logger
	Spotify  *spotify.Client
	DB       *db.DB
	Sessions *sessions.CookieStore
}

func New(app *app.App) *Server {
	return &Server{
		App:      app,
		Logger:   app.Logger,
		Spotify:  app.Spotify,
		DB:       app.DB,
		Sessions: sessions.NewCookieStore([]byte(app.Config.SecretKey)),
	}
}

func Run(app *app.App) {
	New(app).Run()
}

func (s *Server) Run() {
	s.Logger.Printf("serving @ %s", s.App.Config.BaseURL())
	s.Logger.Panic(http.ListenAndServe(fmt.Sprintf(":%s", s.App.Config.Port), s.Handler()))
}

func (s *Server) Handler() http.Handler {
	return ApplyLoggerMiddleware(s.Routes(), s.Logger)
}

func (s *Server) Error(w http.ResponseWriter, err error, code int, format string, a ...interface{}) {
	s.Logger.Println(err)

	error := fmt.Sprintf(format, a...)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": error}); err != nil {
		http.Error(w, error, http.StatusInternalServerError)
	}
}

func (s *Server) Message(w http.ResponseWriter, code int, format string, a ...interface{}) {
	s.JSON(w, code, map[string]string{"message": fmt.Sprintf(format, a...)})
}

func (s *Server) JSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.Error(w, err, http.StatusInternalServerError, err.Error())
	}
}
