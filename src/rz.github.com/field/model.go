package field

type (
	// Side is the side of a shape can be one of top, bottom, left or right.
	Side int

	// State describes a place on the field,
	// a place may be either have a match on it or be a space.
	State bool

	// MatchPosition describes a position of a present match, used during loading.
	MatchPosition struct {
		X int
		Y int
		S Side
	}

	// Copyable represents an object that can be copied.
	Copyable interface {
		Copy(bool) Copyable
	}
)

//noinspection GoExportedElementShouldHaveComment
const (
	Match State = true
	Space State = false

	Top Side = iota
	Bot
	Lft
	Rgt
)
