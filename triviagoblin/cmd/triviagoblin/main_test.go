package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunQuizCommand_Success(t *testing.T) {
	questionsPath := filepath.Join(t.TempDir(), "questions.json")
	writeQuestionsFile(t, questionsPath, `[
		{"prompt":"Capital of France?","answer":"Paris","category":"geo"},
		{"prompt":"2+2","answer":"4","category":"math"},
		{"prompt":"Capital of Spain?","answer":"Madrid","category":"geo"}
	]`)

	in := strings.NewReader("Paris\nMadrid\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exit := run([]string{"run", "--file", questionsPath, "--count", "2", "--seed", "7", "--category", "geo"}, in, &stdout, &stderr)
	if exit != 0 {
		t.Fatalf("run() exit = %d, want 0, stderr=%q", exit, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}

	out := stdout.String()
	for _, want := range []string{
		"Capital of France?",
		"Capital of Spain?",
		"Run summary",
		"Questions available: 2",
		"Questions asked: 2",
		"Correct answers: 2",
		"Seed: 7",
		"Category: geo",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout missing %q:\n%s", want, out)
		}
	}
}

func TestRunQuizCommand_MissingFileFlag(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exit := run([]string{"run"}, strings.NewReader(""), &stdout, &stderr)
	if exit != 2 {
		t.Fatalf("run() exit = %d, want 2", exit)
	}
	if !strings.Contains(stderr.String(), "missing --file") {
		t.Fatalf("stderr = %q, want missing --file", stderr.String())
	}
}

func TestRunQuizCommand_NoQuestionsForCategory(t *testing.T) {
	questionsPath := filepath.Join(t.TempDir(), "questions.json")
	writeQuestionsFile(t, questionsPath, `[
		{"prompt":"Capital of France?","answer":"Paris","category":"geo"}
	]`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exit := run([]string{"run", "--file", questionsPath, "--category", "math"}, strings.NewReader(""), &stdout, &stderr)
	if exit != 1 {
		t.Fatalf("run() exit = %d, want 1", exit)
	}
	if !strings.Contains(stderr.String(), `no questions matched category "math"`) {
		t.Fatalf("stderr = %q, want category mismatch error", stderr.String())
	}
}

func TestRunQuizCommand_InvalidCount(t *testing.T) {
	questionsPath := filepath.Join(t.TempDir(), "questions.json")
	writeQuestionsFile(t, questionsPath, `[
		{"prompt":"Capital of France?","answer":"Paris","category":"geo"}
	]`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exit := run([]string{"run", "--file", questionsPath, "--count", "-1"}, strings.NewReader(""), &stdout, &stderr)
	if exit != 2 {
		t.Fatalf("run() exit = %d, want 2", exit)
	}
	if !strings.Contains(stderr.String(), "count must be >= 0") {
		t.Fatalf("stderr = %q, want invalid count error", stderr.String())
	}
}

func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exit := run([]string{"help"}, strings.NewReader(""), &stdout, &stderr)
	if exit != 0 {
		t.Fatalf("run() exit = %d, want 0", exit)
	}
	if !strings.Contains(stdout.String(), "usage: triviagoblin run") {
		t.Fatalf("stdout = %q, want usage", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func writeQuestionsFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}
