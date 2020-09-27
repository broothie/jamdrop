package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/broothie/queuecumber/config"
	"github.com/broothie/queuecumber/model"
	"github.com/pkg/errors"
)

type Spotify struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	Logger       *log.Logger
}

func New(cfg *config.Config, logger *log.Logger) *Spotify {
	return &Spotify{
		ClientID:     cfg.SpotifyClientID,
		ClientSecret: cfg.SpotifyClientSecret,
		BaseURL:      cfg.BaseURL(),
		Logger:       logger,
	}
}

func (s *Spotify) GetUserByID(accessToken, userID string) (*model.User, error) {
	s.Logger.Println("spotify.GetUserByID", userID)

	req, err := http.NewRequest(http.MethodGet, apiPath(fmt.Sprintf("/v1/users/%s", userID)), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request for user_id '%s'", userID)
	}

	s.setBearerAuth(req, accessToken)
	user := new(model.User)
	if err := s.requestToJSON(req, user); err != nil {
		return nil, errors.Wrapf(err, "failed to read json response for user data with user_id '%s'", userID)
	}

	return user, nil
}

func (s *Spotify) QueueSongForUser(accessToken, songURI string) error {
	s.Logger.Println("spotify.QueueSongForUser", "access_token", accessToken, "song_uri", songURI)

	req, err := http.NewRequest(http.MethodPost, apiPath("/v1/me/player/queue"), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request for song queuing")
	}

	req.URL.RawQuery = url.Values{"uri": {songURI}}.Encode()
	s.setBearerAuth(req, accessToken)
	_, err = s.request(req)
	return errors.Wrapf(err, "failed to make request for song queueing with access_token '%s', song_uri '%s'", accessToken, songURI)
}

func (s *Spotify) SetUserTokens(code string, user *model.User) error {
	s.Logger.Println("spotify.SetUserTokens", code)

	body := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {s.AuthRedirectURI()},
	}

	req, err := http.NewRequest(http.MethodPost, accountsPath("/api/token"), bytes.NewBufferString(body.Encode()))
	if err != nil {
		return errors.Wrap(err, "failed to create token request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	s.setBasicAuth(req)
	return errors.Wrapf(s.requestToJSON(req, user), "failed to make request for token with code '%s'", code)
}

func (s *Spotify) RefreshAccessToken(user *model.User) error {
	s.Logger.Println("spotify.RefreshAccessToken", user.ID)

	body := url.Values{"grant_type": {"refresh_token"}, "refresh_token": {user.RefreshToken}}
	req, err := http.NewRequest(http.MethodPost, accountsPath("/api/token"), bytes.NewBufferString(body.Encode()))
	if err != nil {
		return errors.Wrap(err, "failed to create refresh token request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	s.setBasicAuth(req)
	return errors.Wrapf(s.requestToJSON(req, user), "failed to make request for refresh token with refresh_token '%s'", user.RefreshToken)
}

func (s *Spotify) SetUserData(accessToken string, user *model.User) error {
	s.Logger.Println("spotify.SetUserData", accessToken)

	req, err := http.NewRequest(http.MethodGet, apiPath("/v1/me"), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create user data request")
	}

	s.setBearerAuth(req, accessToken)
	return errors.Wrapf(s.requestToJSON(req, user), "failed to make request for user with token '%s'", accessToken)
}

func (s *Spotify) request(r *http.Request) (*http.Response, error) {
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}

	if res.StatusCode < 200 || 299 < res.StatusCode {
		return nil, fmt.Errorf("bad response: %d", res.StatusCode)
	}

	return res, nil
}

func (s *Spotify) requestToJSON(r *http.Request, v interface{}) error {
	res, err := s.request(r)
	if err != nil {
		return errors.WithStack(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read request body")
	}

	return errors.Wrap(json.Unmarshal(body, v), "failed to unmarshal request response")
}

func (s *Spotify) setBasicAuth(r *http.Request) {
	r.SetBasicAuth(s.ClientID, s.ClientSecret)
}

func (s *Spotify) setBearerAuth(r *http.Request, token string) {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
}
