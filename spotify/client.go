package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/broothie/queuecumber/config"
	"github.com/pkg/errors"
)

type Client struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
}

func New(cfg *config.Config) *Client {
	return &Client{
		ClientID:     cfg.SpotifyClientID,
		ClientSecret: cfg.SpotifyClientSecret,
		BaseURL:      cfg.BaseURL(),
	}
}

func (c *Client) Token(code string) (map[string]string, error) {
	body := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {c.AuthRedirectURI()},
	}.Encode()

	req, err := http.NewRequest(http.MethodPost, APIPath("/api/token"), bytes.NewBufferString(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create token request")
	}

	c.setBasicAuth(req)
	var data map[string]string
	if err := c.RequestToJSON(req, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) Request(req *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}

	body, err := ioutil.ReadAll(res.Body)
	return body, errors.Wrap(err, "failed to read request response")
}

func (c *Client) RequestToJSON(req *http.Request, v interface{}) error {
	data, err := c.Request(req)
	if err != nil {
		return err
	}

	return errors.Wrap(json.Unmarshal(data, v), "failed to unmarshal request response")
}

func (c *Client) setBasicAuth(r *http.Request) {
	r.SetBasicAuth(c.ClientID, c.ClientSecret)
}

func (c *Client) setBearerAuth(r *http.Request, token string) {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
}
