package spotify

import (
	"fmt"
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
	s.Logger.Println("spotify.GetUserByID", currentUser.ID, otherUserID)

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

type SongData struct {
	Name string `json:"name"`
}

func (s *Spotify) GetSongData(user *model.User, songIdentifier string) (SongData, error) {
	s.Logger.Println("spotify.GetSongData", user.ID, songIdentifier)

	if err := s.refreshAccessTokenIfExpired(user); err != nil {
		return SongData{}, err
	}

	songID, err := IDFromIdentifier(songIdentifier)
	if err != nil {
		return SongData{}, err
	}

	req, err := http.NewRequest(http.MethodGet, apiPath("/v1/tracks/%s", songID), nil)
	if err != nil {
		return SongData{}, errors.Wrapf(err, "")
	}

	s.setBearerAuth(req, user.AccessToken)
	var songData SongData
	if err := s.requestToJSON(req, &songData); err != nil {
		return SongData{}, err
	}

	return songData, nil
}

func (s *Spotify) QueueSongForUser(user *model.User, songIdentifier string) error {
	s.Logger.Println("spotify.QueueSongForUser", user.ID, songIdentifier)

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
