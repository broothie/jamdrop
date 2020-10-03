package spotify

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"jamdrop/model"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

func (s *Client) UserAuthorizeURL(state string) string {
	u, _ := url.Parse(AccountsBaseURL)
	u.Path = path.Join(u.Path, "/authorize")
	u.RawQuery = url.Values{
		"client_id":     {s.ClientID},
		"response_type": {"code"},
		"redirect_uri":  {s.AuthRedirectURI()},
		"scope":         {"user-modify-playback-state user-read-currently-playing user-read-playback-state"},
		"state":         {state},
	}.Encode()

	return u.String()
}

func (s *Client) AuthRedirectURI() string {
	return fmt.Sprintf("%s/spotify/authorize/callback", s.BaseURL)
}

func (s *Client) UserFromAuthorizationCode(ctx context.Context, code string) (*model.User, error) {
	s.Logger.Println("spotify.UserFromAuthorizationCode", code)

	user := new(model.User)
	if err := s.setUserTokens(code, user); err != nil {
		return nil, err
	}

	if err := s.setUserData(user.AccessToken, user); err != nil {
		return nil, err
	}

	userExists, err := s.DB.Exists(ctx, model.CollectionUsers, user.ID)
	if err != nil {
		return nil, err
	}

	if userExists {
		updates := []firestore.Update{
			{Path: "display_name", Value: user.DisplayName},
			{Path: "images", Value: user.Images},
			{Path: "access_token", Value: user.AccessToken},
			{Path: "refresh_token", Value: user.RefreshToken},
		}

		if err := s.DB.Update(ctx, user, updates...); err != nil {
			return nil, err
		}
	} else {
		if err := s.DB.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *Client) setUserTokens(code string, user *model.User) error {
	s.Logger.Println("spotify.setUserTokens", code, user.ID)

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

func (s *Client) setUserData(accessToken string, user *model.User) error {
	s.Logger.Println("spotify.setUserData", accessToken, user.ID)

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

func (s *Client) refreshAccessTokenIfExpired(user *model.User) error {
	if user.AccessTokenIsFresh() {
		return nil
	}

	return s.refreshAccessToken(user)
}

func (s *Client) refreshAccessToken(user *model.User) error {
	s.Logger.Println("spotify.refreshAccessToken", user.ID)

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
	go func() {
		updates := []firestore.Update{
			{Path: "refresh_token", Value: user.RefreshToken},
			{Path: "access_token", Value: user.AccessToken},
			{Path: "access_token_expires_at", Value: user.AccessTokenExpiresAt},
		}

		if err := s.DB.Update(context.Background(), user, updates...); err != nil {
			s.Logger.Printf(
				"failed to refresh access_token; user_id: %s, access_token: %s, refresh_token: %s\n",
				user.ID,
				user.AccessToken,
				user.RefreshToken,
			)
		}
	}()

	return nil
}

func (s *Client) setBasicAuth(r *http.Request) {
	r.SetBasicAuth(s.ClientID, s.ClientSecret)
}

func (s *Client) setBearerAuth(r *http.Request, token string) {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
}
