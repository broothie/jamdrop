package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/broothie/queuecumber/config"
	"github.com/broothie/queuecumber/db"
	"github.com/broothie/queuecumber/model"
	"github.com/pkg/errors"
)

type Spotify struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	Logger       *log.Logger
	DB           *db.DB
}

func New(cfg *config.Config, db *db.DB, logger *log.Logger) *Spotify {
	return &Spotify{
		ClientID:     cfg.SpotifyClientID,
		ClientSecret: cfg.SpotifyClientSecret,
		BaseURL:      cfg.BaseURL(),
		Logger:       logger,
		DB:           db,
	}
}

func (s *Spotify) GetUserByID(currentUser *model.User, otherUserID string) (*model.User, error) {
	s.Logger.Println("spotify.GetUserByID")

	if err := s.refreshAccessTokenIfExpired(currentUser); err != nil {
		return nil, errors.Wrapf(err, "failed to refresh access token; user_id: %s", currentUser.ID)
	}

	req, err := http.NewRequest(http.MethodGet, apiPath(fmt.Sprintf("/v1/users/%s", otherUserID)), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request for user_id '%s'", otherUserID)
	}

	s.setBearerAuth(req, currentUser.AccessToken)
	otherUser := new(model.User)
	if err := s.requestToJSON(req, otherUser); err != nil {
		return nil, errors.Wrapf(err, "request failed; user_id: %s, access_token", otherUserID, currentUser.AccessToken)
	}

	otherUser.UpdateAccessTokenExpiration()
	return otherUser, nil
}

func (s *Spotify) QueueSongForUser(user *model.User, songIdentifier string) error {
	s.Logger.Println("spotify.QueueSongForUser")

	if err := s.refreshAccessTokenIfExpired(user); err != nil {
		return errors.Wrapf(err, "failed to refresh access token; user_id: %s", user.ID)
	}

	songID, err := IDFromIdentifier(songIdentifier)
	if err != nil {
		return errors.WithStack(err)
	}

	req, err := http.NewRequest(http.MethodPost, apiPath("/v1/me/player/queue"), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request for song queuing")
	}

	req.URL.RawQuery = url.Values{"uri": {SongURI(songID)}}.Encode()
	s.setBearerAuth(req, user.AccessToken)
	if _, _, err := s.request(req); err != nil {
		return errors.Wrapf(err, "failed to make song queue request; access_token: %s, song_identifier: %s", user.AccessToken, songIdentifier)
	}

	return nil
}

func (s *Spotify) setUserData(accessToken string, user *model.User) error {
	s.Logger.Println("spotify.SetUserData")

	req, err := http.NewRequest(http.MethodGet, apiPath("/v1/me"), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create user data request")
	}

	s.setBearerAuth(req, accessToken)
	if err := s.requestToJSON(req, user); err != nil {
		return errors.Wrapf(err, "failed to make request for user with token '%s'", accessToken)
	}

	user.UpdateAccessTokenExpiration()
	return nil
}

func (s *Spotify) request(r *http.Request) (*http.Response, []byte, error) {
	s.Logger.Printf("%s %s", r.Method, r.URL.String())

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to make request")
	}

	var body []byte
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		return nil, nil, errors.Wrap(err, "failed to read request body")
	}

	if res.StatusCode < 200 || 299 < res.StatusCode {
		return nil, nil, fmt.Errorf("bad response; status %d, body: %s", res.StatusCode, body)
	}

	return res, body, nil
}

func (s *Spotify) requestToJSON(r *http.Request, v interface{}) error {
	_, body, err := s.request(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return errors.Wrap(err, "failed to unmarshal request response")
	}

	return nil
}

func (s *Spotify) setBasicAuth(r *http.Request) {
	r.SetBasicAuth(s.ClientID, s.ClientSecret)
}

func (s *Spotify) setBearerAuth(r *http.Request, token string) {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
}
