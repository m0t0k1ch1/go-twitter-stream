package stream

import (
	"bufio"
	"net/http"
)

const (
	UserStreamEndpoint = "https://userstream.twitter.com/1.1/user.json"
)

type Stream struct {
	Response *http.Response
	Scanner  *bufio.Scanner
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
