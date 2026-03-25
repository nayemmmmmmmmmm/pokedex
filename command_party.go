package main

import (
	"fmt"
	"strings"
	"time"
)

func commandParty(cfg *config, args ...string) error {
	// Check for completed evolutions first
	cfg.evolutionTracker.CheckEvolutions()

	if len(cfg.party.pokemon) == 0 {
		fmt.Println("Your party is empty. Catch some Pokemon first!")
		return nil
	}

	fmt.Println("=== YOUR PARTY ===")
	fmt.Printf("Party: %d/%d Pokemon\n", len(cfg.party.pokemon), cfg.party.maxSize)
	fmt.Println()

	for name, partyPokemon := range cfg.party.pokemon {
		fmt.Printf("🎖️  %s (Lv.%d)\n", strings.Title(name), partyPokemon.level)
		fmt.Printf("   HP: %d/%d\n", partyPokemon.hp, partyPokemon.maxHp)
		fmt.Printf("   EXP: %d/%d\n", partyPokemon.exp, partyPokemon.level*100)
		fmt.Printf("   Types: ")
		for i, t := range partyPokemon.pokemon.Types {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s", strings.Title(t.Type.Name))
		}
		fmt.Println()

		// Show evolution status
		if target, timeRemaining, inProgress := cfg.evolutionTracker.GetEvolutionStatus(name); inProgress {
			if timeRemaining <= 0 {
				fmt.Printf("   🔄 Ready to evolve to %s! Use 'checkevolutions'\n", strings.Title(target))
			} else {
				fmt.Printf("   ⏳ Evolving to %s in %v\n", strings.Title(target), timeRemaining.Round(time.Second))
			}
		} else if partyPokemon.CanEvolve() {
			target, _ := partyPokemon.GetEvolutionTarget()
			fmt.Printf("   ✨ Can evolve to %s (Lv.%d+)\n", strings.Title(target), evolutionChains[strings.ToLower(name)]["minLevel"].(int))
		}
		fmt.Println()
	}

	return nil
}

func commandAddToParty(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: addtoparty <pokemon_name>")
	}

	name := args[0]
	pokemon, exists := cfg.caughtPokemon[name]
	if !exists {
		return fmt.Errorf("you haven't caught %s", name)
	}

	err := cfg.party.AddPokemon(pokemon)
	if err != nil {
		return err
	}

	fmt.Printf("%s was added to your party!\n", name)
	return nil
}

func commandRemoveFromParty(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: removefromparty <pokemon_name>")
	}

	name := args[0]
	err := cfg.party.RemovePokemon(name)
	if err != nil {
		return err
	}

	fmt.Printf("%s was removed from your party.\n", name)
	return nil
}

func commandHeal(cfg *config, args ...string) error {
	if len(cfg.party.pokemon) == 0 {
		return fmt.Errorf("your party is empty")
	}

	if len(args) == 0 {
		// Heal all Pokemon in party
		fmt.Println("Healing all Pokemon in your party...")
		for name, partyPokemon := range cfg.party.pokemon {
			beforeHP := partyPokemon.hp
			partyPokemon.FullHeal()
			fmt.Printf("%s: %d → %d HP\n", strings.Title(name), beforeHP, partyPokemon.hp)
		}
		fmt.Println("All Pokemon have been fully healed!")
	} else {
		// Heal specific Pokemon
		name := args[0]
		partyPokemon, err := cfg.party.GetPokemon(name)
		if err != nil {
			return err
		}

		beforeHP := partyPokemon.hp
		partyPokemon.FullHeal()
		fmt.Printf("%s: %d → %d HP\n", strings.Title(name), beforeHP, partyPokemon.hp)
	}

	return nil
}
