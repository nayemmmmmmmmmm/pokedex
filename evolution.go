package main

import (
	"math/rand/v2"
	"strings"
	"time"

	"github.com/nayemmmmmmmmmm/pokedex/internal/pokeapi"
)

type Evolution struct {
	fromPokemon string
	toPokemon   string
	minLevel    int
	timeNeeded  time.Duration
	evolveTime  time.Time
}

type EvolutionTracker struct {
	evolutions map[string]*Evolution
}

func NewEvolutionTracker() *EvolutionTracker {
	return &EvolutionTracker{
		evolutions: make(map[string]*Evolution),
	}
}

func (et *EvolutionTracker) StartEvolution(pokemonName string, toPokemon string, minLevel int, timeNeeded time.Duration) {
	lowerName := strings.ToLower(pokemonName)
	et.evolutions[lowerName] = &Evolution{
		fromPokemon: pokemonName,
		toPokemon:   toPokemon,
		minLevel:    minLevel,
		timeNeeded:  timeNeeded,
		evolveTime:  time.Now().Add(timeNeeded),
	}
}

func (et *EvolutionTracker) CheckEvolutions() []string {
	var evolved []string
	now := time.Now()

	for name, evolution := range et.evolutions {
		if now.After(evolution.evolveTime) {
			evolved = append(evolved, evolution.toPokemon)
			delete(et.evolutions, name)
		}
	}

	return evolved
}

func (et *EvolutionTracker) GetEvolutionStatus(pokemonName string) (string, time.Duration, bool) {
	lowerName := strings.ToLower(pokemonName)
	evolution, exists := et.evolutions[lowerName]
	if !exists {
		return "", 0, false
	}

	timeRemaining := time.Until(evolution.evolveTime)
	if timeRemaining < 0 {
		timeRemaining = 0
	}

	return evolution.toPokemon, timeRemaining, true
}

func (et *EvolutionTracker) CancelEvolution(pokemonName string) {
	lowerName := strings.ToLower(pokemonName)
	delete(et.evolutions, lowerName)
}

// Evolution chains for common Pokemon (shortened for testing)
var evolutionChains = map[string]map[string]interface{}{
	"pidgey": {
		"to":         "pidgeotto",
		"minLevel":   5, 
		"timeNeeded": 10 * time.Second,
	},
	"pidgeotto": {
		"to":         "pidgeot",
		"minLevel":   10, 
		"timeNeeded": 15 * time.Second,
	},
	"rattata": {
		"to":         "raticate",
		"minLevel":   5, 
		"timeNeeded": 8 * time.Second,
	},
	"bulbasaur": {
		"to":         "ivysaur",
		"minLevel":   5, 
		"timeNeeded": 12 * time.Second,
	},
	"ivysaur": {
		"to":         "venusaur",
		"minLevel":   10, 
		"timeNeeded": 18 * time.Second,
	},
	"charmander": {
		"to":         "charmeleon",
		"minLevel":   5, 
		"timeNeeded": 11 * time.Second,
	},
	"charmeleon": {
		"to":         "charizard",
		"minLevel":   10, 
		"timeNeeded": 16 * time.Second,
	},
	"squirtle": {
		"to":         "wartortle",
		"minLevel":   5, 
		"timeNeeded": 12 * time.Second,
	},
	"wartortle": {
		"to":         "blastoise",
		"minLevel":   10, 
		"timeNeeded": 17 * time.Second,
	},
	"caterpie": {
		"to":         "metapod",
		"minLevel":   3, 
		"timeNeeded": 5 * time.Second,
	},
	"metapod": {
		"to":         "butterfree",
		"minLevel":   5, 
		"timeNeeded": 7 * time.Second,
	},
	"weedle": {
		"to":         "kakuna",
		"minLevel":   3, 
		"timeNeeded": 5 * time.Second,
	},
	"kakuna": {
		"to":         "beedrill",
		"minLevel":   5, 
		"timeNeeded": 7 * time.Second,
	},
	"pikachu": {
		"to":         "raichu",
		"minLevel":   8, 
		"timeNeeded": 14 * time.Second,
	},
	"eevee": {
		"to":         randomEeveeEvolution(),
		"minLevel":   6, 
		"timeNeeded": 13 * time.Second,
	},
}

func randomEeveeEvolution() string {
	evolutions := []string{"flareon", "jolteon", "vaporeon"}
	return evolutions[rand.IntN(len(evolutions))]
}

func (pp *PartyPokemon) CanEvolve() bool {
	chain, exists := evolutionChains[strings.ToLower(pp.pokemon.Name)]
	if !exists {
		return false
	}

	minLevel := chain["minLevel"].(int)
	return pp.level >= minLevel
}

func (pp *PartyPokemon) GetEvolutionTarget() (string, bool) {
	chain, exists := evolutionChains[strings.ToLower(pp.pokemon.Name)]
	if !exists {
		return "", false
	}

	if !pp.CanEvolve() {
		return "", false
	}

	return chain["to"].(string), true
}

func (pp *PartyPokemon) Evolve(toPokemon string) error {

	// Create evolved Pokemon
	evolvedPokemon := pokeapi.Pokemon{
		Name: strings.ToLower(toPokemon),
		Stats: make([]struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}, len(pp.pokemon.Stats)),
		Types: make([]struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		}, len(pp.pokemon.Types)),
	}

	// Copy and boost stats
	for i, stat := range pp.pokemon.Stats {
		boost := 1.2 // 20% stat boost on evolution
		if stat.Stat.Name == "hp" {
			boost = 1.3 // 30% HP boost
		}
		evolvedPokemon.Stats[i].BaseStat = int(float64(stat.BaseStat) * boost)
		evolvedPokemon.Stats[i].Effort = stat.Effort
		evolvedPokemon.Stats[i].Stat = stat.Stat
	}

	// Copy types
	for i, t := range pp.pokemon.Types {
		evolvedPokemon.Types[i] = t
	}

	// Update the party Pokemon
	pp.pokemon = evolvedPokemon
	pp.maxHp = calculateHP(evolvedPokemon)
	pp.hp = pp.maxHp // Full heal on evolution

	return nil
}
