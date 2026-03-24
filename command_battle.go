package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/nayemmmmmmmmmm/pokedex/internal/pokeapi"
)

type BattlePokemon struct {
	pokemon  pokeapi.Pokemon
	hp       int
	maxHp    int
	isPlayer bool
}

type Battle struct {
	playerPokemon   *BattlePokemon
	opponentPokemon *BattlePokemon
	turnCount       int
	isActive        bool
}

func commandBattle(cfg *config, args ...string) error {
	if len(args) != 2 {
		return errors.New("usage: battle <your_pokemon> <opponent_pokemon>")
	}

	playerName := args[0]
	opponentName := args[1]

	playerPokemon, ok := cfg.caughtPokemon[playerName]
	if !ok {
		return fmt.Errorf("you haven't caught %s", playerName)
	}

	opponentPokemon, err := cfg.pokeapiClient.GetPokemon(opponentName)
	if err != nil {
		return fmt.Errorf("couldn't find opponent %s: %v", opponentName, err)
	}

	battle := &Battle{
		playerPokemon: &BattlePokemon{
			pokemon:  playerPokemon,
			maxHp:    calculateHP(playerPokemon),
			hp:       calculateHP(playerPokemon),
			isPlayer: true,
		},
		opponentPokemon: &BattlePokemon{
			pokemon:  opponentPokemon,
			maxHp:    calculateHP(opponentPokemon),
			hp:       calculateHP(opponentPokemon),
			isPlayer: false,
		},
		turnCount: 0,
		isActive:  true,
	}

	return runBattle(cfg, battle)
}

func calculateHP(pokemon pokeapi.Pokemon) int {
	for _, stat := range pokemon.Stats {
		if stat.Stat.Name == "hp" {
			return stat.BaseStat
		}
	}
	return 50
}

func calculateAttack(pokemon pokeapi.Pokemon) int {
	for _, stat := range pokemon.Stats {
		if stat.Stat.Name == "attack" {
			return stat.BaseStat
		}
	}
	return 50
}

func calculateDefense(pokemon pokeapi.Pokemon) int {
	for _, stat := range pokemon.Stats {
		if stat.Stat.Name == "defense" {
			return stat.BaseStat
		}
	}
	return 50
}

func calculateSpeed(pokemon pokeapi.Pokemon) int {
	for _, stat := range pokemon.Stats {
		if stat.Stat.Name == "speed" {
			return stat.BaseStat
		}
	}
	return 50
}

func runBattle(cfg *config, battle *Battle) error {
	fmt.Printf("\n=== BATTLE START ===\n")
	fmt.Printf("%s vs %s\n", battle.playerPokemon.pokemon.Name, battle.opponentPokemon.pokemon.Name)
	fmt.Printf("%s HP: %d/%d\n", battle.playerPokemon.pokemon.Name, battle.playerPokemon.hp, battle.playerPokemon.maxHp)
	fmt.Printf("%s HP: %d/%d\n", battle.opponentPokemon.pokemon.Name, battle.opponentPokemon.hp, battle.opponentPokemon.maxHp)
	fmt.Println("==================")

	for battle.isActive {
		battle.turnCount++
		fmt.Printf("\n--- Turn %d ---\n", battle.turnCount)

		first, second := determineTurnOrder(battle.playerPokemon, battle.opponentPokemon)

		executeTurn(first, second)
		
		if !battle.isActive {
			break
		}

		executeTurn(second, first)

		fmt.Printf("%s HP: %d/%d\n", battle.playerPokemon.pokemon.Name, battle.playerPokemon.hp, battle.playerPokemon.maxHp)
		fmt.Printf("%s HP: %d/%d\n", battle.opponentPokemon.pokemon.Name, battle.opponentPokemon.hp, battle.opponentPokemon.maxHp)

		time.Sleep(1 * time.Second)
	}

	fmt.Printf("\n=== BATTLE END ===\n")
	return nil
}

func determineTurnOrder(attacker, defender *BattlePokemon) (*BattlePokemon, *BattlePokemon) {
	attackerSpeed := calculateSpeed(attacker.pokemon)
	defenderSpeed := calculateSpeed(defender.pokemon)

	if attackerSpeed >= defenderSpeed {
		return attacker, defender
	}
	return defender, attacker
}

func executeTurn(attacker, defender *BattlePokemon) {
	if !attacker.isAlive() || !defender.isAlive() {
		return
	}

	damage := calculateDamage(attacker, defender)
	defender.hp -= damage

	if defender.hp < 0 {
		defender.hp = 0
	}

	fmt.Printf("%s attacks %s for %d damage!\n", attacker.pokemon.Name, defender.pokemon.Name, damage)

	if !defender.isAlive() {
		fmt.Printf("%s has fainted!\n", defender.pokemon.Name)
		if defender.isPlayer {
			fmt.Printf("You lose! %s wins the battle.\n", attacker.pokemon.Name)
		} else {
			fmt.Printf("You win! %s defeated %s.\n", attacker.pokemon.Name, defender.pokemon.Name)
		}
	}
}

func calculateDamage(attacker, defender *BattlePokemon) int {
	attack := calculateAttack(attacker.pokemon)
	defense := calculateDefense(defender.pokemon)
	level := 50

	baseDamage := ((2*level + 10) * attack / defense) / 50 + 2
	randomFactor := rand.IntN(38) + 217
	damage := baseDamage * randomFactor / 255

	typeEffectiveness := getTypeEffectiveness(attacker.pokemon, defender.pokemon)
	damage = int(float64(damage) * typeEffectiveness)

	if typeEffectiveness > 1.0 {
		fmt.Printf("It's super effective! ")
	} else if typeEffectiveness < 1.0 && typeEffectiveness > 0 {
		fmt.Printf("It's not very effective... ")
	} else if typeEffectiveness == 0 {
		fmt.Printf("It has no effect... ")
	}

	return max(1, damage)
}

func getTypeEffectiveness(attacker, defender pokeapi.Pokemon) float64 {
	if len(attacker.Types) == 0 || len(defender.Types) == 0 {
		return 1.0
	}

	attackerType := strings.ToLower(attacker.Types[0].Type.Name)
	defenderTypes := make([]string, len(defender.Types))
	for i, t := range defender.Types {
		defenderTypes[i] = strings.ToLower(t.Type.Name)
	}

	effectiveness := 1.0

	for _, defType := range defenderTypes {
		typeMultiplier := getTypeMultiplier(attackerType, defType)
		effectiveness *= typeMultiplier
	}

	return effectiveness
}

func getTypeMultiplier(attackerType, defenderType string) float64 {
	typeChart := map[string]map[string]float64{
		"normal": {
			"rock":     0.5,
			"ghost":    0,
			"steel":    0.5,
		},
		"fire": {
			"fire":   0.5,
			"water":  0.5,
			"grass":  2,
			"ice":    2,
			"bug":    0.5,
			"rock":   0.5,
			"dragon": 0.5,
			"steel":  2,
		},
		"water": {
			"fire":   2,
			"water":  0.5,
			"grass":  0.5,
			"ground": 2,
			"rock":   2,
			"dragon": 0.5,
		},
		"electric": {
			"water":   2,
			"electric": 0.5,
			"grass":   0.5,
			"ground":  0,
			"flying":  2,
			"dragon":  0.5,
		},
		"grass": {
			"fire":   0.5,
			"water":  2,
			"grass":  0.5,
			"poison":  0.5,
			"flying": 0.5,
			"bug":    0.5,
			"rock":   2,
			"dragon": 0.5,
			"steel":  0.5,
		},
		"ice": {
			"fire":   0.5,
			"water":  0.5,
			"grass":  2,
			"ice":    0.5,
			"ground": 2,
			"flying": 2,
			"dragon": 2,
			"steel":  0.5,
		},
		"fighting": {
			"normal": 2,
			"ice":    2,
			"poison": 0.5,
			"flying": 0.5,
			"psychic": 0.5,
			"bug":    0.5,
			"rock":   2,
			"ghost":  0,
			"dark":   2,
			"steel":  2,
			"fairy":  0.5,
		},
		"poison": {
			"grass":  2,
			"poison": 0.5,
			"ground": 0.5,
			"bug":    2,
			"rock":   0.5,
			"ghost":  0.5,
			"steel":  0,
			"fairy":  2,
		},
		"ground": {
			"fire":    2,
			"grass":   0.5,
			"electric": 2,
			"poison":  2,
			"flying":  0,
			"bug":    0.5,
			"rock":   2,
			"ghost":  0,
			"steel":  2,
		},
		"flying": {
			"grass":  2,
			"electric": 0.5,
			"poison": 2,
			"fighting": 2,
			"bug":    2,
			"rock":   0.5,
			"steel":  0.5,
		},
		"psychic": {
			"fighting": 2,
			"poison":  2,
			"ground":  2,
			"psychic": 0.5,
			"bug":    0.5,
			"ghost":  0,
			"dark":   0,
			"steel":  0.5,
		},
		"bug": {
			"fire":    0.5,
			"grass":   2,
			"fighting": 0.5,
			"poison":  0.5,
			"flying":  0.5,
			"psychic": 2,
			"ghost":  0.5,
			"dark":    2,
			"steel":   0.5,
			"fairy":   0.5,
		},
		"rock": {
			"fire":    2,
			"ice":     2,
			"fighting": 0.5,
			"ground":  2,
			"flying":  2,
			"bug":     2,
			"steel":   0.5,
		},
		"ghost": {
			"normal":  0,
			"psychic": 2,
			"ghost":   2,
			"bug":     0.5,
			"poison":  0.5,
			"dark":    0.5,
			"fairy":   0,
		},
		"dragon": {
			"dragon": 2,
			"steel":  0.5,
			"fairy":  0,
		},
		"steel": {
			"fire":   0.5,
			"water":  0.5,
			"ice":    2,
			"steel":  0.5,
			"fairy":  2,
		},
		"fairy": {
			"fire":    0.5,
			"poison":  0.5,
			"fighting": 2,
			"dragon":  2,
			"dark":    2,
			"steel":   0.5,
		},
	}

	if attackerChart, exists := typeChart[attackerType]; exists {
		if multiplier, exists := attackerChart[defenderType]; exists {
			return multiplier
		}
	}

	return 1.0
}

func (bp *BattlePokemon) isAlive() bool {
	return bp.hp > 0
}
