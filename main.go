package main

import (
	"time"

	"github.com/nayemmmmmmmmmm/pokedex/internal/pokeapi"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := &config{
		caughtPokemon:    map[string]pokeapi.Pokemon{},
		pokeapiClient:    pokeClient,
		history:          []string{},
		historyIndex:     -1,
		party:            NewParty(6), // Max 6 Pokemon in party
		evolutionTracker: NewEvolutionTracker(),
		explorationState: NewExplorationState(),
	}

	startRepl(cfg)
}
