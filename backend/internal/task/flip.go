package task

type FlipDirection string

const (
	FlipDirectionHorizontal FlipDirection = "horizontal"
	FlipDirectionVertical   FlipDirection = "vertical"
)

func (d FlipDirection) Valid() bool {
	switch d {
	case FlipDirectionHorizontal, FlipDirectionVertical:
		return true

	default:
		return false
	}
}
