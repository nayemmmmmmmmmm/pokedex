package main

import (
	"fmt"
	"strings"
	"time"
)

func commandEvolve(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: evolve <pokemon_name>")
	}

	name := args[0]
	partyPokemon, err := cfg.party.GetPokemon(name)
	if err != nil {
		return fmt.Errorf("you don't have %s in your party", name)
	}

	if !partyPokemon.CanEvolve() {
		return fmt.Errorf("%s cannot evolve yet", name)
	}

	target, canEvolve := partyPokemon.GetEvolutionTarget()
	if !canEvolve {
		return fmt.Errorf("%s has no evolution path", name)
	}

	// Check if evolution is already in progress
	_, _, inProgress := cfg.evolutionTracker.GetEvolutionStatus(name)
	if inProgress {
		return fmt.Errorf("%s is already evolving", name)
	}

	// Start evolution process
	chain := evolutionChains[strings.ToLower(partyPokemon.pokemon.Name)]
	timeNeeded := chain["timeNeeded"].(time.Duration)
	
	cfg.evolutionTracker.StartEvolution(name, target, partyPokemon.level, timeNeeded)
	
	fmt.Printf("%s started evolving to %s!\n", name, target)
	fmt.Printf("Evolution will complete in %v\n", timeNeeded)
	
	return nil
}

func commandEvolutionStatus(cfg *config, args ...string) error {
	if len(args) == 0 {
		// Show all evolution statuses
		if len(cfg.evolutionTracker.evolutions) == 0 {
			fmt.Println("No Pokemon are currently evolving.")
			return nil
		}
		
		fmt.Println("=== EVOLUTION STATUS ===")
		for name, evolution := range cfg.evolutionTracker.evolutions {
			timeRemaining := time.Until(evolution.evolveTime)
			if timeRemaining < 0 {
				timeRemaining = 0
			}
			fmt.Printf("%s → %s: %v remaining\n", 
				strings.Title(name), 
				strings.Title(evolution.toPokemon), 
				timeRemaining.Round(time.Second))
		}
		return nil
	}

	name := args[0]
	target, timeRemaining, exists := cfg.evolutionTracker.GetEvolutionStatus(name)
	if !exists {
		return fmt.Errorf("%s is not evolving", name)
	}

	if timeRemaining <= 0 {
		fmt.Printf("%s is ready to evolve to %s!\n", name, target)
	} else {
		fmt.Printf("%s will evolve to %s in %v\n", name, target, timeRemaining.Round(time.Second))
	}

	return nil
}

func commandCheckEvolutions(cfg *config, args ...string) error {
	evolved := cfg.evolutionTracker.CheckEvolutions()
	
	if len(evolved) == 0 {
		fmt.Println("No evolutions completed.")
		return nil
	}

	fmt.Println("=== EVOLUTIONS COMPLETED ===")
	for _, evolvedName := range evolved {
		// Find the party Pokemon that evolved
		for name, partyPokemon := range cfg.party.pokemon {
			target, canEvolve := partyPokemon.GetEvolutionTarget()
			if canEvolve && strings.ToLower(target) == strings.ToLower(evolvedName) {
				fmt.Printf("\n🎉 %s is evolving! 🎉\n", strings.Title(name))
				
				// Perform the evolution
				err := partyPokemon.Evolve(evolvedName)
				if err != nil {
					fmt.Printf("Error evolving %s: %v\n", name, err)
					continue
				}
				
				fmt.Printf("%s evolved into %s!\n", strings.Title(name), strings.Title(evolvedName))
				fmt.Printf("Stats increased and HP fully restored!\n")
				
				// Update the party map key
				delete(cfg.party.pokemon, name)
				cfg.party.pokemon[strings.ToLower(evolvedName)] = partyPokemon
				
				break
			}
		}
	}

	return nil
}

func commandCancelEvolution(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: cancelevolve <pokemon_name>")
	}

	name := args[0]
	
	// Check if evolution exists
	_, _, exists := cfg.evolutionTracker.GetEvolutionStatus(name)
	if !exists {
		return fmt.Errorf("%s is not evolving", name)
	}

	cfg.evolutionTracker.CancelEvolution(name)
	fmt.Printf("Cancelled evolution for %s\n", name)
	
	return nil
}
