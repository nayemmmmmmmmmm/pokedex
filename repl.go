package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
	"github.com/nayemmmmmmmmmm/pokedex/internal/pokeapi"
)

type config struct {
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
	caughtPokemon    map[string]pokeapi.Pokemon
	history          []string
	historyIndex     int
	party            *Party
	evolutionTracker *EvolutionTracker
}

func startRepl(cfg *config) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "Pokedex > ",
		HistoryFile:     "pokedex_history.tmp",
		AutoComplete:    nil,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}
			fmt.Println("Error reading input:", err)
			continue
		}

		words := cleanInput(line)
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(cfg, args...)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"party": {
			name:        "party",
			description: "View your Pokemon party with levels and HP",
			callback:    commandParty,
		},
		"evolve": {
			name:        "evolve <pokemon_name>",
			description: "Start evolving a Pokemon (requires minimum level)",
			callback:    commandEvolve,
		},
		"evolutionstatus": {
			name:        "evolutionstatus [pokemon_name]",
			description: "Check evolution status of your Pokemon",
			callback:    commandEvolutionStatus,
		},
		"checkevolutions": {
			name:        "checkevolutions",
			description: "Check for completed evolutions",
			callback:    commandCheckEvolutions,
		},
		"cancelevolve": {
			name:        "cancelevolve <pokemon_name>",
			description: "Cancel an ongoing evolution",
			callback:    commandCancelEvolution,
		},
		"addtoparty": {
			name:        "addtoparty <pokemon_name>",
			description: "Add a caught Pokemon to your party",
			callback:    commandAddToParty,
		},
		"removefromparty": {
			name:        "removefromparty <pokemon_name>",
			description: "Remove a Pokemon from your party",
			callback:    commandRemoveFromParty,
		},
		"heal": {
			name:        "heal [pokemon_name]",
			description: "Heal your Pokemon (all or specific)",
			callback:    commandHeal,
		},
		"battle": {
			name:        "battle <your_pokemon> <opponent_pokemon>",
			description: "Simulate a battle between your pokemon and an opponent",
			callback:    commandBattle,
		},
		"catch": {
			name:        "catch <pokemon_name",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon_name",
			description: "View details about a caught Pokemon",
			callback:    commandInspect,
		},
		"explore": {
			name:        "explore <location_name",
			description: "Explore a location",
			callback:    commandExplore,
		},
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapb,
		},
		"pokedex": {
			name:        "pokedex",
			description: "See all the pokemon you've caught",
			callback:    commandPokedex,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}
