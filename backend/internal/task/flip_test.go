package task

import "testing"

func TestFlipDirectionValid(t *testing.T) {
	tests := []struct {
		name      string
		direction FlipDirection
		want      bool
	}{
		{
			name:      "horizontal",
			direction: FlipDirectionHorizontal,
			want:      true,
		},
		{
			name:      "vertical",
			direction: FlipDirectionVertical,
			want:      true,
		},
		{
			name:      "empty",
			direction: "",
			want:      false,
		},
		{
			name:      "unknown",
			direction: "diagonal",
			want:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.direction.Valid(); got != test.want {
				t.Fatalf(
					"FlipDirection(%q).Valid() = %t, want %t",
					test.direction,
					got,
					test.want,
				)
			}
		},
		)
	}
}
