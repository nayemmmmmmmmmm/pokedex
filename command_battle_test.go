package main

import (
	"testing"

	"github.com/nayemmmmmmmmmm/pokedex/internal/pokeapi"
)

func TestCalculateHP(t *testing.T) {
	pokemon := pokeapi.Pokemon{
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}{
			{BaseStat: 45, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "hp"}},
		},
	}

	hp := calculateHP(pokemon)
	if hp != 45 {
		t.Errorf("Expected HP 45, got %d", hp)
	}
}

func TestCalculateAttack(t *testing.T) {
	pokemon := pokeapi.Pokemon{
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}{
			{BaseStat: 49, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "attack"}},
		},
	}

	attack := calculateAttack(pokemon)
	if attack != 49 {
		t.Errorf("Expected Attack 49, got %d", attack)
	}
}

func TestBattlePokemonIsAlive(t *testing.T) {
	bp := &BattlePokemon{
		pokemon: pokeapi.Pokemon{Name: "test"},
		hp:      50,
		maxHp:   50,
	}

	if !bp.isAlive() {
		t.Error("Pokemon should be alive with HP > 0")
	}

	bp.hp = 0
	if bp.isAlive() {
		t.Error("Pokemon should not be alive with HP = 0")
	}

	bp.hp = -10
	if bp.isAlive() {
		t.Error("Pokemon should not be alive with HP < 0")
	}
}

func TestGetTypeEffectiveness(t *testing.T) {
	// Test super effective
	fire := pokeapi.Pokemon{
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
			}{Name: "fire"}},
		},
	}

	grass := pokeapi.Pokemon{
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
			}{Name: "grass"}},
		},
	}

	effectiveness := getTypeEffectiveness(fire, grass)
	if effectiveness != 2.0 {
		t.Errorf("Expected fire vs grass to be 2.0, got %f", effectiveness)
	}

	// Test not very effective
	water := pokeapi.Pokemon{
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
			}{Name: "water"}},
		},
	}

	effectiveness = getTypeEffectiveness(fire, water)
	if effectiveness != 0.5 {
		t.Errorf("Expected fire vs water to be 0.5, got %f", effectiveness)
	}
}

func TestDetermineTurnOrder(t *testing.T) {
	fastPokemon := pokeapi.Pokemon{
		Name: "fast",
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}{
			{BaseStat: 100, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "speed"}},
		},
	}

	slowPokemon := pokeapi.Pokemon{
		Name: "slow",
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}{
			{BaseStat: 50, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "speed"}},
		},
	}

	fast := &BattlePokemon{pokemon: fastPokemon, hp: 100, maxHp: 100, isPlayer: true}
	slow := &BattlePokemon{pokemon: slowPokemon, hp: 100, maxHp: 100, isPlayer: false}

	first, second := determineTurnOrder(fast, slow)
	if first != fast {
		t.Error("Fast pokemon should go first")
	}
	if second != slow {
		t.Error("Slow pokemon should go second")
	}
}
