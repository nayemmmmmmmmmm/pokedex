# 🎮 Pokedex CLI Game

A command-line Pokemon game built in Go that allows you to catch, battle, train, and explore the world of Pokemon with your own party!

## 📋 Table of Contents

- [Features](#-features)
- [Installation](#-installation)
- [Getting Started](#-getting-started)
- [Commands](#-commands)
- [Game Systems](#-game-systems)
- [Project Structure](#-project-structure)
- [API Integration](#-api-integration)
- [Testing](#-testing)
- [Contributing](#-contributing)

## ✨ Features

### 🎯 Core Gameplay
- **Catch Pokemon** - Encounter and catch wild Pokemon using Pokeballs
- **Battle System** - Turn-based battles with type effectiveness and damage calculation
- **Party Management** - Build and manage a team of up to 6 Pokemon
- **Evolution System** - Time-based evolution with stat boosts
- **World Exploration** - Navigate through interconnected areas with random encounters

### 🎮 Advanced Features
- **Leveling System** - Gain experience and level up Pokemon through battles
- **Type Effectiveness** - Realistic Pokemon type matchups (super effective, not very effective, etc.)
- **Random Encounters** - Wild Pokemon appear while traveling between areas
- **Healing System** - Restore your Pokemon's health
- **Command History** - Navigate through previous commands with arrow keys

## 🚀 Installation

### Prerequisites
- Go 1.21 or higher
- Internet connection (for PokeAPI integration)

### Build from Source
```bash
git clone https://github.com/nayemmmmmmmmmm/pokedex.git
cd pokedex
go build
```

### Run the Game
```bash
./pokedex
```

## 🎮 Getting Started

1. **Start your journey** by exploring areas and catching Pokemon
2. **Build your party** with up to 6 Pokemon
3. **Train your Pokemon** through battles to gain experience
4. **Evolve your Pokemon** by meeting level requirements and waiting for evolution time
5. **Explore the world** from Pallet Town to Cerulean City

## 📚 Commands

### 🗺️ Exploration Commands
```bash
explore                    # Show current area
explore start             # Begin your journey
explore go <direction>    # Travel north/south/east/west
explore search            # Look for Pokemon in wild areas
explore map               # Show world map
explore back              # Return to previous area
areas                     # Show all available areas
```

### 🎯 Pokemon Commands
```bash
catch <pokemon>           # Attempt to catch a wild Pokemon
battle <your_pokemon> <opponent>  # Battle against another Pokemon
party                     # View your Pokemon party
addtoparty <pokemon>      # Add caught Pokemon to party
removefromparty <pokemon> # Remove Pokemon from party
heal [pokemon]            # Heal your Pokemon (all or specific)
```

### ⚡ Evolution Commands
```bash
evolve <pokemon>          # Start evolution process
evolutionstatus [pokemon] # Check evolution progress
checkevolutions           # Complete ready evolutions
cancelevolve <pokemon>    # Cancel ongoing evolution
```

### 📋 Information Commands
```bash
pokedex                   # View all caught Pokemon
inspect <pokemon>         # View detailed Pokemon information
help                      # Show available commands
exit                      # Exit the game
```

### 🗺️ Navigation Commands
```bash
map                       # Show next page of locations
mapb                      # Show previous page of locations
```

## 🎮 Game Systems

### ⚔️ Battle System
- **Turn-based combat** with speed-based turn order
- **Type effectiveness** calculations (2x, 0.5x, 0x damage)
- **Damage calculation** based on Attack vs Defense stats
- **Experience gain** for winning battles
- **Fainting mechanics** - Pokemon with 0 HP can't battle

### 🎯 Party System
- **Maximum 6 Pokemon** per party
- **Level and experience tracking** (100 EXP per level)
- **HP management** with healing capabilities
- **Stat increases** on level up (HP, Attack, Defense, Speed)
- **Full heal** on evolution

### ⚡ Evolution System
- **Time-based evolution** - Takes real time to complete
- **Level requirements** - Each Pokemon has minimum level for evolution
- **Stat boosts** - 20-30% stat increases on evolution
- **Evolution chains** - Multiple evolution stages for many Pokemon
- **Random Eevee evolutions** - Eevee can evolve into Flareon, Jolteon, or Vaporeon

### 🗺️ World Exploration
- **10 interconnected areas** from Pallet Town to Cerulean City
- **Area-specific Pokemon** - Each area has unique Pokemon populations
- **Random encounters** - 30% chance of wild Pokemon while traveling
- **Wild vs Safe areas** - Only wild areas have Pokemon encounters
- **Navigation system** - Use compass directions to travel

### 🎲 Random Encounters
- **Travel encounters** - Wild Pokemon appear while traveling
- **Multiple options** - Battle, catch, or run away
- **Escape mechanics** - 50% success rate for running away
- **Area-specific encounters** - Pokemon match the area's population

## 📁 Project Structure

```
pokedex/
├── main.go                    # Application entry point
├── repl.go                    # REPL interface and command registry
├── party.go                   # Party system and Pokemon management
├── evolution.go               # Evolution mechanics and chains
├── exploration.go             # World exploration and random encounters
├── command_battle.go          # Battle command implementation
├── command_catch.go           # Catch command implementation
├── command_party.go           # Party management commands
├── command_evolution.go       # Evolution commands
├── command_inspect.go         # Pokemon inspection command
├── command_pokedex.go         # Pokedex command
├── command_explore.go         # Legacy explore command (removed)
├── internal/
│   └── pokeapi/
│       ├── client.go          # API client
│       ├── types_pokemon.go   # Pokemon data structures
│       └── types_locations.go # Location data structures
├── *_test.go                  # Unit tests for all systems
└── README.md                  # This file
```

## 🔌 API Integration

This game integrates with the **PokeAPI** (pokeapi.co) to provide:
- **Real Pokemon data** - Stats, types, evolution chains
- **Location information** - Area details and Pokemon encounters
- **Comprehensive Pokedex** - All Pokemon species data

### API Features Used
- Pokemon endpoint: `/api/v2/pokemon/{id}/`
- Location endpoint: `/api/v2/location/{id}/`
- Type effectiveness calculations
- Base experience for EXP calculations

## 🧪 Testing

The project includes comprehensive unit tests for all major systems:

```bash
# Run all tests
go test -v

# Run specific test packages
go test -v ./...  # All tests
go test -v command_battle_test.go  # Battle tests
go test -v party_test.go  # Party tests
go test -v evolution_test.go  # Evolution tests
```

### Test Coverage
- ✅ Battle system mechanics
- ✅ Party management
- ✅ Evolution system
- ✅ Type effectiveness
- ✅ Damage calculation
- ✅ Exploration mechanics
- ✅ Random encounters

## 🎯 Game World

### Available Areas
1. **Pallet Town** - Starting area (Pidgey, Rattata, Pikachu)
2. **Route 1** - Grass path (Pidgey, Rattata, Caterpie, Weedle)
3. **Viridian City** - City with Gym (Pidgey, Rattata)
4. **Route 2** - Forest path (Caterpie, Weedle, Metapod, Kakuna)
5. **Viridian Forest** - Bug-type paradise (Caterpie, Weedle, Metapod, Kakuna, Pikachu)
6. **Pewter City** - Rock-type city (Geodude, Sandshrew)
7. **Route 3** - Mountain path (Pidgey, Rattata, Spearow, Geodude, Sandshrew)
8. **Mt. Moon** - Cave area (Zubat, Geodude, Clefairy, Paras)
9. **Route 4** - Descending path (Rattata, Spearow, Zubat)
10. **Cerulean City** - Water city (Goldeen, Magikarp, Staryu)

### Evolution Chains
- **Pidgey → Pidgeotto → Pidgeot**
- **Charmander → Charmeleon → Charizard**
- **Bulbasaur → Ivysaur → Venusaur**
- **Squirtle → Wartortle → Blastoise**
- **Caterpie → Metapod → Butterfree**
- **Weedle → Kakuna → Beedrill**
- **Pikachu → Raichu**
- **Eevee** → (Flareon/Jolteon/Vaporeon - random)

## 🎮 Tips & Strategies

### For Beginners
1. **Start by catching** a few different Pokemon to build your party
2. **Explore Route 1** early for diverse Pokemon types
3. **Level up your Pokemon** through battles before attempting evolutions
4. **Keep your Pokemon healed** - fainted Pokemon can't battle

### Advanced Strategies
1. **Type advantages** - Use super effective attacks in battles
2. **Evolution timing** - Plan evolutions around important battles
3. **Party diversity** - Maintain different types for various situations
4. **Random encounters** - Be prepared for wild Pokemon while traveling

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup
```bash
git clone https://github.com/nayemmmmmmmmmm/pokedex.git
cd pokedex
go mod tidy
go build
go test -v
```

## 📄 License

This project is for educational purposes. Pokemon and related characters are trademarks of Nintendo/Game Freak.

## 🔗 Links

- [PokeAPI](https://pokeapi.co/) - The Pokemon API
- [Go Documentation](https://golang.org/doc/) - Go programming language
- [Original Repository](https://github.com/nayemmmmmmmmmm/pokedex) - Project source code

---

**Enjoy your Pokemon adventure! 🎮✨**
