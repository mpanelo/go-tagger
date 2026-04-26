package structfind_test

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/mpanelo/go-tagger/internal/config"
	"github.com/mpanelo/go-tagger/internal/structfind"
)

func TestFind(t *testing.T) {
	src := `package main
	type Example struct {
		A string
		B int /* target */
	}
	func main() {}`

	tests := []struct {
		name          string
		src           string
		line          string
		offset        int
		structName    string
		expectedStart int
		expectedEnd   int
		expectedError bool
	}{
		{
			name:          "Sad Case: config fields not set",
			expectedError: true,
		},
		{
			name:          "Sad Case: struct not found",
			src:           src,
			structName:    "eXamPLE",
			expectedError: true,
		},
		{
			name:          "Happy Case: struct found",
			src:           src,
			structName:    "Example",
			expectedStart: 2,
			expectedEnd:   5,
		},
		{
			name:          "Sad Case: offset is not within a struct definition",
			src:           src,
			offset:        1,
			expectedError: true,
		},
		{
			name:          "Happy Case: offset is within a struct definition",
			src:           src,
			offset:        55,
			expectedStart: 2,
			expectedEnd:   5,
		},
		{
			name:          "Sad Case: empty line field",
			line:          "",
			expectedError: true,
		},
		{
			name:          "Sad Case: non-numeric start value",
			line:          "a,b",
			expectedError: true,
		},
		{
			name:          "Sad Case: more than 2 comma separated items provided",
			line:          "1,2,3,4,5",
			expectedError: true,
		},
		{
			name:          "Sad Case: non-numeric end value",
			line:          "2,c",
			expectedError: true,
		},
		{
			name:          "Sad Case: invalid range where start is bigger than end",
			line:          "2,1",
			expectedError: true,
		},
		{
			name:          "Happy Case: received start and end values",
			line:          "1,100",
			expectedStart: 1,
			expectedEnd:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			cfg := &config.Config{
				Fset:       token.NewFileSet(),
				Lines:      tt.line,
				Offset:     tt.offset,
				StructName: tt.structName,
			}

			cfg.File, err = parser.ParseFile(cfg.Fset, "example.go", src, parser.ParseComments)
			if err != nil {
				t.Fatalf("unexpected parser error: %v", err)
			}

			start, end, err := structfind.Find(cfg)
			if err != nil {
				if tt.expectedError {
					return
				}
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectedStart != start {
				t.Fatalf("expected start %d, got %d", tt.expectedStart, start)
			}
			if tt.expectedEnd != end {
				t.Fatalf("expected end %d, got %d", tt.expectedEnd, end)
			}
		})
	}
}
