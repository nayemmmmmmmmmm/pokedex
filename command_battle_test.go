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

func TestCalculateHP_Default(t *testing.T) {
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

	hp := calculateHP(pokemon)
	if hp != 50 {
		t.Errorf("Expected default HP 50, got %d", hp)
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

func TestCalculateAttack_Default(t *testing.T) {
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

	attack := calculateAttack(pokemon)
	if attack != 50 {
		t.Errorf("Expected default Attack 50, got %d", attack)
	}
}

func TestCalculateDefense(t *testing.T) {
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
			}{Name: "defense"}},
		},
	}

	defense := calculateDefense(pokemon)
	if defense != 49 {
		t.Errorf("Expected Defense 49, got %d", defense)
	}
}

func TestCalculateDefense_Default(t *testing.T) {
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

	defense := calculateDefense(pokemon)
	if defense != 50 {
		t.Errorf("Expected default Defense 50, got %d", defense)
	}
}

func TestCalculateSpeed(t *testing.T) {
	pokemon := pokeapi.Pokemon{
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}{
			{BaseStat: 60, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "speed"}},
		},
	}

	speed := calculateSpeed(pokemon)
	if speed != 60 {
		t.Errorf("Expected Speed 60, got %d", speed)
	}
}

func TestCalculateSpeed_Default(t *testing.T) {
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

	speed := calculateSpeed(pokemon)
	if speed != 50 {
		t.Errorf("Expected default Speed 50, got %d", speed)
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

	// Test normal effectiveness
	normal := pokeapi.Pokemon{
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
			}{Name: "normal"}},
		},
	}

	effectiveness = getTypeEffectiveness(normal, normal)
	if effectiveness != 1.0 {
		t.Errorf("Expected normal vs normal to be 1.0, got %f", effectiveness)
	}

	// Test no effect (ghost vs normal)
	ghost := pokeapi.Pokemon{
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
			}{Name: "ghost"}},
		},
	}

	effectiveness = getTypeEffectiveness(ghost, normal)
	if effectiveness != 0 {
		t.Errorf("Expected ghost vs normal to be 0, got %f", effectiveness)
	}

	// Test dual types
	grassPoison := pokeapi.Pokemon{
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
			{Type: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "poison"}},
		},
	}

	effectiveness = getTypeEffectiveness(fire, grassPoison)
	// fire vs grass = 2.0, fire vs poison = 1.0, so 2.0 * 1.0 = 2.0
	if effectiveness != 2.0 {
		t.Errorf("Expected fire vs grass/poison to be 2.0, got %f", effectiveness)
	}

	// Test 4x effectiveness (fire vs grass/poison where both are weak)
	psychic := pokeapi.Pokemon{
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
			}{Name: "psychic"}},
		},
	}

	fightingPoison := pokeapi.Pokemon{
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
			}{Name: "fighting"}},
			{Type: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "poison"}},
		},
	}

	effectiveness = getTypeEffectiveness(psychic, fightingPoison)
	// psychic vs fighting = 2.0, psychic vs poison = 2.0, so 2.0 * 2.0 = 4.0
	if effectiveness != 4.0 {
		t.Errorf("Expected psychic vs fighting/poison to be 4.0, got %f", effectiveness)
	}
}

func TestGetTypeMultiplier(t *testing.T) {
	// Test existing type matchup
	multiplier := getTypeMultiplier("fire", "grass")
	if multiplier != 2.0 {
		t.Errorf("Expected fire vs grass multiplier to be 2.0, got %f", multiplier)
	}

	// Test non-existent type (should return 1.0)
	multiplier = getTypeMultiplier("unknown", "grass")
	if multiplier != 1.0 {
		t.Errorf("Expected unknown vs grass multiplier to be 1.0, got %f", multiplier)
	}

	// Test existing type but non-existent defender type
	multiplier = getTypeMultiplier("fire", "unknown")
	if multiplier != 1.0 {
		t.Errorf("Expected fire vs unknown multiplier to be 1.0, got %f", multiplier)
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

	// Test equal speed (first parameter should go first)
	equalSpeedPokemon1 := pokeapi.Pokemon{
		Name: "equal1",
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}{
			{BaseStat: 75, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "speed"}},
		},
	}

	equalSpeedPokemon2 := pokeapi.Pokemon{
		Name: "equal2",
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}{
			{BaseStat: 75, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "speed"}},
		},
	}

	equal1 := &BattlePokemon{pokemon: equalSpeedPokemon1, hp: 100, maxHp: 100, isPlayer: true}
	equal2 := &BattlePokemon{pokemon: equalSpeedPokemon2, hp: 100, maxHp: 100, isPlayer: false}

	first, second = determineTurnOrder(equal1, equal2)
	if first != equal1 {
		t.Error("When speeds are equal, first parameter should go first")
	}
	if second != equal2 {
		t.Error("When speeds are equal, second parameter should go second")
	}
}

func TestCalculateDamage(t *testing.T) {
	attacker := &BattlePokemon{
		pokemon: pokeapi.Pokemon{
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
				}{Name: "normal"}},
			},
		},
		hp:       100,
		maxHp:    100,
		isPlayer: true,
	}

	defender := &BattlePokemon{
		pokemon: pokeapi.Pokemon{
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
				}{Name: "defense"}},
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
				}{Name: "normal"}},
			},
		},
		hp:       100,
		maxHp:    100,
		isPlayer: false,
	}

	damage := calculateDamage(attacker, defender)
	if damage <= 0 {
		t.Errorf("Expected positive damage, got %d", damage)
	}

	// Test super effective damage
	attacker.pokemon.Types[0].Type.Name = "fire"
	defender.pokemon.Types[0].Type.Name = "grass"

	superDamage := calculateDamage(attacker, defender)
	if superDamage <= damage {
		t.Errorf("Expected super effective damage (%d) to be greater than normal damage (%d)", superDamage, damage)
	}

	// Test not very effective damage
	defender.pokemon.Types[0].Type.Name = "water"

	weakDamage := calculateDamage(attacker, defender)
	if weakDamage >= damage {
		t.Errorf("Expected not very effective damage (%d) to be less than normal damage (%d)", weakDamage, damage)
	}
}

func TestExecuteTurn(t *testing.T) {
	attacker := &BattlePokemon{
		pokemon: pokeapi.Pokemon{
			Name: "attacker",
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
				}{Name: "attack"}},
				{BaseStat: 50, Stat: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{Name: "defense"}},
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
				}{Name: "normal"}},
			},
		},
		hp:       100,
		maxHp:    100,
		isPlayer: true,
	}

	defender := &BattlePokemon{
		pokemon: pokeapi.Pokemon{
			Name: "defender",
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
				}{Name: "attack"}},
				{BaseStat: 10, Stat: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{Name: "defense"}},
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
				}{Name: "normal"}},
			},
		},
		hp:       30,
		maxHp:    30,
		isPlayer: false,
	}

	battle := &Battle{
		playerPokemon:   attacker,
		opponentPokemon: defender,
		turnCount:       1,
		isActive:        true,
	}

	initialHP := defender.hp
	executeTurn(attacker, defender, battle)

	if defender.hp >= initialHP {
		t.Error("Defender HP should decrease after attack")
	}

	if !battle.isActive && defender.isAlive() {
		t.Error("Battle should remain active if defender is still alive")
	}

	// Test fainting
	defender.hp = 5
	executeTurn(attacker, defender, battle)

	if defender.hp != 0 {
		t.Errorf("Defender HP should be 0, got %d", defender.hp)
	}

	if battle.isActive {
		t.Error("Battle should end when defender faints")
	}
}

func TestExecuteTurn_SkipIfFainted(t *testing.T) {
	attacker := &BattlePokemon{
		pokemon:  pokeapi.Pokemon{Name: "attacker"},
		hp:       0,
		maxHp:    100,
		isPlayer: true,
	}

	defender := &BattlePokemon{
		pokemon:  pokeapi.Pokemon{Name: "defender"},
		hp:       100,
		maxHp:    100,
		isPlayer: false,
	}

	battle := &Battle{isActive: true}
	initialDefenderHP := defender.hp

	executeTurn(attacker, defender, battle)

	if defender.hp != initialDefenderHP {
		t.Error("Defender HP should not change if attacker is fainted")
	}
}
