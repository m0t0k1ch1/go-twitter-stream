package stream

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/mrjones/oauth"
)

const (
	UserStreamEndpoint = "https://userstream.twitter.com/1.1/user.json"
)

type Client struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    *oauth.AccessToken
}

type Stream struct {
	Response *http.Response
	Scanner  *bufio.Scanner
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

func (s *Stream) Scan() ([]byte, error) {
	var line []byte
	for s.Scanner.Scan() {
		line = s.Scanner.Bytes()
		if len(line) > 0 {
			break
		}
	}
	if err := s.Scanner.Err(); err != nil {
		return nil, err
	}
	return line, nil
}

func (s *Stream) Close() {
	s.Response.Body.Close()
}
