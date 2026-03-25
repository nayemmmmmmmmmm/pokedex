package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

type ExplorationState struct {
	currentArea            string
	lastArea               string
	availableAreas         []AreaChoice
	pokemonFound           []string
	isExploring            bool
	lastExploreTime        time.Time
	inRandomEncounter      bool
	randomEncounterPokemon string
	randomEncounterArea    string
}

type AreaChoice struct {
	name        string
	description string
	direction   string
	url         string
	pokemon     []string
}

type Area struct {
	name        string
	description string
	neighbors   map[string]string // direction -> area name
	pokemon     []string
	isWild      bool
}

var gameWorld = map[string]Area{
	"pallet-town": {
		name:        "Pallet Town",
		description: "A peaceful town where new trainers begin their journey",
		neighbors:   map[string]string{"north": "route-1"},
		pokemon:     []string{"pidgey", "rattata", "pikachu"},
		isWild:      false,
	},
	"route-1": {
		name:        "Route 1",
		description: "A grassy path connecting Pallet Town to Viridian City",
		neighbors:   map[string]string{"south": "pallet-town", "north": "viridian-city"},
		pokemon:     []string{"pidgey", "rattata", "caterpie", "weedle"},
		isWild:      true,
	},
	"viridian-city": {
		name:        "Viridian City",
		description: "A bustling city with a Pokemon Gym",
		neighbors:   map[string]string{"south": "route-1", "north": "route-2", "west": "route-22"},
		pokemon:     []string{"pidgey", "rattata"},
		isWild:      false,
	},
	"route-2": {
		name:        "Route 2",
		description: "A path leading through dense forest",
		neighbors:   map[string]string{"south": "viridian-city", "north": "viridian-forest"},
		pokemon:     []string{"caterpie", "weedle", "metapod", "kakuna"},
		isWild:      true,
	},
	"viridian-forest": {
		name:        "Viridian Forest",
		description: "A mysterious forest filled with bug-type Pokemon",
		neighbors:   map[string]string{"south": "route-2", "north": "pewter-city"},
		pokemon:     []string{"caterpie", "weedle", "metapod", "kakuna", "pikachu"},
		isWild:      true,
	},
	"pewter-city": {
		name:        "Pewter City",
		description: "A city known for its rock-type Pokemon Gym",
		neighbors:   map[string]string{"south": "viridian-forest", "east": "route-3"},
		pokemon:     []string{"geodude", "sandshrew"},
		isWild:      false,
	},
	"route-3": {
		name:        "Route 3",
		description: "A mountainous path with many rock formations",
		neighbors:   map[string]string{"west": "pewter-city", "east": "mt-moon"},
		pokemon:     []string{"pidgey", "rattata", "spearow", "geodude", "sandshrew"},
		isWild:      true,
	},
	"mt-moon": {
		name:        "Mt. Moon",
		description: "A mysterious mountain cave filled with rare Pokemon",
		neighbors:   map[string]string{"west": "route-3", "east": "route-4"},
		pokemon:     []string{"zubat", "geodude", "clefairy", "paras"},
		isWild:      true,
	},
	"route-4": {
		name:        "Route 4",
		description: "A path descending from Mt. Moon",
		neighbors:   map[string]string{"west": "mt-moon", "east": "cerulean-city"},
		pokemon:     []string{"rattata", "spearow", "zubat"},
		isWild:      true,
	},
	"cerulean-city": {
		name:        "Cerulean City",
		description: "A beautiful city with a water-type Pokemon Gym",
		neighbors:   map[string]string{"west": "route-4", "north": "route-24", "south": "route-5", "east": "route-9"},
		pokemon:     []string{"goldeen", "magikarp", "staryu"},
		isWild:      false,
	},
}

func NewExplorationState() *ExplorationState {
	return &ExplorationState{
		currentArea:    "pallet-town",
		availableAreas: []AreaChoice{},
		pokemonFound:   []string{},
		isExploring:    false,
	}
}

func tryExploreArea(cfg *config, areaName string) error {
	state := cfg.explorationState

	// Convert area name to kebab-case and check if it's the current area
	normalizedArea := strings.ReplaceAll(strings.ToLower(areaName), " ", "-")
	noSpaceArea := strings.ReplaceAll(strings.ToLower(areaName), " ", "")

	if normalizedArea == state.currentArea || noSpaceArea == state.currentArea {
		return showCurrentArea(cfg)
	}

	// Check if it's a neighboring area
	currentArea := gameWorld[state.currentArea]
	for direction, neighbor := range currentArea.neighbors {
		if neighbor == normalizedArea || neighbor == noSpaceArea {
			return travelToArea(cfg, direction)
		}
	}

	return fmt.Errorf("unknown area or exploration command: %s", areaName)
}

func showCurrentArea(cfg *config) error {
	state := cfg.explorationState
	area := gameWorld[state.currentArea]

	fmt.Printf("\n=== %s ===\n", strings.Title(area.name))
	fmt.Printf("%s\n", area.description)
	fmt.Println()

	if area.isWild {
		fmt.Printf("🌲 This is a wild area - you may encounter Pokemon here!\n")
	} else {
		fmt.Printf("🏘️  This is a safe area - a good place to rest.\n")
	}

	fmt.Println("\n📍 Available directions:")
	for direction, neighbor := range area.neighbors {
		neighborArea := gameWorld[neighbor]
		fmt.Printf("  %s → %s\n", strings.Title(direction), strings.Title(neighborArea.name))
	}

	if len(state.pokemonFound) > 0 {
		fmt.Printf("\n🔍 Pokemon recently found in this area: %s\n", strings.Join(state.pokemonFound, ", "))
	}

	fmt.Println("\n🎯 Commands:")
	fmt.Println("  explore go <direction> - Travel in that direction")
	fmt.Println("  explore search - Look for Pokemon")
	fmt.Println("  explore map - Show world map")
	fmt.Println("  explore back - Return to previous area")

	return nil
}

func startExploration(cfg *config) error {
	state := cfg.explorationState
	state.isExploring = true
	state.lastExploreTime = time.Now()

	fmt.Println("\n🎒 You begin your Pokemon journey!")
	fmt.Println("You're starting in Pallet Town, ready to explore the world!")
	fmt.Println("Use 'explore look' to see your current location.")
	fmt.Println("Use 'explore go <direction>' to travel between areas.")
	fmt.Println("Use 'explore search' to look for Pokemon in wild areas.")

	return showCurrentArea(cfg)
}

func travelToArea(cfg *config, direction string) error {
	state := cfg.explorationState
	currentArea := gameWorld[state.currentArea]

	neighbor, exists := currentArea.neighbors[direction]
	if !exists {
		return fmt.Errorf("you cannot go %s from here", direction)
	}

	// Store previous area for back command
	state.lastArea = state.currentArea
	state.currentArea = neighbor
	state.pokemonFound = []string{} // Reset found Pokemon for new area

	fmt.Printf("\n🚶 Traveling %s...\n", direction)
	time.Sleep(1 * time.Second)

	// Check for random encounter while traveling (only if traveling to/from wild areas)
	shouldCheckEncounter := currentArea.isWild || gameWorld[neighbor].isWild
	if shouldCheckEncounter && rand.IntN(100) < 30 { // 30% chance of random encounter
		return handleRandomEncounter(cfg, state, state.currentArea, neighbor, direction)
	}

	newArea := gameWorld[neighbor]
	fmt.Printf("You arrived at %s!\n", strings.Title(newArea.name))

	return showCurrentArea(cfg)
}

func handleRandomEncounter(cfg *config, state *ExplorationState, fromAreaName, toAreaName, direction string) error {
	// Get the actual Area objects
	fromArea := gameWorld[fromAreaName]
	toArea := gameWorld[toAreaName]

	// Determine which area's Pokemon to use for encounter
	var encounterArea Area
	if fromArea.isWild {
		encounterArea = fromArea
	} else {
		encounterArea = toArea
	}

	// Select random Pokemon from the area
	pokemon := encounterArea.pokemon[rand.IntN(len(encounterArea.pokemon))]

	fmt.Printf("\n⚠️  A wild %s appeared while traveling!\n", strings.Title(pokemon))
	fmt.Printf("The %s blocks your path to %s!\n", strings.Title(pokemon), strings.Title(toArea.name))

	// 50% chance the Pokemon wants to battle
	if rand.IntN(100) < 50 {
		fmt.Printf("⚔️  The %s wants to battle! Use 'battle <your_pokemon> %s'\n", strings.Title(pokemon), pokemon)
		fmt.Printf("Or you can try to 'catch %s' if you're feeling lucky!\n", pokemon)
		fmt.Printf("Type 'explore continue' to try running away (50%% success rate)\n")

		// Store encounter state
		state.randomEncounterPokemon = pokemon
		state.randomEncounterArea = toAreaName
		state.inRandomEncounter = true
	} else {
		fmt.Printf("👀 The %s observes you cautiously before disappearing into the wild.\n", strings.Title(pokemon))
		fmt.Printf("You continue your journey to %s.\n", strings.Title(toArea.name))

		// Complete the travel after the encounter
		fmt.Printf("You arrived at %s!\n", strings.Title(toArea.name))
		return showCurrentArea(cfg)
	}

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if cfg.explorationState == nil {
		cfg.explorationState = NewExplorationState()
	}

	if len(args) == 0 {
		return showCurrentArea(cfg)
	}

	// Join all args for area names like "Route 4"
	fullCommand := strings.Join(args, " ")
	command := strings.ToLower(args[0])

	switch command {
	case "start", "begin":
		return startExploration(cfg)
	case "look", "area":
		return showCurrentArea(cfg)
	case "go", "move", "travel":
		if len(args) < 2 {
			return fmt.Errorf("usage: explore go <direction>")
		}
		return travelToArea(cfg, args[1])
	case "search", "find", "hunt":
		return searchForPokemon(cfg)
	case "back", "return":
		return returnToLastArea(cfg)
	case "map":
		return showWorldMap(cfg)
	case "continue":
		return handleRandomEncounterEscape(cfg)
	case "run", "flee":
		return handleRandomEncounterEscape(cfg)
	default:
		// Try to treat it as an area name
		return tryExploreArea(cfg, fullCommand)
	}
}

func handleRandomEncounterEscape(cfg *config) error {
	state := cfg.explorationState

	if !state.inRandomEncounter {
		return fmt.Errorf("you're not in a random encounter")
	}

	if rand.IntN(100) < 50 { // 50% chance to escape
		fmt.Printf("\n🏃 You successfully ran away from the wild %s!\n", strings.Title(state.randomEncounterPokemon))
		fmt.Printf("You continue your journey to %s.\n", strings.Title(state.randomEncounterArea))

		// Complete the travel
		newArea := gameWorld[state.randomEncounterArea]
		fmt.Printf("You arrived at %s!\n", strings.Title(newArea.name))

		// Clear encounter state
		state.inRandomEncounter = false
		state.randomEncounterPokemon = ""
		state.randomEncounterArea = ""

		return showCurrentArea(cfg)
	} else {
		fmt.Printf("\n❌ You couldn't escape from the wild %s!\n", strings.Title(state.randomEncounterPokemon))
		fmt.Printf("The %s is still blocking your path. You'll need to deal with it.\n", strings.Title(state.randomEncounterPokemon))
		fmt.Printf("Options: 'battle <your_pokemon> %s', 'catch %s', or try 'explore run' again\n",
			state.randomEncounterPokemon, state.randomEncounterPokemon)
		return nil
	}
}

func searchForPokemon(cfg *config) error {
	state := cfg.explorationState
	area := gameWorld[state.currentArea]

	if !area.isWild {
		return fmt.Errorf("you can only search for Pokemon in wild areas")
	}

	// Add cooldown to prevent spamming
	if time.Since(state.lastExploreTime) < 5*time.Second {
		remaining := 5*time.Second - time.Since(state.lastExploreTime)
		return fmt.Errorf("please wait %v before searching again", remaining.Round(time.Second))
	}

	state.lastExploreTime = time.Now()

	fmt.Printf("\n🔍 Searching for wild Pokemon in %s...\n", strings.Title(area.name))
	time.Sleep(2 * time.Second)

	// 70% chance to find Pokemon
	if rand.IntN(100) < 70 {
		// Select random Pokemon from area
		pokemon := area.pokemon[rand.IntN(len(area.pokemon))]
		state.pokemonFound = append(state.pokemonFound, pokemon)

		fmt.Printf("✨ You found a wild %s!\n", strings.Title(pokemon))

		// 30% chance it's ready to battle
		if rand.IntN(100) < 30 {
			fmt.Printf("⚔️  The %s wants to battle! Use 'battle <your_pokemon> %s'\n", strings.Title(pokemon), pokemon)
		} else {
			fmt.Printf("👀 The %s observed you curiously before disappearing into the wild.\n", strings.Title(pokemon))
		}
	} else {
		fmt.Println("🍃 You didn't find any Pokemon this time. Try again!")
	}

	return nil
}

func returnToLastArea(cfg *config) error {
	state := cfg.explorationState
	if state.lastArea == "" {
		return fmt.Errorf("you haven't traveled anywhere yet")
	}

	current := state.currentArea
	state.currentArea = state.lastArea
	state.lastArea = current
	state.pokemonFound = []string{}

	fmt.Printf("\n🔙 Returning to %s...\n", strings.Title(gameWorld[state.currentArea].name))
	return showCurrentArea(cfg)
}

func showWorldMap(cfg *config) error {
	fmt.Println("\n🗺️  WORLD MAP")
	fmt.Println("=" + strings.Repeat("=", 40))

	visited := make(map[string]bool)
	if cfg.explorationState != nil {
		// Mark visited areas (simplified - just show current area)
		visited[cfg.explorationState.currentArea] = true
	}

	areas := []string{
		"pallet-town", "route-1", "viridian-city", "route-2", "viridian-forest",
		"pewter-city", "route-3", "mt-moon", "route-4", "cerulean-city",
	}

	for _, areaKey := range areas {
		area := gameWorld[areaKey]
		if visited[areaKey] {
			fmt.Printf("✅ %s - %s\n", strings.Title(area.name), area.description)
		} else {
			fmt.Printf("❓ %s - ???\n", strings.Title(area.name))
		}
	}

	fmt.Println("\n💡 Tip: Explore areas to discover them and find new Pokemon!")
	return nil
}

func commandAreas(cfg *config, args ...string) error {
	fmt.Println("\n🗺️  AVAILABLE AREAS")
	fmt.Println("=" + strings.Repeat("=", 30))

	for _, areaKey := range []string{"pallet-town", "route-1", "viridian-city", "route-2", "viridian-forest", "pewter-city"} {
		area := gameWorld[areaKey]
		fmt.Printf("📍 %s\n", strings.Title(area.name))
		fmt.Printf("   %s\n", area.description)
		fmt.Printf("   Pokemon: %s\n", strings.Join(area.pokemon, ", "))
		fmt.Printf("   Directions: ")
		directions := []string{}
		for dir := range area.neighbors {
			directions = append(directions, dir)
		}
		fmt.Printf("%s\n", strings.Join(directions, ", "))
		fmt.Println()
	}

	return nil
}
