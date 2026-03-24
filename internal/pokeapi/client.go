package pokeapi

import (
	"net/http"
	"time"

	"github.com/nayemmmmmmmmmm/pokedex/internal/pokecache"
)

// Client _
type Client struct {
	cache      pokecache.Cache
	httpClient http.Client
}

// NewClient _
func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		cache: pokecache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}
