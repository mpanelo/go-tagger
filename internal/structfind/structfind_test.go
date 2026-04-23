package structfind_test

import (
	"errors"
	"testing"

	"github.com/mpanelo/custom-go-tool/internal/structfind"
)

func TestLineSelection(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		expectedStart int
		expectedEnd   int
		expectedError error
	}{
		{
			name:          "Sad Case: Empty Line field",
			line:          "",
			expectedError: structfind.ErrInvalidLineSelection,
		},
		{
			name:          "Sad Case: Non-numeric start value",
			line:          "a,b",
			expectedError: structfind.ErrInvalidLineSelection,
		},
		{
			name:          "Sad Case: More than 2 comma separated items provided",
			line:          "1,2,3,4,5",
			expectedError: structfind.ErrInvalidLineSelection,
		},
		{
			name:          "Sad Case: Non-numeric end value",
			line:          "2,c",
			expectedError: structfind.ErrInvalidLineSelection,
		},
		{
			name:          "Sad Case: Invalid range where start is bigger than end",
			line:          "2,1",
			expectedError: structfind.ErrInvalidLineSelection,
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
				if tt.expectedError == nil {
					t.Fatalf("unexpected error %v", err)
				}

				if !errors.Is(err, tt.expectedError) {
					t.Fatalf("expected error %v, got %v", tt.expectedError, err)
				}

				// This means the err and tt.expectedError matched
				return
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
}

func TestStructSelection(t *testing.T) {
}
