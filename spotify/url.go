package spotify

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

const (
	APIBaseURL      = "https://api.spotify.com"
	AccountsBaseURL = "https://accounts.spotify.com"
)

func apiPath(path string) string {
	return fmt.Sprintf("%s/%s", APIBaseURL, strings.TrimPrefix(path, "/"))
}

func accountsPath(path string) string {
	return fmt.Sprintf("%s/%s", AccountsBaseURL, strings.TrimPrefix(path, "/"))
}

func (s *Spotify) UserAuthorizeURL() string {
	u, _ := url.Parse(AccountsBaseURL)
	u.Path = path.Join(u.Path, "/authorize")
	u.RawQuery = url.Values{
		"client_id":     {s.ClientID},
		"response_type": {"code"},
		"redirect_uri":  {s.AuthRedirectURI()},
		"scope":         {"user-modify-playback-state user-read-currently-playing"},
	}.Encode()

	return u.String()
}

func (s *Spotify) AuthRedirectURI() string {
	return fmt.Sprintf("%s/spotify/authorize/callback", s.BaseURL)
}
