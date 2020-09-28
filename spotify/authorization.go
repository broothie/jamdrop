package spotify

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/broothie/queuecumber/model"
	"github.com/pkg/errors"
)

func (s *Spotify) SetUserTokens(code string, user *model.User) error {
	s.Logger.Println("spotify.SetUserTokens")

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
	if err := s.requestToJSON(req, user); err != nil {
		return errors.Wrapf(err, "failed to make request for token with code '%s'", code)
	}

	user.UpdateAccessTokenExpiration()
	return nil
}

func (s *Spotify) RefreshAccessTokenIfExpired(user *model.User) error {
	if user.AccessTokenIsFresh() {
		return nil
	}

	return s.RefreshAccessToken(user)
}

func (s *Spotify) RefreshAccessToken(user *model.User) error {
	s.Logger.Println("spotify.RefreshAccessToken")

	body := url.Values{"grant_type": {"refresh_token"}, "refresh_token": {user.RefreshToken}}
	req, err := http.NewRequest(http.MethodPost, accountsPath("/api/token"), bytes.NewBufferString(body.Encode()))
	if err != nil {
		return errors.Wrap(err, "failed to create refresh token request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	s.setBasicAuth(req)
	if err := s.requestToJSON(req, user); err != nil {
		return errors.Wrapf(err, "failed to make request for refresh token with refresh_token '%s'", user.RefreshToken)
	}

	user.UpdateAccessTokenExpiration()
	return nil
}
