package dungeondice_test

import (
	"errors"
	"testing"

	"github.com/pekomon/go-sandbox/dungeondice/internal/dungeondice"
)

func TestClassByName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "warrior", want: "Warrior"},
		{name: "ROGUE", want: "Rogue"},
		{name: " Mystic ", want: "Mystic"},
	}

	for _, tc := range tests {
		got, err := dungeondice.ClassByName(tc.name)
		if err != nil {
			t.Fatalf("ClassByName(%q) unexpected error: %v", tc.name, err)
		}
		if got.Name != tc.want {
			t.Fatalf("ClassByName(%q) name mismatch: got %q want %q", tc.name, got.Name, tc.want)
		}
	}
}

func TestClassByName_Invalid(t *testing.T) {
	_, err := dungeondice.ClassByName("paladin")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, dungeondice.ErrUnknownClass) {
		t.Fatalf("expected ErrUnknownClass, got %v", err)
	}
}

func TestSimulateRun(t *testing.T) {
	summary, err := dungeondice.SimulateRun(dungeondice.RunConfig{
		Class: "warrior",
		Seed:  42,
		Rooms: 3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Class != "Warrior" {
		t.Fatalf("unexpected class: got %q", summary.Class)
	}
	if summary.Seed != 42 {
		t.Fatalf("unexpected seed: got %d", summary.Seed)
	}
	if summary.Rooms != 3 {
		t.Fatalf("unexpected room count: got %d", summary.Rooms)
	}
	if summary.Cleared < 0 || summary.Cleared > summary.Rooms {
		t.Fatalf("cleared out of range: got %d", summary.Cleared)
	}
	if summary.State == dungeondice.RunInProgress {
		t.Fatalf("run should be resolved")
	}
	if summary.FinalHP < 0 || summary.FinalHP > summary.FinalMaxHP {
		t.Fatalf("final HP out of range: got %d", summary.FinalHP)
	}
	if summary.Rounds == 0 {
		t.Fatalf("expected at least one round")
	}
	if summary.State == dungeondice.RunVictory && summary.Cleared != summary.Rooms {
		t.Fatalf("victory should clear all rooms: got %d", summary.Cleared)
	}
	if summary.State == dungeondice.RunDefeat && summary.Cleared >= summary.Rooms {
		t.Fatalf("defeat should clear fewer rooms: got %d", summary.Cleared)
	}
}

func TestSimulateRun_InvalidRooms(t *testing.T) {
	_, err := dungeondice.SimulateRun(dungeondice.RunConfig{
		Class: "warrior",
		Rooms: 0,
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, dungeondice.ErrInvalidRooms) {
		t.Fatalf("expected ErrInvalidRooms, got %v", err)
	}
}
