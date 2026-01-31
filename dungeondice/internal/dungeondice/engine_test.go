package dungeondice_test

import (
	"testing"

	"github.com/pekomon/go-sandbox/dungeondice/internal/dungeondice"
)

func TestResolveAction_Attack(t *testing.T) {
	attacker := dungeondice.Combatant{Name: "hero", HP: 10, MaxHP: 10, Defense: 0}
	defender := dungeondice.Combatant{Name: "slime", HP: 10, MaxHP: 10, Defense: 2}
	action := dungeondice.Action{Type: dungeondice.ActionAttack, Value: 5}

	gotAttacker, gotDefender, err := dungeondice.ResolveAction(attacker, defender, action)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAttacker != attacker {
		t.Fatalf("attacker should be unchanged: got %+v", gotAttacker)
	}
	if gotDefender.HP != 7 {
		t.Fatalf("unexpected defender HP: got %d want %d", gotDefender.HP, 7)
	}
	if gotDefender.Defense != 2 {
		t.Fatalf("defense should persist: got %d", gotDefender.Defense)
	}
}

func TestResolveAction_Defend(t *testing.T) {
	attacker := dungeondice.Combatant{Name: "hero", HP: 10, MaxHP: 10, Defense: 1}
	defender := dungeondice.Combatant{Name: "slime", HP: 10, MaxHP: 10, Defense: 0}
	action := dungeondice.Action{Type: dungeondice.ActionDefend, Value: 3}

	gotAttacker, gotDefender, err := dungeondice.ResolveAction(attacker, defender, action)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotDefender != defender {
		t.Fatalf("defender should be unchanged: got %+v", gotDefender)
	}
	if gotAttacker.Defense != 4 {
		t.Fatalf("unexpected defense: got %d want %d", gotAttacker.Defense, 4)
	}
}

func TestResolveAction_Heal(t *testing.T) {
	attacker := dungeondice.Combatant{Name: "hero", HP: 6, MaxHP: 10, Defense: 0}
	defender := dungeondice.Combatant{Name: "slime", HP: 10, MaxHP: 10, Defense: 0}
	action := dungeondice.Action{Type: dungeondice.ActionHeal, Value: 6}

	gotAttacker, gotDefender, err := dungeondice.ResolveAction(attacker, defender, action)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotDefender != defender {
		t.Fatalf("defender should be unchanged: got %+v", gotDefender)
	}
	if gotAttacker.HP != 10 {
		t.Fatalf("unexpected healed HP: got %d want %d", gotAttacker.HP, 10)
	}
}

func TestResolveAction_InvalidValue(t *testing.T) {
	attacker := dungeondice.Combatant{Name: "hero", HP: 10, MaxHP: 10}
	defender := dungeondice.Combatant{Name: "slime", HP: 10, MaxHP: 10}
	action := dungeondice.Action{Type: dungeondice.ActionAttack, Value: 0}

	_, _, err := dungeondice.ResolveAction(attacker, defender, action)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestResolveRound_TurnOrder(t *testing.T) {
	player := dungeondice.Combatant{Name: "hero", HP: 6, MaxHP: 10}
	enemy := dungeondice.Combatant{Name: "slime", HP: 6, MaxHP: 10}
	playerAction := dungeondice.Action{Type: dungeondice.ActionAttack, Value: 6}
	enemyAction := dungeondice.Action{Type: dungeondice.ActionAttack, Value: 6}

	gotPlayer, gotEnemy, err := dungeondice.ResolveRound(player, enemy, playerAction, enemyAction, 1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPlayer.HP != 0 {
		t.Fatalf("player should be defeated before acting: got %d", gotPlayer.HP)
	}
	if gotEnemy.HP != 6 {
		t.Fatalf("enemy should be unchanged: got %d", gotEnemy.HP)
	}
}

func TestAdvanceRun_Victory(t *testing.T) {
	run := dungeondice.Run{
		Rooms: []dungeondice.Room{{Enemy: dungeondice.Combatant{Name: "slime", HP: 0, MaxHP: 10}}},
		Position: 0,
		Player: dungeondice.Combatant{Name: "hero", HP: 5, MaxHP: 10},
		State: dungeondice.RunInProgress,
	}

	updated := dungeondice.AdvanceRun(run, run.Player, run.Rooms[0].Enemy)
	if updated.State != dungeondice.RunVictory {
		t.Fatalf("expected victory state, got %v", updated.State)
	}
	if updated.Position != 1 {
		t.Fatalf("expected position 1, got %d", updated.Position)
	}
	if !updated.Rooms[0].Cleared {
		t.Fatalf("expected room cleared")
	}
}

func TestAdvanceRun_Defeat(t *testing.T) {
	run := dungeondice.Run{
		Rooms: []dungeondice.Room{{Enemy: dungeondice.Combatant{Name: "slime", HP: 4, MaxHP: 10}}},
		Position: 0,
		Player: dungeondice.Combatant{Name: "hero", HP: 0, MaxHP: 10},
		State: dungeondice.RunInProgress,
	}

	updated := dungeondice.AdvanceRun(run, run.Player, run.Rooms[0].Enemy)
	if updated.State != dungeondice.RunDefeat {
		t.Fatalf("expected defeat state, got %v", updated.State)
	}
	if updated.Position != 0 {
		t.Fatalf("expected position 0, got %d", updated.Position)
	}
}

func TestAdvanceRun_InProgress(t *testing.T) {
	run := dungeondice.Run{
		Rooms: []dungeondice.Room{{Enemy: dungeondice.Combatant{Name: "slime", HP: 4, MaxHP: 10}}},
		Position: 0,
		Player: dungeondice.Combatant{Name: "hero", HP: 5, MaxHP: 10},
		State: dungeondice.RunInProgress,
	}

	updated := dungeondice.AdvanceRun(run, run.Player, run.Rooms[0].Enemy)
	if updated.State != dungeondice.RunInProgress {
		t.Fatalf("expected in-progress state, got %v", updated.State)
	}
	if updated.Position != 0 {
		t.Fatalf("expected position 0, got %d", updated.Position)
	}
}
