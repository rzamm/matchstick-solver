package run

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rzamm/matchstick-solver/display"
	"github.com/rzamm/matchstick-solver/field"
	"github.com/rzamm/matchstick-solver/logg"
)

// testing level that runs quickly and has multiple solutions
func multipleSolutionsLevel(bit bool) *Level {
	var matches []*field.MatchPosition
	matches = append(matches, placeSquare(0, 0)...)

	return returnLevel(bit, moveGame, 4, 1, 4, 4, matches)
}

func Test_LvlTestMultiSol(t *testing.T) {
	doRun(t, multipleSolutionsLevel, true, false)
}

func Test_Lvl6(t *testing.T) {
	doRun(t, Lvl6, false, false)
}

func Test_Lvl6_Bit(t *testing.T) {
	doRun(t, Lvl6, true, false)
}

func Test_Lvl16Test(t *testing.T) {
	doRun(t, Lvl16Test, false, false)
}

func Test_Lvl16Test_Bit(t *testing.T) {
	doRun(t, Lvl16Test, true, false)
}

func Test_Lvl16Test_Bit_One(t *testing.T) {
	doRun(t, Lvl16, true, true)
}

func doRun(t *testing.T, newLevel func(bool) *Level, bitwise, oneSolution bool) {
	logg.Println("Starting Layout:")
	lvl := newLevel(bitwise)
	display.Draw(lvl.Field)
	logg.Flush()
	runner := NewRun(lvl)
	runner.PrintStats()

	fs := runner.SolveGame(oneSolution)
	assert.NotEmpty(t, fs)
	logg.Println("\n\nSolution:")
	for _, f := range fs {
		display.Draw(f)
	}
	logg.Flush()
}

func BenchmarkMoveGame(b *testing.B) {
	lvl := Lvl16Test(true)
	runner := NewRun(lvl)
	runner.PrintStats()
	for i := 0; i < b.N; i++ {
		runner.MoveGame(true)
		fmt.Println(i)
	}
}
