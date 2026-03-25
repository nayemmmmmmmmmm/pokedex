package main

import (
	"testing"
	"time"

	"github.com/nayemmmmmmmmmm/pokedex/internal/pokeapi"
)

func TestNewEvolutionTracker(t *testing.T) {
	tracker := NewEvolutionTracker()
	
	if len(tracker.evolutions) != 0 {
		t.Errorf("Expected empty evolution tracker, got %d evolutions", len(tracker.evolutions))
	}
}

func TestEvolutionTrackerStartEvolution(t *testing.T) {
	tracker := NewEvolutionTracker()
	
	tracker.StartEvolution("pikachu", "raichu", 25, 30*time.Second)
	
	if len(tracker.evolutions) != 1 {
		t.Errorf("Expected 1 evolution, got %d", len(tracker.evolutions))
	}
	
	evolution := tracker.evolutions["pikachu"]
	if evolution.fromPokemon != "pikachu" {
		t.Errorf("Expected fromPokemon 'pikachu', got '%s'", evolution.fromPokemon)
	}
	if evolution.toPokemon != "raichu" {
		t.Errorf("Expected toPokemon 'raichu', got '%s'", evolution.toPokemon)
	}
	if evolution.minLevel != 25 {
		t.Errorf("Expected minLevel 25, got %d", evolution.minLevel)
	}
}

func TestEvolutionTrackerGetEvolutionStatus(t *testing.T) {
	tracker := NewEvolutionTracker()
	
	// Test non-existent evolution
	target, timeRemaining, exists := tracker.GetEvolutionStatus("pikachu")
	if exists {
		t.Error("Expected false for non-existent evolution")
	}
	
	// Test existing evolution
	tracker.StartEvolution("pikachu", "raichu", 25, 30*time.Second)
	
	target, timeRemaining, exists = tracker.GetEvolutionStatus("pikachu")
	if !exists {
		t.Error("Expected true for existing evolution")
	}
	if target != "raichu" {
		t.Errorf("Expected target 'raichu', got '%s'", target)
	}
	if timeRemaining <= 0 || timeRemaining > 30*time.Second {
		t.Errorf("Expected timeRemaining between 0 and 30s, got %v", timeRemaining)
	}
}

func TestEvolutionTrackerCheckEvolutions(t *testing.T) {
	tracker := NewEvolutionTracker()
	
	// Test no evolutions ready
	tracker.StartEvolution("pikachu", "raichu", 25, 30*time.Second)
	evolved := tracker.CheckEvolutions()
	if len(evolved) != 0 {
		t.Errorf("Expected 0 evolved pokemon, got %d", len(evolved))
	}
	
	// Test evolution ready (using very short time)
	tracker.StartEvolution("charmander", "charmeleon", 16, 1*time.Millisecond)
	time.Sleep(10 * time.Millisecond) // Wait for evolution to complete
	evolved = tracker.CheckEvolutions()
	if len(evolved) != 1 {
		t.Errorf("Expected 1 evolved pokemon, got %d", len(evolved))
	}
	if evolved[0] != "charmeleon" {
		t.Errorf("Expected 'charmeleon', got '%s'", evolved[0])
	}
	
	// Check that evolution was removed
	_, _, exists := tracker.GetEvolutionStatus("charmander")
	if exists {
		t.Error("Expected evolution to be removed after completion")
	}
}

func TestEvolutionTrackerCancelEvolution(t *testing.T) {
	tracker := NewEvolutionTracker()
	
	tracker.StartEvolution("pikachu", "raichu", 25, 30*time.Second)
	
	// Verify evolution exists
	_, _, exists := tracker.GetEvolutionStatus("pikachu")
	if !exists {
		t.Error("Expected evolution to exist")
	}
	
	// Cancel evolution
	tracker.CancelEvolution("pikachu")
	
	// Verify evolution is gone
	_, _, exists = tracker.GetEvolutionStatus("pikachu")
	if exists {
		t.Error("Expected evolution to be cancelled")
	}
}

func TestPartyPokemonCanEvolve(t *testing.T) {
	pokemon := pokeapi.Pokemon{
		Name: "pikachu",
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}{
			{BaseStat: 35, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "hp"}},
		},
	}
	
	// Test too low level
	pp := &PartyPokemon{
		pokemon: pokemon,
		level:   10,
		exp:     0,
		hp:      35,
		maxHp:   35,
	}
	
	if pp.CanEvolve() {
		t.Error("Pokemon should not be able to evolve at level 10")
	}
	
	// Test sufficient level
	pp.level = 25
	if !pp.CanEvolve() {
		t.Error("Pokemon should be able to evolve at level 25")
	}
	
	// Test Pokemon with no evolution chain
	noEvoPokemon := pokeapi.Pokemon{Name: "mew"}
	pp.pokemon = noEvoPokemon
	if pp.CanEvolve() {
		t.Error("Pokemon with no evolution chain should not be able to evolve")
	}
}

func TestPartyPokemonGetEvolutionTarget(t *testing.T) {
	pokemon := pokeapi.Pokemon{Name: "pikachu"}
	pp := &PartyPokemon{
		pokemon: pokemon,
		level:   25,
		exp:     0,
		hp:      35,
		maxHp:   35,
	}
	
	target, canEvolve := pp.GetEvolutionTarget()
	if !canEvolve {
		t.Error("Pikachu should be able to evolve at level 25")
	}
	if target != "raichu" {
		t.Errorf("Expected evolution target 'raichu', got '%s'", target)
	}
	
	// Test too low level
	pp.level = 10
	target, canEvolve = pp.GetEvolutionTarget()
	if canEvolve {
		t.Error("Pikachu should not be able to evolve at level 10")
	}
	if target != "" {
		t.Errorf("Expected empty target when can't evolve, got '%s'", target)
	}
}

func TestPartyPokemonEvolve(t *testing.T) {
	pokemon := pokeapi.Pokemon{
		Name: "pikachu",
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}{
			{BaseStat: 35, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "hp"}},
			{BaseStat: 55, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "attack"}},
		},
		Types: []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		}{
			{Type: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "electric"}},
		},
	}
	
	pp := &PartyPokemon{
		pokemon: pokemon,
		level:   25,
		exp:     0,
		hp:      20, // Not full health
		maxHp:   35,
	}
	
	initialHP := pp.maxHp
	err := pp.Evolve("raichu")
	if err != nil {
		t.Errorf("Unexpected error evolving: %v", err)
	}
	
	// Check that Pokemon was updated
	if pp.pokemon.Name != "raichu" {
		t.Errorf("Expected pokemon name 'raichu', got '%s'", pp.pokemon.Name)
	}
	
	// Check that stats increased
	if pp.maxHp <= initialHP {
		t.Errorf("Expected max HP to increase on evolution, got %d (was %d)", pp.maxHp, initialHP)
	}
	
	// Check that HP was fully restored
	if pp.hp != pp.maxHp {
		t.Errorf("Expected full heal on evolution, got %d/%d", pp.hp, pp.maxHp)
	}
}

func TestEvolutionChains(t *testing.T) {
	// Test that evolution chains exist for common Pokemon
	testPokemon := []string{"pidgey", "charmander", "bulbasaur", "squirtle", "pikachu"}
	
	for _, name := range testPokemon {
		chain, exists := evolutionChains[name]
		if !exists {
			t.Errorf("Expected evolution chain for %s", name)
			continue
		}
		
		if _, hasTo := chain["to"]; !hasTo {
			t.Errorf("Expected 'to' field in evolution chain for %s", name)
		}
		
		if _, hasMinLevel := chain["minLevel"]; !hasMinLevel {
			t.Errorf("Expected 'minLevel' field in evolution chain for %s", name)
		}
		
		if _, hasTimeNeeded := chain["timeNeeded"]; !hasTimeNeeded {
			t.Errorf("Expected 'timeNeeded' field in evolution chain for %s", name)
		}
	}
}
