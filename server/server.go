package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/broothie/queuecumber/db"

	"github.com/broothie/queuecumber/spotify"

	"github.com/broothie/queuecumber/app"
)

type Server struct {
	App     *app.App
	Logger  *log.Logger
	Spotify *spotify.Spotify
	DB      *db.DB
	Session *sessions.CookieStore
}

func New(app *app.App) *Server {
	return &Server{
		App:     app,
		Logger:  app.Logger,
		Spotify: app.Spotify,
		DB:      app.DB,
		Session: sessions.NewCookieStore([]byte(app.Config.SecretKey)),
	}
}

func Run(app *app.App) {
	New(app).Run()
}

func (s *Server) Run() {
	s.Logger.Printf("serving @ %s", s.App.Config.BaseURL())
	s.Logger.Panic(http.ListenAndServe(fmt.Sprintf(":%d", s.App.Config.Port), s.Handler()))
}

func (s *Server) Handler() http.Handler {
	return ApplyLoggerMiddleware(s.Routes(), s.Logger)
}

func (s *Server) Error(w http.ResponseWriter, error string, code int) {
	s.Logger.Println(error)
	http.Error(w, error, code)
}
