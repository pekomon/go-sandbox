package dungeondice

import "errors"

// ActionType defines the type of dice action taken in combat.
type ActionType int

const (
	ActionAttack ActionType = iota
	ActionDefend
	ActionHeal
)

// Action represents a resolved dice action.
type Action struct {
	Type  ActionType
	Value int
}

// Combatant is a combat participant in a run.
type Combatant struct {
	Name    string
	HP      int
	MaxHP   int
	Defense int
}

// RunState tracks run progression.
type RunState int

const (
	RunInProgress RunState = iota
	RunVictory
	RunDefeat
)

// Room represents a single encounter in a run.
type Room struct {
	Enemy   Combatant
	Cleared bool
}

// Run is the full run state.
type Run struct {
	Rooms    []Room
	Position int
	Player   Combatant
	State    RunState
}

// ResolveAction applies a single action from attacker to defender.
func ResolveAction(attacker, defender Combatant, action Action) (Combatant, Combatant, error) {
	if action.Value <= 0 {
		return attacker, defender, errors.New("action value must be positive")
	}

	switch action.Type {
	case ActionAttack:
		damage := action.Value - defender.Defense
		if damage < 0 {
			damage = 0
		}
		defender.HP -= damage
		if defender.HP < 0 {
			defender.HP = 0
		}
		return attacker, defender, nil
	case ActionDefend:
		attacker.Defense += action.Value
		return attacker, defender, nil
	case ActionHeal:
		attacker.HP += action.Value
		if attacker.HP > attacker.MaxHP {
			attacker.HP = attacker.MaxHP
		}
		return attacker, defender, nil
	default:
		return attacker, defender, errors.New("unknown action type")
	}
}

// ResolveRound resolves a round of combat with initiative ordering.
func ResolveRound(player, enemy Combatant, playerAction, enemyAction Action, playerInit, enemyInit int) (Combatant, Combatant, error) {
	var (
		updatedPlayer Combatant
		updatedEnemy  Combatant
		err           error
	)
	playerFirst := playerInit >= enemyInit
	if playerFirst {
		updatedPlayer, updatedEnemy, err = ResolveAction(player, enemy, playerAction)
		if err != nil {
			return player, enemy, err
		}
		player = updatedPlayer
		enemy = updatedEnemy
		if enemy.HP <= 0 {
			return player, enemy, nil
		}
		updatedEnemy, updatedPlayer, err = ResolveAction(enemy, player, enemyAction)
		if err != nil {
			return player, enemy, err
		}
		enemy = updatedEnemy
		player = updatedPlayer
		return player, enemy, nil
	}

	updatedEnemy, updatedPlayer, err = ResolveAction(enemy, player, enemyAction)
	if err != nil {
		return player, enemy, err
	}
	enemy = updatedEnemy
	player = updatedPlayer
	if player.HP <= 0 {
		return player, enemy, nil
	}
	updatedPlayer, updatedEnemy, err = ResolveAction(player, enemy, playerAction)
	if err != nil {
		return player, enemy, err
	}
	player = updatedPlayer
	enemy = updatedEnemy
	return player, enemy, nil
}

// AdvanceRun updates the run based on the latest room outcome.
func AdvanceRun(run Run, player Combatant, enemy Combatant) Run {
	run.Player = player
	if run.Position < 0 || run.Position >= len(run.Rooms) {
		return run
	}

	run.Rooms[run.Position].Enemy = enemy

	if player.HP <= 0 {
		run.State = RunDefeat
		return run
	}

	if enemy.HP <= 0 {
		run.Rooms[run.Position].Cleared = true
		run.Position++
		if run.Position >= len(run.Rooms) {
			run.State = RunVictory
		} else {
			run.State = RunInProgress
		}
		return run
	}

	run.State = RunInProgress
	return run
}
