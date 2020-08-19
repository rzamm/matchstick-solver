package run

import "SolveMatches/src/rz.github.com/field"

// GameType represents the game type.
type GameType int

const (
	// RemoveGame only remove matches.
	RemoveGame GameType = iota
	// MoveGame remove then place matches.
	MoveGame
)

// A Level describes an initial state, a game type, the number of removable/movable matches
// and the number of shapes required.
type Level struct {
	Field          FieldI
	GameType       GameType
	Movable        int
	ShapesRequired int
}

// LvlTestMultipleSolutions testing level that runs quickly and has multiple solutions.
//goland:noinspection GoUnusedExportedFunction
func LvlTestMultipleSolutions(bit bool) *Level {
	var matches []*field.MatchPosition
	matches = append(matches, placeSquare(0, 0)...)

	return returnLevel(bit, MoveGame, 4, 1, 4, 4, matches)
}

// Lvl6 represents level 6.
//goland:noinspection GoUnusedExportedFunction
func Lvl6(bit bool) *Level {
	var matches []*field.MatchPosition
	matches = append(matches, placeSquare(1, 1)...)
	matches = append(matches, placeSquare(2, 1)...)
	matches = append(matches, placeSquare(0, 2)...)
	matches = append(matches, placeSquare(1, 2)...)
	matches = append(matches, placeSquare(2, 2)...)
	matches = append(matches, placeSquare(3, 2)...)

	return returnLevel(bit, RemoveGame, 6, 3, 4, 5, matches)
}

// Lvl16 represents level 16.
//goland:noinspection GoUnusedExportedFunction
func Lvl16(bit bool) *Level {
	matches := []*field.MatchPosition{
		{0, 0, field.Rgt},
		{0, 0, field.Bot},
		{0, 3, field.Top},
		{0, 3, field.Rgt},
		{3, 3, field.Top},
		{3, 3, field.Lft},
		{3, 0, field.Lft},
		{3, 0, field.Bot},

		{1, 0, field.Rgt},
		{1, 3, field.Rgt},
		{0, 1, field.Bot},
		{3, 1, field.Bot},
	}
	matches = append(matches, placeSquare(1, 1)...)
	matches = append(matches, placeSquare(1, 2)...)
	matches = append(matches, placeSquare(2, 1)...)
	matches = append(matches, placeSquare(2, 2)...)

	return returnLevel(bit, MoveGame, 8, 5, 4, 5, matches)
}

// Lvl16Test a faster version of level 16 due to less movable matches.
//goland:noinspection GoUnusedExportedFunction
func Lvl16Test(bit bool) *Level {
	matches := []*field.MatchPosition{
		{0, 0, field.Rgt},
		{0, 0, field.Bot},
		{0, 3, field.Top},
		{0, 3, field.Rgt},
		{3, 3, field.Top},
		{3, 3, field.Lft},
		{3, 0, field.Lft},
		{3, 0, field.Bot},

		{3, 0, field.Rgt},
		{3, 3, field.Rgt},
		{0, 3, field.Bot},
		{3, 3, field.Bot},
	}

	matches = append(matches, placeSquare(1, 1)...)
	matches = append(matches, placeSquare(1, 2)...)
	matches = append(matches, placeSquare(2, 1)...)
	matches = append(matches, placeSquare(2, 2)...)

	return returnLevel(bit, MoveGame, 4, 5, 4, 5, matches)
}

// Lvl19 represents level 19.
//goland:noinspection GoUnusedExportedFunction
func Lvl19(bit bool) *Level {
	matches := []*field.MatchPosition{
		{0, 0, field.Bot},
		{0, 1, field.Bot},
		{0, 2, field.Bot},
		{0, 3, field.Bot},

		{3, 0, field.Bot},
		{3, 1, field.Bot},
		{3, 2, field.Bot},
		{3, 3, field.Bot},
	}

	matches = append(matches, placeSquare(1, 1)...)
	matches = append(matches, placeSquare(1, 2)...)
	matches = append(matches, placeSquare(1, 3)...)
	matches = append(matches, placeSquare(2, 1)...)
	matches = append(matches, placeSquare(2, 2)...)
	matches = append(matches, placeSquare(2, 3)...)

	return returnLevel(bit, MoveGame, 6, 6, 4, 5, matches)
}

func returnLevel(bit bool, gameType GameType, movable, shapesRequired, width, height int,
	matches []*field.MatchPosition) *Level {

	var removableMatches int
	switch gameType {
	case RemoveGame:
		removableMatches = movable
	case MoveGame:
		removableMatches = 0
	default:
		panic("Unknown Game Type")
	}

	var f FieldI
	if bit {
		f = field.NewBitField(width, height, matches)
	} else {
		f = field.NewField(width, height, removableMatches, matches)
	}

	return &Level{
		Field:          f,
		GameType:       gameType,
		Movable:        movable,
		ShapesRequired: shapesRequired,
	}
}

func placeSquare(x, y int) []*field.MatchPosition {
	return []*field.MatchPosition{
		{x, y, field.Top},
		{x, y, field.Bot},
		{x, y, field.Lft},
		{x, y, field.Rgt},
	}
}
