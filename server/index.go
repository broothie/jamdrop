package server

import (
	"html/template"
	"net/http"

	"github.com/broothie/queuecumber/model"
)

func (s *Server) AppRedirect(w http.ResponseWriter, r *http.Request) {
	s.Logger.Println("server.AppRedirect")
	http.Redirect(w, r, "/app", http.StatusPermanentRedirect)
}

func (s *Server) Index() http.HandlerFunc {
	type Data struct {
		Flashes []interface{}
		User    *model.User
		Shares  []*model.User
		Sharers []*model.User
	}

	index := template.Must(template.ParseFiles("views/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println("server.Index")

		user, _ := model.UserFromContext(r.Context())
		var sharers []*model.User
		sharersDone := make(chan struct{})
		go func() {
			defer close(sharersDone)

			users, err := s.DB.GetUserSharers(r.Context(), user)
			if err != nil {
				s.Logger.Println(err)
				return
			}

			sharers = users
		}()

		shares, err := s.DB.GetUserShares(r.Context(), user)
		if err != nil {
			s.Flash(w, r, "Failed to retrieve followers.")
		}

		<-sharersDone
		data := Data{
			Flashes: s.GetFlashes(w, r),
			User:    user,
			Shares:  shares,
			Sharers: sharers,
		}

		// Allows "hot" page reloading
		if s.App.Config.IsDevelopment() {
			index = template.Must(template.ParseFiles("views/index.html"))
		}

		if err := index.Execute(w, data); err != nil {
			s.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
