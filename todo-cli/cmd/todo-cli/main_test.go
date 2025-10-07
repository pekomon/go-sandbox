package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type commandTestCase struct {
	name    string
	args    []string
	want    Command
	wantErr error
}

func TestParseCommand(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "tasks.json")

	cases := []commandTestCase{
		{
			name: "add task",
			args: []string{"todo-cli", "add", "write tests"},
			want: Command{Name: "add", Args: []string{"write tests"}},
		},
		{
			name: "list tasks",
			args: []string{"todo-cli", "list"},
			want: Command{Name: "list"},
		},
		{
			name: "mark done",
			args: []string{"todo-cli", "done", "2"},
			want: Command{Name: "done", Args: []string{"2"}},
		},
		{
			name: "remove task",
			args: []string{"todo-cli", "rm", "1"},
			want: Command{Name: "rm", Args: []string{"1"}},
		},
		{
			name: "clear tasks",
			args: []string{"todo-cli", "clear"},
			want: Command{Name: "clear"},
		},
		{
			name:    "unknown command",
			args:    []string{"todo-cli", "bogus"},
			wantErr: ErrUnknownCommand,
		},
		{
			name:    "done without id",
			args:    []string{"todo-cli", "done"},
			wantErr: ErrUnknownCommand,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = append([]string(nil), tc.args...)
			cmd, err := ParseCommand(os.Args)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
			if err != nil {
				return
			}
			if cmd.Name != tc.want.Name {
				t.Fatalf("expected command name %q, got %q", tc.want.Name, cmd.Name)
			}
			if fmt.Sprint(cmd.Args) != fmt.Sprint(tc.want.Args) {
				t.Fatalf("expected args %v, got %v", tc.want.Args, cmd.Args)
			}
		})
	}

	// ensure parse respects default storage path flag even when unused
	os.Args = []string{"todo-cli", "list", "--storage", tmpFile}
	cmd, err := ParseCommand(os.Args)
	if err != nil {
		t.Fatalf("unexpected error parsing storage flag: %v", err)
	}
	if cmd.Name != "list" {
		t.Fatalf("expected list command when parsing storage flag, got %q", cmd.Name)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != tmpFile {
		t.Fatalf("expected storage path in args, got %v", cmd.Args)
	}
}
