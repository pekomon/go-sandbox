package quiz

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

type Question struct {
	Prompt   string `json:"prompt"`
	Answer   string `json:"answer"`
	Category string `json:"category,omitempty"`
}

type Config struct {
	Seed  int64
	Count int
}

type Summary struct {
	Total   int
	Asked   int
	Correct int
}

func LoadQuestions(r io.Reader) ([]Question, error) {
	var questions []Question
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&questions); err != nil {
		return nil, fmt.Errorf("decode questions: %w", err)
	}

	for i, question := range questions {
		question.Prompt = strings.TrimSpace(question.Prompt)
		question.Answer = strings.TrimSpace(question.Answer)
		question.Category = strings.TrimSpace(question.Category)

		if question.Prompt == "" {
			return nil, fmt.Errorf("question %d: prompt is required", i)
		}
		if question.Answer == "" {
			return nil, fmt.Errorf("question %d: answer is required", i)
		}

		questions[i] = question
	}

	return questions, nil
}

func CheckAnswer(question Question, answer string) bool {
	return strings.EqualFold(strings.TrimSpace(question.Answer), strings.TrimSpace(answer))
}

func ShuffleQuestions(questions []Question, seed int64) []Question {
	shuffled := append([]Question(nil), questions...)
	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func Run(in io.Reader, out io.Writer, questions []Question, cfg Config) (Summary, error) {
	summary := Summary{Total: len(questions)}
	if len(questions) == 0 {
		return summary, nil
	}

	ordered := ShuffleQuestions(questions, cfg.Seed)
	limit := len(ordered)
	if cfg.Count > 0 && cfg.Count < limit {
		limit = cfg.Count
	}

	reader := bufio.NewReader(in)
	for i := 0; i < limit; i++ {
		question := ordered[i]
		if _, err := fmt.Fprintf(out, "%d. %s\n> ", i+1, question.Prompt); err != nil {
			return summary, fmt.Errorf("write prompt: %w", err)
		}

		answer, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return summary, fmt.Errorf("read answer: %w", err)
		}

		summary.Asked++
		if CheckAnswer(question, answer) {
			summary.Correct++
		}

		if err == io.EOF {
			break
		}
	}

	return summary, nil
}
