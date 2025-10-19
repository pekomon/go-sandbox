package game

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	Max      int   // default: 100
	Attempts int   // default: 7
	Seed     int64 // 0 => non-deterministic (time-based)
}

type StatsStore interface {
	Load() (Stats, error)
	Save(Stats) error
}

type Stats struct {
	Games          int     `json:"games"`
	Wins           int     `json:"wins"`
	TotalGuesses   int     `json:"total_guesses"`
	AverageGuesses float64 `json:"average_guesses"`
}

// Run plays one game reading guesses (one per line) from r and writing prompts/results to w.
// Behavior:
//   - target in [1..Max]
//   - print "higher"/"lower" hints
//   - on success: print "Correct!"
//   - on attempts exhausted: print "Game over! The number was X."
//   - non-integer input: print "Invalid input" and DO NOT consume an attempt
//   - stats: increment Games, Wins (if success), accumulate TotalGuesses and AverageGuesses
func Run(r io.Reader, w io.Writer, opts Options, stats StatsStore) error {
	max := opts.Max
	if max <= 0 {
		max = 100
	}
	attempts := opts.Attempts
	if attempts <= 0 {
		attempts = 7
	}

	// RNG (deterministic if Seed!=0)
	var src rand.Source
	if opts.Seed != 0 {
		src = rand.NewSource(opts.Seed)
	} else {
		src = rand.NewSource(time.Now().UnixNano())
	}
	rng := rand.New(src)
	target := rng.Intn(max) + 1

	// Load stats (ignore error by treating as zeroed – the stats tests cover persistence separately)
	var st Stats
	if stats != nil {
		if loaded, err := stats.Load(); err == nil {
			st = loaded
		}
	}

	sc := bufio.NewScanner(r)
	usedGuesses := 0

	fmt.Fprintln(w, "Guess the number!")
	for left := attempts; left > 0; {
		fmt.Fprintf(w, "Attempts left: %d\n", left)
		if !sc.Scan() {
			// EOF or read error => stop the game loop
			break
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			fmt.Fprintln(w, "Invalid input")
			return fmt.Errorf("invalid input")
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintln(w, "Invalid input")
			return fmt.Errorf("invalid input: %w", err)
		}

		usedGuesses++
		left--

		if n == target {
			fmt.Fprintln(w, "Correct!")
			// Update stats on success and persist
			if stats != nil {
				st.Games++
				st.Wins++
				st.TotalGuesses += usedGuesses
				if st.Games > 0 {
					st.AverageGuesses = float64(st.TotalGuesses) / float64(st.Games)
				}
				_ = stats.Save(st)
			}
			return nil
		}
		if n < target {
			fmt.Fprintln(w, "higher")
		} else {
			fmt.Fprintln(w, "lower")
		}
	}

	// Out of attempts (or input ended)
	fmt.Fprintf(w, "Game over! The number was %d.\n", target)
	if stats != nil {
		st.Games++
		st.TotalGuesses += usedGuesses
		if st.Games > 0 {
			st.AverageGuesses = float64(st.TotalGuesses) / float64(st.Games)
		}
		_ = stats.Save(st)
	}
	return nil
}
