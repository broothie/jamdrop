package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/broothie/queuecumber/spotify"

	"github.com/broothie/queuecumber/app"
)

type Server struct {
	App     *app.App
	Logger  *log.Logger
	Spotify *spotify.Client
}

func New(app *app.App) *Server {
	return &Server{
		App:     app,
		Logger:  app.Logger,
		Spotify: app.Spotify,
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
