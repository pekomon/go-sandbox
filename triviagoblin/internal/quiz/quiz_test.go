package quiz_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pekomon/go-sandbox/triviagoblin/internal/quiz"
)

func TestLoadQuestions_Valid(t *testing.T) {
	input := `[
		{"prompt": "Capital of France?", "answer": "Paris", "category": "geo"},
		{"prompt": "2+2", "answer": "4"}
	]`

	got, err := quiz.LoadQuestions(strings.NewReader(input))
	if err != nil {
		t.Fatalf("LoadQuestions() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("LoadQuestions() length = %d, want 2", len(got))
	}
	if got[0].Prompt != "Capital of France?" || got[0].Answer != "Paris" || got[0].Category != "geo" {
		t.Fatalf("unexpected first question: %+v", got[0])
	}
	if got[1].Prompt != "2+2" || got[1].Answer != "4" || got[1].Category != "" {
		t.Fatalf("unexpected second question: %+v", got[1])
	}
}

func TestLoadQuestions_InvalidJSON(t *testing.T) {
	input := `[{` // malformed JSON

	_, err := quiz.LoadQuestions(strings.NewReader(input))
	if err == nil {
		t.Fatalf("LoadQuestions() error = nil, want non-nil")
	}
}

func TestLoadQuestions_MissingFields(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "missing prompt",
			input: `[{"answer": "Paris"}]`,
		},
		{
			name:  "missing answer",
			input: `[{"prompt": "Capital?"}]`,
		},
		{
			name:  "blank prompt",
			input: `[{"prompt": " ", "answer": "Paris"}]`,
		},
		{
			name:  "blank answer",
			input: `[{"prompt": "Capital?", "answer": " "}]`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := quiz.LoadQuestions(strings.NewReader(tc.input))
			if err == nil {
				t.Fatalf("LoadQuestions() error = nil, want non-nil")
			}
		})
	}
}

func TestCheckAnswer_TrimAndCase(t *testing.T) {
	question := quiz.Question{Prompt: "Capital?", Answer: "Paris"}
	if !quiz.CheckAnswer(question, "  paris  ") {
		t.Fatalf("CheckAnswer() = false, want true")
	}
	if quiz.CheckAnswer(question, "Lyon") {
		t.Fatalf("CheckAnswer() = true, want false")
	}
}

func TestShuffleQuestions_Deterministic(t *testing.T) {
	questions := []quiz.Question{
		{Prompt: "A", Answer: "ok"},
		{Prompt: "B", Answer: "ok"},
		{Prompt: "C", Answer: "ok"},
		{Prompt: "D", Answer: "ok"},
	}

	first := quiz.ShuffleQuestions(questions, 42)
	second := quiz.ShuffleQuestions(questions, 42)
	if len(first) != len(second) {
		t.Fatalf("ShuffleQuestions() length mismatch: %d vs %d", len(first), len(second))
	}

	for i := range first {
		if first[i].Prompt != second[i].Prompt {
			t.Fatalf("ShuffleQuestions() order mismatch at %d: %q vs %q", i, first[i].Prompt, second[i].Prompt)
		}
	}
}

func TestRunSummaryCounts(t *testing.T) {
	questions := []quiz.Question{
		{Prompt: "Q1", Answer: "ok"},
		{Prompt: "Q2", Answer: "ok"},
		{Prompt: "Q3", Answer: "ok"},
	}
	input := bytes.NewBufferString("ok\nok\n")
	output := &bytes.Buffer{}

	summary, err := quiz.Run(input, output, questions, quiz.Config{
		Seed:  99,
		Count: 2,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if summary.Total != 3 {
		t.Fatalf("summary.Total = %d, want 3", summary.Total)
	}
	if summary.Asked != 2 {
		t.Fatalf("summary.Asked = %d, want 2", summary.Asked)
	}
	if summary.Correct != 2 {
		t.Fatalf("summary.Correct = %d, want 2", summary.Correct)
	}
}
