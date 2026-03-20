package pokeapi

import (
	"net/http"
	"time"
)

// Client _
type Client struct {
	httpClient http.Client
}

// NewClient _
func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}