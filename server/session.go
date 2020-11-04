package server

import (
	"context"
	"fmt"
	"jamdrop/logger"
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
	s.Logger.Info("setting cookie", logger.Field("token", token))
	session, _ := s.Sessions.Get(r, sessionName)
	session.Values[sessionTokenName] = token
	if err := session.Save(r, w); err != nil {
		s.Error(w, err, http.StatusInternalServerError, "failed to log in")
		return false
	}

	return true
}

func (s *Server) LogOut(w http.ResponseWriter, r *http.Request) bool {
	s.Logger.Info("removing cookie")
	session, _ := s.Sessions.Get(r, sessionName)
	delete(session.Values, sessionTokenName)
	if err := session.Save(r, w); err != nil {
		s.Error(w, err, http.StatusInternalServerError, "failed to log out")
		return false
	}

	return true
}

func (s *Server) RequireLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Sessions.Get(r, sessionName)
		token := session.Values[sessionTokenName]
		if token == nil {
			s.Logger.Info("no session token")
			s.SpotifyAuthorizeRedirect(w, r)
			return
		}

		user, err := s.DB.GetUserBySessionToken(r.Context(), token.(string))
		if err != nil {
			if db.IsNotFound(err) {
				s.Logger.Info("no user for session_token", logger.Field("token", token))
				s.SpotifyAuthorizeRedirect(w, r)
				return
			}

			if model.IsExpiredSessionTokenError(err) {
				s.Logger.Info("session_token expired", logger.Field("token", token))
				s.LogOut(w, r)
				s.SpotifyAuthorizeRedirect(w, r)
				return
			}

			s.Error(w, err, http.StatusInternalServerError, "failed to find user")
			return
		}

		go func() {
			sessionToken := &model.SessionToken{Base: model.Base{ID: token.(string)}}
			if err := s.DB.Touch(context.Background(), sessionToken); err != nil {
				s.Logger.Err(err, "failed to touch session_token", logger.Field("token", token))
			}
		}()

		next.ServeHTTP(w, r.WithContext(user.Context(r.Context())))
	})
}
