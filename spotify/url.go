package spotify

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	APIBaseURL      = "https://api.spotify.com"
	AccountsBaseURL = "https://accounts.spotify.com"
)

func APIPath(path string) string {
	return fmt.Sprintf("%s/%s", APIBaseURL, strings.TrimPrefix(path, "/"))
}

func AccountsPath(path string) string {
	return fmt.Sprintf("%s/%s", AccountsBaseURL, strings.TrimPrefix(path, "/"))
}

func (c *Client) UserAuthorizeURL() string {
	u, _ := url.Parse(AccountsBaseURL)
	u.RawQuery = url.Values{
		"client_id":     {c.ClientID},
		"response_type": {"code"},
		"redirect_uri":  {c.AuthRedirectURI()},
		"scope":         {"user-modify-playback-state user-read-currently-playing"},
	}.Encode()

	return u.String()
}

func (c *Client) AuthRedirectURI() string {
	return fmt.Sprintf("%s/spotify/authorize/callback", c.BaseURL)
}
