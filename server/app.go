package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/broothie/queuecumber/spotify"

	"github.com/broothie/queuecumber/db"

	"github.com/broothie/queuecumber/model"
)

var (
	index = template.Must(template.ParseFiles("views/index.html"))
)

func (s *Server) Index() http.HandlerFunc {
	type Data struct {
		Flashes []interface{}
		User    *model.User
	}

	return s.RequireLoggedIn(func(w http.ResponseWriter, r *http.Request) {
		user, _ := model.UserFromContext(r.Context())
		s.Logger.Println("server.Index", "user_id", user.ID)

		data := Data{Flashes: s.GetFlashes(w, r), User: user}
		if err := index.Execute(w, data); err != nil {
			s.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (s *Server) AddFriend() http.HandlerFunc {
	return s.RequireLoggedIn(func(w http.ResponseWriter, r *http.Request) {
		defer http.Redirect(w, r, "/app", http.StatusPermanentRedirect)

		user, _ := model.UserFromContext(r.Context())
		friendIdentifier := r.FormValue("user_identifier")
		s.Logger.Println("server.AddFriend", "user_id", user.ID, "friend_identifier", friendIdentifier)

		friendID, err := spotify.IDFromIdentifier(friendIdentifier)
		if err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, fmt.Sprintf("invalid user identifier: %s", friendIdentifier))
			return
		}

		friend, err := s.DB.GetUserByID(r.Context(), friendID)
		if err != nil {
			s.Logger.Println(err)

			flash := "An error occurred. Please try again."
			if db.IsNotFound(err) {
				friend, err := s.Spotify.GetUserByID(user.AccessToken, friendID)
				if err == nil {
					flash = fmt.Sprintf("%s has not connected their Spotify account to queuecumber yet.", friend.DisplayName)
				}
			}

			s.Flash(w, r, flash)
			return
		}

		user.AddFriend(friend)
		if err := s.DB.UpsertUser(r.Context(), user); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		s.Flash(w, r, "Friend added successfully!")
	})
}

func (s *Server) QueueSong() http.HandlerFunc {
	return s.RequireLoggedIn(func(w http.ResponseWriter, r *http.Request) {
		defer http.Redirect(w, r, "/app", http.StatusPermanentRedirect)

		user, _ := model.UserFromContext(r.Context())
		friendID := r.FormValue("user_id")
		songIdentifier := r.FormValue("song_identifier")
		s.Logger.Println("server.QueueSong", "user_id", user.ID, "song_identifier", songIdentifier)

		songID, err := spotify.IDFromIdentifier(songIdentifier)
		if err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, fmt.Sprintf("invalid song identifier: %s", songIdentifier))
			return
		}

		songURI := spotify.SongURI(songID)
		friend, err := s.DB.GetUserByID(r.Context(), friendID)
		if err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		if err := s.Spotify.QueueSongForUser(friend.AccessToken, songURI); err != nil {
			s.Logger.Println(err)
			s.Flash(w, r, "An error occurred. Please try again.")
			return
		}

		s.Flash(w, r, "Song queued successfully!")
	})
}
