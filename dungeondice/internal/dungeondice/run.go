package dungeondice

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// ErrInvalidRooms indicates a run length that is not supported.
var ErrInvalidRooms = errors.New("rooms must be positive")

// RunConfig defines how to build and simulate a run.
type RunConfig struct {
	Class string
	Seed  int64
	Rooms int
}

// RunSummary captures the outcome of a simulated run.
type RunSummary struct {
	Class      string
	Seed       int64
	Rooms      int
	Cleared    int
	State      RunState
	FinalHP    int
	FinalMaxHP int
	Rounds     int
}

const maxDefense = 6

// SimulateRun builds and resolves a full run based on config.
func SimulateRun(cfg RunConfig) (RunSummary, error) {
	if cfg.Rooms <= 0 {
		return RunSummary{}, ErrInvalidRooms
	}

	class, err := ClassByName(cfg.Class)
	if err != nil {
		return RunSummary{}, err
	}

	seed := cfg.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed))

	run := buildRun(class, cfg.Rooms, rng)
	summary := RunSummary{
		Class:      class.Name,
		Seed:       seed,
		Rooms:      cfg.Rooms,
		FinalMaxHP: class.MaxHP,
	}

	rounds := 0
	for run.State == RunInProgress {
		room := run.Rooms[run.Position]
		player := run.Player
		enemy := room.Enemy

		for player.HP > 0 && enemy.HP > 0 {
			playerAction := chooseAction(rng, player, class.Bonus)
			enemyAction := chooseAction(rng, enemy, ActionBonus{})
			playerInit := rollDie(rng, 6)
			enemyInit := rollDie(rng, 6)

			nextPlayer, nextEnemy, err := ResolveRound(player, enemy, playerAction, enemyAction, playerInit, enemyInit)
			if err != nil {
				return RunSummary{}, err
			}
			player = nextPlayer
			enemy = nextEnemy
			player = clampDefense(player)
			enemy = clampDefense(enemy)
			rounds++
		}

		run = AdvanceRun(run, player, enemy)
	}

	summary.State = run.State
	summary.FinalHP = run.Player.HP
	summary.Cleared = countCleared(run.Rooms)
	summary.Rounds = rounds
	return summary, nil
}

func buildRun(class Class, rooms int, rng *rand.Rand) Run {
	run := Run{
		Rooms:    make([]Room, rooms),
		Position: 0,
		Player: Combatant{
			Name:    class.Name,
			HP:      class.MaxHP,
			MaxHP:   class.MaxHP,
			Defense: class.Defense,
		},
		State: RunInProgress,
	}

	for i := 0; i < rooms; i++ {
		run.Rooms[i] = Room{
			Enemy: generateEnemy(rng, i),
		}
	}

	return run
}

func generateEnemy(rng *rand.Rand, index int) Combatant {
	baseHP := 6 + rng.Intn(4) + index
	defense := rng.Intn(3)
	name := fmt.Sprintf("foe-%d", index+1)
	return Combatant{
		Name:    name,
		HP:      baseHP,
		MaxHP:   baseHP,
		Defense: defense,
	}
}

func chooseAction(rng *rand.Rand, self Combatant, bonus ActionBonus) Action {
	actionType := pickActionType(rng, self)
	value := rollActionValue(rng, actionType, bonus)
	return Action{Type: actionType, Value: value}
}

func pickActionType(rng *rand.Rand, self Combatant) ActionType {
	lowHP := self.HP > 0 && self.HP <= self.MaxHP/3
	roll := rng.Intn(100)
	if lowHP {
		switch {
		case roll < 50:
			return ActionHeal
		case roll < 80:
			return ActionAttack
		default:
			return ActionDefend
		}
	}

	switch {
	case roll < 60:
		return ActionAttack
	case roll < 85:
		return ActionDefend
	default:
		return ActionHeal
	}
}

func rollActionValue(rng *rand.Rand, actionType ActionType, bonus ActionBonus) int {
	base := 1
	value := 1
	switch actionType {
	case ActionAttack:
		base = rollDie(rng, 6)
		value = base + bonus.Attack
	case ActionDefend:
		base = rollDie(rng, 4)
		value = base + bonus.Defense
	case ActionHeal:
		base = rollDie(rng, 4)
		value = base + bonus.Heal
	default:
		value = base
	}

	if value < 1 {
		value = 1
	}
	return value
}

func rollDie(rng *rand.Rand, sides int) int {
	if sides <= 1 {
		return 1
	}
	return rng.Intn(sides) + 1
}

func countCleared(rooms []Room) int {
	count := 0
	for _, room := range rooms {
		if room.Cleared {
			count++
		}
	}
	return count
}

func clampDefense(combatant Combatant) Combatant {
	if combatant.Defense > maxDefense {
		combatant.Defense = maxDefense
	}
	if combatant.Defense < 0 {
		combatant.Defense = 0
	}
	return combatant
}
