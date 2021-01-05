# go-twitter-stream

a Twitter streaming API client for Go

## Example

``` go
package main

import (
	"log"

	"github.com/m0t0k1ch1/go-twitter-stream"
)

func main() {
	c := stream.NewClient(
		"your consumer key",
		"your consumer secret",
		"your access token",
		"your access token secret",
	)

	s, err := c.UserStream(map[string]string{"stringify_friend_ids": "true"})
	if err != nil {
		panic(err)
	}
	defer s.Close()

	for {
		line, err := s.Scan()
		if err != nil {
			log.Println(err)
		}
		log.Println(string(line))
	}
}
```
