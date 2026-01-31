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
	return attacker, defender, errors.New("not implemented")
}

// ResolveRound resolves a round of combat with initiative ordering.
func ResolveRound(player, enemy Combatant, playerAction, enemyAction Action, playerInit, enemyInit int) (Combatant, Combatant, error) {
	return player, enemy, errors.New("not implemented")
}

// AdvanceRun updates the run based on the latest room outcome.
func AdvanceRun(run Run, player Combatant, enemy Combatant) Run {
	return run
}
