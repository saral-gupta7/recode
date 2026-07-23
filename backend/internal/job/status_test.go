package job

import "testing"

func TestStatusValid(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{name: "queued", status: StatusQueued, want: true},
		{name: "processing", status: StatusProcessing, want: true},
		{name: "cancelling", status: StatusCancelling, want: true},
		{name: "completed", status: StatusCompleted, want: true},
		{name: "failed", status: StatusFailed, want: true},
		{name: "cancelled", status: StatusCancelled, want: true},
		{name: "expired", status: StatusExpired, want: true},
		{name: "empty", status: Status(""), want: false},
		{name: "unknown", status: Status("almost_finished"), want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.status.Valid(); got != test.want {
				t.Fatalf("Status(%q).Valid() = %t, want %t", test.status, got, test.want)
			}
		})
	}
}

func TestStatusCanTransitionTo(t *testing.T) {
	statuses := []Status{
		StatusQueued,
		StatusProcessing,
		StatusCancelling,
		StatusCompleted,
		StatusFailed,
		StatusCancelled,
		StatusExpired,
	}

	allowed := map[Status]map[Status]bool{
		StatusQueued: {
			StatusProcessing: true,
			StatusCancelled:  true,
		},
		StatusProcessing: {
			StatusQueued:     true,
			StatusCancelling: true,
			StatusCompleted:  true,
			StatusFailed:     true,
		},
		StatusCancelling: {
			StatusCancelled: true,
		},
		StatusCompleted: {
			StatusExpired: true,
		},
		StatusFailed: {
			StatusExpired: true,
		},
		StatusCancelled: {
			StatusExpired: true,
		},
		StatusExpired: {},
	}

	for _, from := range statuses {
		for _, to := range statuses {
			want := allowed[from][to]

			t.Run(string(from)+"_to_"+string(to), func(t *testing.T) {
				if got := from.CanTransitionTo(to); got != want {
					t.Fatalf("%q.CanTransitionTo(%q) = %t, want %t", from, to, got, want)
				}
			})
		}
	}
}

func TestStatusCanTransitionToRejectsInvalidStatuses(t *testing.T) {
	invalid := Status("unknown")

	tests := []struct {
		name string
		from Status
		to   Status
	}{
		{name: "invalid source", from: invalid, to: StatusQueued},
		{name: "invalid target", from: StatusQueued, to: invalid},
		{name: "invalid source and target", from: invalid, to: invalid},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.from.CanTransitionTo(test.to) {
				t.Fatalf("%q.CanTransitionTo(%q) = true, want false", test.from, test.to)
			}
		})
	}
}
