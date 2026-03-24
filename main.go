package main

import (
	"time"

	"github.com/nayemmmmmmmmmm/pokedex/internal/pokeapi"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := &config{
		caughtPokemon: map[string]pokeapi.Pokemon{},
		pokeapiClient: pokeClient,
		history:       []string{},
		historyIndex:  -1,
	}

	startRepl(cfg)
}
