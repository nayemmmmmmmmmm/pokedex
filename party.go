package main

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/nayemmmmmmmmmm/pokedex/internal/pokeapi"
)

type PartyPokemon struct {
	pokemon pokeapi.Pokemon
	level   int
	exp     int
	hp      int
	maxHp   int
}

type Party struct {
	pokemon map[string]*PartyPokemon
	maxSize int
}

func NewParty(maxSize int) *Party {
	return &Party{
		pokemon: make(map[string]*PartyPokemon),
		maxSize: maxSize,
	}
}

func (p *Party) AddPokemon(pokemon pokeapi.Pokemon) error {
	if len(p.pokemon) >= p.maxSize {
		return fmt.Errorf("party is full (max %d pokemon)", p.maxSize)
	}

	name := strings.ToLower(pokemon.Name)
	if _, exists := p.pokemon[name]; exists {
		return fmt.Errorf("%s is already in your party", pokemon.Name)
	}

	partyPokemon := &PartyPokemon{
		pokemon: pokemon,
		level:   5, // Start at level 5
		exp:     0,
		hp:      calculateHP(pokemon),
		maxHp:   calculateHP(pokemon),
	}

	p.pokemon[name] = partyPokemon
	return nil
}

func (p *Party) RemovePokemon(name string) error {
	lowerName := strings.ToLower(name)
	if _, exists := p.pokemon[lowerName]; !exists {
		return fmt.Errorf("%s is not in your party", name)
	}

	delete(p.pokemon, lowerName)
	return nil
}

func (p *Party) GetPokemon(name string) (*PartyPokemon, error) {
	lowerName := strings.ToLower(name)
	pokemon, exists := p.pokemon[lowerName]
	if !exists {
		return nil, fmt.Errorf("%s is not in your party", name)
	}
	return pokemon, nil
}

func (p *Party) ListPokemon() []string {
	var names []string
	for name := range p.pokemon {
		names = append(names, name)
	}
	return names
}

func (p *Party) IsFull() bool {
	return len(p.pokemon) >= p.maxSize
}

func (pp *PartyPokemon) GainExp(amount int) {
	pp.exp += amount
	
	// Check for level up (simple formula: 100 exp per level)
	expNeeded := pp.level * 100
	for pp.exp >= expNeeded {
		pp.exp -= expNeeded
		pp.levelUp()
		expNeeded = pp.level * 100
	}
}

func (pp *PartyPokemon) levelUp() {
	pp.level++
	
	// Increase stats on level up
	hpIncrease := rand.IntN(5) + 3 // 3-7 HP increase
	pp.maxHp += hpIncrease
	pp.hp = pp.maxHp // Full heal on level up
	
	// Update base stats in the Pokemon struct
	for i, stat := range pp.pokemon.Stats {
		switch stat.Stat.Name {
		case "hp":
			pp.pokemon.Stats[i].BaseStat += hpIncrease
		case "attack":
			pp.pokemon.Stats[i].BaseStat += rand.IntN(3) + 1 // 1-3 attack increase
		case "defense":
			pp.pokemon.Stats[i].BaseStat += rand.IntN(3) + 1 // 1-3 defense increase
		case "speed":
			pp.pokemon.Stats[i].BaseStat += rand.IntN(2) + 1 // 1-2 speed increase
		}
	}
	
	fmt.Printf("\n🎉 %s reached level %d! 🎉\n", pp.pokemon.Name, pp.level)
	fmt.Printf("HP increased by %d (max HP: %d)\n", hpIncrease, pp.maxHp)
}

func (pp *PartyPokemon) IsAlive() bool {
	return pp.hp > 0
}

func (pp *PartyPokemon) TakeDamage(damage int) {
	pp.hp -= damage
	if pp.hp < 0 {
		pp.hp = 0
	}
}

func (pp *PartyPokemon) Heal(amount int) {
	pp.hp += amount
	if pp.hp > pp.maxHp {
		pp.hp = pp.maxHp
	}
}

func (pp *PartyPokemon) FullHeal() {
	pp.hp = pp.maxHp
}
