package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type menuTestCase struct {
	name               string
	script             string
	wantExit           int
	wantStdout         []string
	wantStdoutInOrder  []string
	wantStdoutMinCount map[string]int
	wantStderr         []string
}

func TestMenuScenarios(t *testing.T) {
	baseMenu := []string{
		"TODO CLI MENU",
		"1) Add task",
		"2) List tasks",
		"3) Mark done",
		"4) Remove task",
		"5) Clear tasks",
		"0) Exit",
		"Select an option:",
	}

	cases := []menuTestCase{
		{
			name:     "1 exit via zero",
			script:   "0\n",
			wantExit: 0,
			wantStdout: []string{
				"Goodbye",
			},
		},
		{
			name:     "1 exit via exit command",
			script:   "exit\n",
			wantExit: 0,
			wantStdout: []string{
				"Goodbye",
			},
		},
		{
			name:     "2 add then list",
			script:   "1\nRead Go docs\n2\n0\n",
			wantExit: 0,
			wantStdout: []string{
				"Enter task description:",
				"added #1",
				"[ ] #1 Read Go docs",
			},
		},
		{
			name:     "3 add twice done first then list",
			script:   "1\nFirst task\n1\nSecond task\n3\n1\n2\n0\n",
			wantExit: 0,
			wantStdout: []string{
				"added #1",
				"added #2",
				"Enter task ID to mark done:",
				"done #1",
			},
			wantStdoutInOrder: []string{
				"[ ] #2 Second task",
				"[x] #1 First task",
			},
		},
		{
			name:     "4 add remove list empty",
			script:   "1\nDisposable task\n4\n1\n2\n0\n",
			wantExit: 0,
			wantStdout: []string{
				"added #1",
				"Enter task ID to remove:",
				"removed #1",
				"No tasks found.",
			},
		},
		{
			name:     "5 clear removes all",
			script:   "1\nAlpha\n1\nBeta\n5\n2\n0\n",
			wantExit: 0,
			wantStdout: []string{
				"added #1",
				"added #2",
				"cleared",
				"No tasks found.",
			},
		},
		{
			name:     "6 invalid selection re-prompts",
			script:   "9\n0\n",
			wantExit: 0,
			wantStdout: []string{
				"Goodbye",
			},
			wantStdoutMinCount: map[string]int{
				"Select an option:": 2,
			},
			wantStderr: []string{
				"invalid selection",
			},
		},
	}

	entrypoints := []struct {
		name string
		args []string
		env  map[string]string
	}{
		{
			name: "menu command",
			args: []string{"menu"},
		},
		{
			name: "env trigger",
			env: map[string]string{
				"TODO_CLI_MENU": "1",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for _, ep := range entrypoints {
				ep := ep
				t.Run(ep.name, func(t *testing.T) {
					storageDir := t.TempDir()
					t.Setenv("TODO_CLI_PATH", filepath.Join(storageDir, "tasks.json"))
					for k, v := range ep.env {
						t.Setenv(k, v)
					}

					stdout, stderr, exit := runMenuHarness(t, ep.args, tc.script)
					if exit != tc.wantExit {
						t.Fatalf("expected exit %d, got %d. stdout=%q stderr=%q", tc.wantExit, exit, stdout, stderr)
					}

					requireContainsAll(t, stdout, baseMenu)
					requireContainsAll(t, stdout, tc.wantStdout)
					requireContainsAll(t, stderr, tc.wantStderr)
					if len(tc.wantStdoutInOrder) > 0 {
						requireContainsInOrder(t, stdout, tc.wantStdoutInOrder)
					}
					if len(tc.wantStdoutMinCount) > 0 {
						requireCountAtLeast(t, stdout, tc.wantStdoutMinCount)
					}
				})
			}
		})
	}
}

func runMenuHarness(t *testing.T, args []string, input string) (string, string, int) {
	t.Helper()

	inR, inW, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdin pipe: %v", err)
	}
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stderr pipe: %v", err)
	}

	defer func() {
		inR.Close()
		outR.Close()
		errR.Close()
	}()

	originalStdin := os.Stdin
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	os.Stdin = inR
	os.Stdout = outW
	os.Stderr = errW

	defer func() {
		os.Stdin = originalStdin
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}()

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(&stdoutBuf, outR)
		close(stdoutDone)
	}()

	stderrDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(&stderrBuf, errR)
		close(stderrDone)
	}()

	inputDone := make(chan struct{})
	go func() {
		if input != "" {
			_, _ = io.WriteString(inW, input)
		}
		inW.Close()
		close(inputDone)
	}()

	exit := Run(args)

	outW.Close()
	errW.Close()

	<-stdoutDone
	<-stderrDone
	<-inputDone

	return stdoutBuf.String(), stderrBuf.String(), exit
}

func requireContainsAll(t *testing.T, output string, want []string) {
	t.Helper()
	for _, w := range want {
		if w == "" {
			continue
		}
		if !strings.Contains(output, w) {
			t.Fatalf("expected output to contain %q, got %q", w, output)
		}
	}
}

func requireContainsInOrder(t *testing.T, output string, want []string) {
	t.Helper()
	searchStart := 0
	for _, w := range want {
		idx := strings.Index(output[searchStart:], w)
		if idx < 0 {
			t.Fatalf("expected output to contain %q after position %d, got %q", w, searchStart, output)
		}
		searchStart += idx + len(w)
	}
}

func requireCountAtLeast(t *testing.T, output string, wants map[string]int) {
	t.Helper()
	for w, min := range wants {
		if count := strings.Count(output, w); count < min {
			t.Fatalf("expected output to contain %q at least %d times, got %d. output=%q", w, min, count, output)
		}
	}
}
