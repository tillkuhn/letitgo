package httpserver

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"time"
)

// ExpensiveHTTPClient Rate Limited HTTP Client
// thanks https://medium.com/mflow/rate-limiting-in-golang-http-client-a22fba15861a
// from GIST https://gist.githubusercontent.com/MelchiSalins/27c11566184116ec1629a0726e0f9af5/raw/71c0b5a9451548013f111c428fe9291e5596b99d/http-rate-limit.go

type ExpensiveHTTPClient struct {
	client      *http.Client
	Ratelimiter *rate.Limiter
}

// Do dispatches the HTTP request to the network
func (c *ExpensiveHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Comment out the below 5 lines to turn off ratelimiting
	ctx := context.Background()
	err := c.Ratelimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// NewClient return http client with a ratelimiter
func NewClient(rl *rate.Limiter) *ExpensiveHTTPClient {
	c := &ExpensiveHTTPClient{
		client:      http.DefaultClient,
		Ratelimiter: rl,
	}
	return c
}

func Run() {
	rl := rate.NewLimiter(rate.Every(2*time.Second), 6) // 4 request every 2 seconds
	c := NewClient(rl)
	reqURL := "https://httpbin.org/json"
	req, _ := http.NewRequest("GET", reqURL, nil)
	for i := 0; i < 100; i++ {
		resp, err := c.Do(req)
		if err != nil {
			fmt.Printf("err %s", err.Error())
			fmt.Println(resp.StatusCode)
			return
		}
		if resp.StatusCode == 429 {
			fmt.Printf("Rate limit reached after %d requests\n", i)
			return
		}
		b, _ := io.ReadAll(resp.Body)
		s := string(b)
		fmt.Printf("#%d %v %10v\n", i, time.Now(), s[0:20])
	}
}
