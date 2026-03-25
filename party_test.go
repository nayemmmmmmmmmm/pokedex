package main

import (
	"testing"

	"github.com/nayemmmmmmmmmm/pokedex/internal/pokeapi"
)

func TestNewParty(t *testing.T) {
	party := NewParty(6)

	if party.maxSize != 6 {
		t.Errorf("Expected max size 6, got %d", party.maxSize)
	}

	if len(party.pokemon) != 0 {
		t.Errorf("Expected empty party, got %d pokemon", len(party.pokemon))
	}

	if party.IsFull() {
		t.Error("New party should not be full")
	}
}

func TestPartyAddPokemon(t *testing.T) {
	party := NewParty(2)
	pokemon := pokeapi.Pokemon{Name: "pikachu"}

	err := party.AddPokemon(pokemon)
	if err != nil {
		t.Errorf("Unexpected error adding pokemon: %v", err)
	}

	if len(party.pokemon) != 1 {
		t.Errorf("Expected 1 pokemon, got %d", len(party.pokemon))
	}

	// Test adding duplicate
	err = party.AddPokemon(pokemon)
	if err == nil {
		t.Error("Expected error when adding duplicate pokemon")
	}

	// Test adding to full party
	pokemon2 := pokeapi.Pokemon{Name: "charizard"}
	party.AddPokemon(pokemon2)

	pokemon3 := pokeapi.Pokemon{Name: "bulbasaur"}
	err = party.AddPokemon(pokemon3)
	if err == nil {
		t.Error("Expected error when adding to full party")
	}
}

func TestPartyRemovePokemon(t *testing.T) {
	party := NewParty(6)
	pokemon := pokeapi.Pokemon{Name: "pikachu"}

	party.AddPokemon(pokemon)

	err := party.RemovePokemon("pikachu")
	if err != nil {
		t.Errorf("Unexpected error removing pokemon: %v", err)
	}

	if len(party.pokemon) != 0 {
		t.Errorf("Expected empty party, got %d pokemon", len(party.pokemon))
	}

	// Test removing non-existent pokemon
	err = party.RemovePokemon("charizard")
	if err == nil {
		t.Error("Expected error when removing non-existent pokemon")
	}
}

func TestPartyGetPokemon(t *testing.T) {
	party := NewParty(6)
	pokemon := pokeapi.Pokemon{Name: "pikachu"}

	party.AddPokemon(pokemon)

	retrieved, err := party.GetPokemon("pikachu")
	if err != nil {
		t.Errorf("Unexpected error getting pokemon: %v", err)
	}

	if retrieved.pokemon.Name != "pikachu" {
		t.Errorf("Expected pikachu, got %s", retrieved.pokemon.Name)
	}

	// Test getting non-existent pokemon
	_, err = party.GetPokemon("charizard")
	if err == nil {
		t.Error("Expected error when getting non-existent pokemon")
	}
}

func TestPartyListPokemon(t *testing.T) {
	party := NewParty(6)
	pokemon1 := pokeapi.Pokemon{Name: "pikachu"}
	pokemon2 := pokeapi.Pokemon{Name: "charizard"}

	party.AddPokemon(pokemon1)
	party.AddPokemon(pokemon2)

	names := party.ListPokemon()
	if len(names) != 2 {
		t.Errorf("Expected 2 names, got %d", len(names))
	}
}

func TestPartyPokemonGainExp(t *testing.T) {
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
	}

	pp := &PartyPokemon{
		pokemon: pokemon,
		level:   5,
		exp:     450, // Close to leveling up
		hp:      35,
		maxHp:   35,
	}

	// Test gaining exp without leveling up
	pp.GainExp(30)
	if pp.exp != 480 {
		t.Errorf("Expected 480 exp, got %d", pp.exp)
	}
	if pp.level != 5 {
		t.Errorf("Expected level 5, got %d", pp.level)
	}

	// Test gaining enough exp to level up
	initialHP := pp.maxHp
	pp.GainExp(120) // 480 + 120 = 600, need 500 for level 5, so should level up to 6
	if pp.level != 6 {
		t.Errorf("Expected level 6, got %d", pp.level)
	}
	if pp.exp < 0 {
		t.Errorf("Expected non-negative exp, got %d", pp.exp)
	}
	if pp.maxHp <= initialHP {
		t.Errorf("Expected max HP to increase on level up, got %d (was %d)", pp.maxHp, initialHP)
	}
	if pp.hp != pp.maxHp {
		t.Errorf("Expected full heal on level up, got %d/%d", pp.hp, pp.maxHp)
	}
}

func TestPartyPokemonHP(t *testing.T) {
	pp := &PartyPokemon{
		pokemon: pokeapi.Pokemon{Name: "pikachu"},
		level:   5,
		exp:     0,
		hp:      35,
		maxHp:   35,
	}

	// Test is alive
	if !pp.IsAlive() {
		t.Error("Pokemon should be alive with HP > 0")
	}

	// Test take damage
	pp.TakeDamage(10)
	if pp.hp != 25 {
		t.Errorf("Expected 25 HP after taking 10 damage, got %d", pp.hp)
	}

	// Test take damage below 0
	pp.TakeDamage(30)
	if pp.hp != 0 {
		t.Errorf("Expected 0 HP after taking 30 damage from 25 HP, got %d", pp.hp)
	}

	if pp.IsAlive() {
		t.Error("Pokemon should not be alive with 0 HP")
	}

	// Test heal
	pp.hp = 10
	pp.Heal(15)
	if pp.hp != 25 {
		t.Errorf("Expected 25 HP after healing 15 from 10, got %d", pp.hp)
	}

	// Test heal above max
	pp.Heal(20)
	if pp.hp != pp.maxHp {
		t.Errorf("Expected max HP %d after overhealing, got %d", pp.maxHp, pp.hp)
	}

	// Test full heal
	pp.hp = 5
	pp.FullHeal()
	if pp.hp != pp.maxHp {
		t.Errorf("Expected max HP %d after full heal, got %d", pp.maxHp, pp.hp)
	}
}

func TestCalculateExpGain(t *testing.T) {
	defeatedPokemon := pokeapi.Pokemon{
		Name:           "pikachu",
		BaseExperience: 64,
	}

	// Test at low level
	exp := calculateExpGain(defeatedPokemon, 5)
	expected := int(float64(64) * (1.0 + 0.5)) // 1.0 + (5 * 0.1)
	if exp != expected {
		t.Errorf("Expected %d exp at level 5, got %d", expected, exp)
	}

	// Test at higher level
	exp = calculateExpGain(defeatedPokemon, 10)
	expected = int(float64(64) * (1.0 + 1.0)) // 1.0 + (10 * 0.1)
	if exp != expected {
		t.Errorf("Expected %d exp at level 10, got %d", expected, exp)
	}
}
