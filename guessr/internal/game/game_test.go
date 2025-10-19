package game_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/pekomon/go-sandbox/guessr/internal/game"
)

type fakeStats struct {
	load    game.Stats
	saved   game.Stats
	saveErr error
}

func (f *fakeStats) Load() (game.Stats, error) { return f.load, nil }
func (f *fakeStats) Save(s game.Stats) error {
	f.saved = s
	return f.saveErr
}

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		opts         game.Options
		prepare      func(target, max int, opts game.Options) []int
		wantSubstr   []string
		wantStats    game.Stats
		wantErr      bool
		errContains  string
		prependInput []string
	}{
		{
			name: "win within attempts",
			opts: game.Options{Max: 10, Attempts: 5, Seed: 42},
			prepare: func(target, max int, opts game.Options) []int {
				return hintGuesses(target, max)
			},
			wantSubstr: []string{"higher", "lower", "correct"},
			wantStats:  game.Stats{Games: 1, Wins: 1, TotalGuesses: 3, AverageGuesses: 3},
		},
		{
			name: "out of attempts",
			opts: game.Options{Max: 10, Attempts: 3, Seed: 99},
			prepare: func(target, max int, opts game.Options) []int {
				return losingGuesses(target, max, opts.Attempts)
			},
			wantSubstr: []string{"game over"},
			wantStats:  game.Stats{Games: 1, Wins: 0, TotalGuesses: 3, AverageGuesses: 3},
		},
		{
			name:         "invalid input",
			opts:         game.Options{Seed: 7},
			wantErr:      true,
			errContains:  "invalid",
			prependInput: []string{"oops"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := tt.opts
			max := opts.Max
			if max == 0 {
				max = 100
			}
			if opts.Attempts == 0 {
				opts.Attempts = 7
			}

			guesses := make([]int, 0)
			if tt.prepare != nil {
				target := expectedTarget(opts.Seed, max)
				guesses = append(guesses, tt.prepare(target, max, opts)...)
			}

			inputs := make([]string, 0, len(tt.prependInput)+len(guesses))
			inputs = append(inputs, tt.prependInput...)
			for _, g := range guesses {
				inputs = append(inputs, fmt.Sprintf("%d", g))
			}

			out, fs, err := runGame(t, opts, strings.Join(inputs, "\n")+"\n")
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Run() error = nil, want non-nil")
				}
				if tt.errContains != "" && !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.errContains)) {
					t.Fatalf("Run() error = %q, want contains %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("Run() error = %v, want nil", err)
			}

			lowerOut := strings.ToLower(out)
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(lowerOut, substr) {
					t.Errorf("output %q does not contain %q", out, substr)
				}
			}

			if tt.wantStats != (game.Stats{}) {
				if fs.saved != tt.wantStats {
					t.Errorf("saved stats = %+v, want %+v", fs.saved, tt.wantStats)
				}
			}
		})
	}
}

func runGame(t *testing.T, opts game.Options, input string) (string, *fakeStats, error) {
	t.Helper()

	in := bytes.NewBufferString(input)
	out := &bytes.Buffer{}
	stats := &fakeStats{}

	err := game.Run(in, out, opts, stats)
	return out.String(), stats, err
}

func expectedTarget(seed int64, max int) int {
	if seed == 0 {
		seed = 1
	}
	r := rand.New(rand.NewSource(seed))
	return r.Intn(max) + 1
}

func hintGuesses(target, max int) []int {
	guesses := make([]int, 0, 3)
	if target > 1 {
		guesses = append(guesses, target-1)
	} else {
		guesses = append(guesses, target+1)
	}
	if target < max {
		guesses = append(guesses, target+1)
	} else {
		guesses = append(guesses, target-1)
	}
	guesses = append(guesses, target)
	return guesses
}

func losingGuesses(target, max, attempts int) []int {
	guesses := make([]int, 0, attempts)
	current := 1
	for len(guesses) < attempts {
		if current == target {
			current++
			if current > max {
				current = 1
			}
			continue
		}
		guesses = append(guesses, current)
		current++
		if current > max {
			current = 1
		}
	}
	return guesses
}
