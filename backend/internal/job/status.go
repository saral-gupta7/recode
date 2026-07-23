package job

type Status string

const (
	StatusQueued     Status = "queued"
	StatusProcessing Status = "processing"
	StatusCancelling Status = "cancelling"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
	StatusCancelled  Status = "cancelled"
	StatusExpired    Status = "expired"
)

func (s Status) Valid() bool {
	switch s {
	case StatusQueued,
		StatusProcessing,
		StatusCancelling,
		StatusCompleted,
		StatusFailed,
		StatusCancelled,
		StatusExpired:
		return true

	default:
		return false
	}
}

func (s Status) CanTransitionTo(next Status) bool {
	if !s.Valid() || !next.Valid() {
		return false
	}

	switch s {
	case StatusQueued:
		return next == StatusProcessing ||
			next == StatusCancelled

	case StatusProcessing:
		return next == StatusQueued ||
			next == StatusCancelling ||
			next == StatusCompleted ||
			next == StatusFailed

	case StatusCancelling:
		return next == StatusCancelled

	case StatusCompleted,
		StatusFailed,
		StatusCancelled:
		return next == StatusExpired

	case StatusExpired:
		return false

	default:
		return false
	}
}
