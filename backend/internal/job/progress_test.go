package job

import (
	"errors"
	"testing"
)

func TestNewProgress(t *testing.T) {
	tests := []struct {
		name  string
		value int
	}{
		{name: "minimum", value: 0},
		{name: "middle", value: 50},
		{name: "maximum", value: 100},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewProgress(test.value)
			if err != nil {
				t.Fatalf("NewProgress(%d) error = %v", test.value, err)
			}
			if got.Value() != test.value {
				t.Errorf("Value() = %d, want %d", got.Value(), test.value)
			}
		})
	}
}

func TestNewProgressRejectsOutOfRangeValue(t *testing.T) {
	values := []int{-1, 101}

	for _, value := range values {
		got, err := NewProgress(value)

		if got != 0 {
			t.Errorf("NewProgress(%d) = %d, want zero", value, got)
		}
		if !errors.Is(err, ErrInvalidProgress) {
			t.Fatalf(
				"NewProgress(%d) error = %v, want errors.Is(_, ErrInvalidProgress)",
				value,
				err,
			)
		}
	}
}
