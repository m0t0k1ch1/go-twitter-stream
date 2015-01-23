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

func (c *Client) NewStream(endpoint string) (*Stream, error) {
	consumer := oauth.NewConsumer(
		c.ConsumerKey,
		c.ConsumerSecret,
		oauth.ServiceProvider{},
	)

	res, err := consumer.Post(endpoint, nil, c.AccessToken)
	if err != nil {
		return nil, err
	}

	return &Stream{
		Response: res,
		Scanner:  bufio.NewScanner(res.Body),
	}, nil
}

func (c *Client) UserStream() (*Stream, error) {
	return c.NewStream(UserStreamEndpoint)
}

func (s *Stream) Listen() <-chan []byte {
	ch := make(chan []byte)

	go func() {
		for s.Scanner.Scan() {
			line := s.Scanner.Bytes()
			if len(line) > 0 {
				ch <- line
			}
		}
	}()

	return ch
}

func (s *Stream) Close() {
	s.Response.Body.Close()
}
