package stream

import (
	"bufio"
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

func (c *Client) NewStream(endpoint string, param map[string]string) (*Stream, error) {
	consumer := oauth.NewConsumer(
		c.ConsumerKey,
		c.ConsumerSecret,
		oauth.ServiceProvider{},
	)

	res, err := consumer.Post(endpoint, param, c.AccessToken)
	if err != nil {
		return nil, err
	}

	return &Stream{
		Response: res,
		Scanner:  bufio.NewScanner(res.Body),
	}, nil
}

func (c *Client) UserStream(param map[string]string) (*Stream, error) {
	return c.NewStream(UserStreamEndpoint, param)
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
