package stream

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/mrjones/oauth"
)

type Client struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    *oauth.AccessToken
}

func NewClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *Client {
	return &Client{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		AccessToken: &oauth.AccessToken{
			Token:  accessToken,
			Secret: accessTokenSecret,
		},
	}
}

func (c *Client) NewStream(method, endpoint string, param map[string]string) (*Stream, error) {
	consumer := oauth.NewConsumer(
		c.ConsumerKey,
		c.ConsumerSecret,
		oauth.ServiceProvider{},
	)

	var resp *http.Response
	var err error

	switch method {
	case "GET":
		resp, err = consumer.Get(endpoint, param, c.AccessToken)
	case "POST":
		resp, err = consumer.Post(endpoint, param, c.AccessToken)
	default:
		return nil, fmt.Errorf("\"%s\" method is not supported", method)
	}
	if err != nil {
		return nil, err
	}

	return &Stream{
		Response: resp,
		Scanner:  bufio.NewScanner(resp.Body),
	}, nil
}

func (c *Client) UserStream(param map[string]string) (*Stream, error) {
	return c.NewStream("GET", UserStreamEndpoint, param)
}
