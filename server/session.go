package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/broothie/queuecumber/db"
	"github.com/broothie/queuecumber/model"
	"github.com/pkg/errors"
)

const sessionName = "session"
const sessionTokenName = "session_token"

func (s *Server) LogInUser(ctx context.Context, w http.ResponseWriter, r *http.Request, user *model.User) error {
	token, err := s.DB.CreateUserSession(ctx, user.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	if !s.LogIn(w, r, token) {
		return fmt.Errorf("failed to log user in with token '%s'", token)
	}

	return nil
}

func (s *Server) LogIn(w http.ResponseWriter, r *http.Request, token string) bool {
	s.Logger.Println("setting cookie", token)
	session, _ := s.Session.Get(r, sessionName)
	session.Values[sessionTokenName] = token
	if err := session.Save(r, w); err != nil {
		s.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
}

func (s *Server) LogOut(w http.ResponseWriter) {
	s.Logger.Println("removing cookie")
	http.SetCookie(w, &http.Cookie{
		Name:     sessionName,
		Expires:  time.Now(),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (s *Server) RequireLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return s.RequireLoggedInMiddleware(next)
}

func (s *Server) RequireLoggedInMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Session.Get(r, sessionName)
		token := session.Values[sessionTokenName]
		if token == nil {
			s.Logger.Println("no session token")
			http.Redirect(w, r, "/spotify/authorize", http.StatusTemporaryRedirect)
			return
		}

		user, err := s.DB.GetUserBySessionToken(r.Context(), token.(string))
		if err != nil {
			if db.IsNotFound(err) {
				s.Logger.Println("no user for session_token", token)
				http.Redirect(w, r, "/spotify/authorize", http.StatusTemporaryRedirect)
			} else {
				s.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		next.ServeHTTP(w, r.WithContext(user.Context(r.Context())))
	}
}

func (s *Server) GetFlashes(w http.ResponseWriter, r *http.Request) []interface{} {
	session, _ := s.Session.Get(r, sessionName)
	flashes := session.Flashes()
	if err := session.Save(r, w); err != nil {
		s.Logger.Println("failed to save flashes:", err)
	}

	return flashes
}

func (s *Server) Flash(w http.ResponseWriter, r *http.Request, value interface{}, vars ...string) {
	session, _ := s.Session.Get(r, sessionName)
	session.AddFlash(value, vars...)
	if err := session.Save(r, w); err != nil {
		s.Logger.Println("failed to save flashes:", err)
	}
}
