package server

import (
	"context"
	"fmt"
	"net/http"

	"jamdrop/db"
	"jamdrop/model"

	"github.com/pkg/errors"
)

const (
	sessionName      = "session"
	sessionTokenName = "session_token"
)

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
	session, _ := s.Sessions.Get(r, sessionName)
	session.Values[sessionTokenName] = token
	if err := session.Save(r, w); err != nil {
		s.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
}

func (s *Server) LogOut(w http.ResponseWriter, r *http.Request) bool {
	s.Logger.Println("removing cookie")
	session, _ := s.Sessions.Get(r, sessionName)
	delete(session.Values, sessionTokenName)
	if err := session.Save(r, w); err != nil {
		s.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
}

func (s *Server) RequireLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Sessions.Get(r, sessionName)
		token := session.Values[sessionTokenName]
		if token == nil {
			s.Logger.Println("no session token")
			s.SpotifyAuthorizeRedirect(w, r)
			return
		}

		user, err := s.DB.GetUserBySessionToken(r.Context(), token.(string))
		if err != nil {
			if db.IsNotFound(err) {
				s.Logger.Println("no user for session_token", token)
				s.SpotifyAuthorizeRedirect(w, r)
				return
			}

			s.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		go func() {
			if err := s.DB.Touch(context.Background(), &model.SessionToken{Base: model.Base{ID: token.(string)}}); err != nil {
				s.Logger.Printf("failed to touch session_token; id: %v: %v\n", token, err)
			}
		}()

		next.ServeHTTP(w, r.WithContext(user.Context(r.Context())))
	})
}
