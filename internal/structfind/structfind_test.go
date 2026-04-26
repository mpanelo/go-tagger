package structfind_test

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/mpanelo/go-tagger/internal/structfind"
)

func TestLineSelection(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		expectedStart int
		expectedEnd   int
		expectedError bool
	}{
		{
			name:          "Sad Case: Empty Line field",
			line:          "",
			expectedError: true,
		},
		{
			name:          "Sad Case: Non-numeric start value",
			line:          "a,b",
			expectedError: true,
		},
		{
			name:          "Sad Case: More than 2 comma separated items provided",
			line:          "1,2,3,4,5",
			expectedError: true,
		},
		{
			name:          "Sad Case: Non-numeric end value",
			line:          "2,c",
			expectedError: true,
		},
		{
			name:          "Sad Case: Invalid range where start is bigger than end",
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
			sf := structfind.StructFinder{
				Line: tt.line,
			}
			start, end, err := sf.LineSelection()
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

func TestOffsetSelection(t *testing.T) {
	src := `package main
	type Example struct {
		A string
		B int /* target */
	}
	func main() {}`

	tests := []struct {
		name          string
		offset        int
		src           string
		expectedStart int
		expectedEnd   int
		expectedError bool
	}{
		{
			name:          "Sad Case: offset is not within a struct definition",
			src:           src,
			offset:        0,
			expectedError: true,
		},
		{
			name:          "Happy Case: offset is within a struct definition",
			src:           src,
			offset:        55,
			expectedStart: 2,
			expectedEnd:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()

			file, err := parser.ParseFile(fset, "example.go", src, parser.ParseComments)
			if err != nil {
				t.Fatalf("unexpected parser error: %v", err)
			}

			st := structfind.StructFinder{
				Offset:  tt.offset,
				FileSet: fset,
			}
			start, end, err := st.OffsetSelection(file)
			if err != nil {
				if tt.expectedError {
					return
				}
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

func TestStructSelection(t *testing.T) {
}
