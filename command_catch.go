package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}

	name := args[0]
	pokemon, err := cfg.pokeapiClient.GetPokemon(name)
	if err != nil {
		return err
	}

	res := rand.IntN(pokemon.BaseExperience)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if res > 40 {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemon.Name)

	// Add to caught Pokemon (for compatibility)
	cfg.caughtPokemon[pokemon.Name] = pokemon

	// Try to add to party
	err = cfg.party.AddPokemon(pokemon)
	if err != nil {
		fmt.Printf("Couldn't add %s to party: %v\n", pokemon.Name, err)
		fmt.Printf("%s was sent to storage.\n", pokemon.Name)
	} else {
		fmt.Printf("%s was added to your party!\n", pokemon.Name)
	}

	return nil
}
