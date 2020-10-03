package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"jamdrop/model"

	"github.com/pkg/errors"
)

type SongData struct {
	Name string `json:"name"`
}

func (s *Client) GetUserByID(currentUser *model.User, otherUserID string) (*model.User, error) {
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
		return nil, errors.Wrapf(err, "request failed; user_id: %s, access_token: %s", otherUserID, currentUser.AccessToken)
	}

	otherUser.UpdateAccessTokenExpiration()
	return otherUser, nil
}

func (s *Client) GetSongData(user *model.User, songIdentifier string) (SongData, error) {
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

func (s *Client) QueueSong(user *model.User, songIdentifier string) error {
	s.Logger.Println("spotify.QueueSong", user.ID, songIdentifier)

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

	songURI := SongURI(songID)
	req.URL.RawQuery = url.Values{"uri": {SongURI(songID)}}.Encode()
	s.setBearerAuth(req, user.AccessToken)
	if _, _, err := s.request(req); err != nil {
		if err, isSpotifyError := err.(SpotifyError); isSpotifyError && err.Reason == noActiveDevice {
			return s.setCurrentSong(user, songURI)
		}

		return errors.Wrapf(err, "failed to make song queue request; access_token: %s, song_identifier: %s", user.AccessToken, songIdentifier)
	}

	return nil
}

func (s *Client) setCurrentSong(user *model.User, songURI string) error {
	s.Logger.Println("spotify.setCurrentSong", user.ID, songURI)

	body := fmt.Sprintf(`{"uris":["%s"]}`, songURI)
	req, err := http.NewRequest(http.MethodPut, apiPath("v1/me/player/play"), bytes.NewBufferString(body))
	if err != nil {
		return errors.Wrap(err, "failed to create request for setting song")
	}

	s.setBearerAuth(req, user.AccessToken)
	if _, _, err := s.request(req); err != nil {
		return errors.Wrap(err, "failed to set current song for user")
	}

	return nil
}

func (s *Client) GetCurrentlyPlaying(user *model.User) (bool, error) {
	s.Logger.Println("spotify.GetCurrentlyPlaying", user.ID)

	if err := s.refreshAccessTokenIfExpired(user); err != nil {
		return false, errors.Wrapf(err, "failed to refresh access token; user_id: %s", user.ID)
	}

	req, err := http.NewRequest(http.MethodGet, apiPath("/v1/me/player"), nil)
	if err != nil {
		return false, errors.Wrap(err, "failed to create request for setting song")
	}

	var playerData struct {
		IsPlaying bool `json:"is_playing"`
	}

	s.setBearerAuth(req, user.AccessToken)
	_, body, err := s.request(req)
	if err != nil {
		return false, errors.Wrap(err, "failed to get user player status")
	}

	if len(body) == 0 {
		return false, nil
	}

	if err := json.Unmarshal(body, &playerData); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal request response: %s", body)
	}

	return playerData.IsPlaying, nil
}
