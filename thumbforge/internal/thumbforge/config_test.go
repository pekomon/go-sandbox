package thumbforge_test

import (
	"testing"

	"github.com/pekomon/go-sandbox/thumbforge/internal/thumbforge"
)

func TestParseSize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     thumbforge.Size
		wantErr  bool
	}{
		{
			name:  "valid",
			input: "320x240",
			want: thumbforge.Size{Width: 320, Height: 240},
		},
		{
			name:    "missing_separator",
			input:   "320-240",
			wantErr: true,
		},
		{
			name:    "zero_dimension",
			input:   "0x240",
			wantErr: true,
		},
		{
			name:    "non_numeric",
			input:   "axb",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := thumbforge.ParseSize(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("unexpected size: got %+v want %+v", got, tt.want)
			}
		})
	}
}
