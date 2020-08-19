// Package run provides functions for creating and solving levels.
package run

import (
	"fmt"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gonum.org/v1/gonum/stat/combin"

	"SolveMatches/src/rz.github.com/display"
	"SolveMatches/src/rz.github.com/ec"
	"SolveMatches/src/rz.github.com/field"
)

type (
	// FieldI represents a field.
	FieldI interface {
		display.FieldI
		GetSpacesCount() int
		GetMatchesCount() int
		ChangeToState(list []int, fromState field.State, toState field.State)
		CheckSquares(requiredShapes int) bool
		Copy(bool) field.Copyable
	}
	// Run is a collection of information needed to find a solution to a level
	// as well as metadata.
	Run struct {
		field             FieldI
		matchCount        int
		spaceCount        int
		movable           int
		removeCombsTotal  int
		placeCombsTotal   int
		totalCombinations int
		shapesRequired    int
		gameType          GameType
		printer           *message.Printer
	}
)

// NewRun creates a new Run from a Level.
func NewRun(lvl *Level) *Run {
	mCount := lvl.Field.GetMatchesCount()
	sCount := lvl.Field.GetSpacesCount()
	movable := lvl.Movable
	removeCombs := combin.Binomial(mCount, movable)
	var placeCombs int
	var totalCombinations int
	switch lvl.GameType {
	case RemoveGame:
		placeCombs = 0
		totalCombinations = removeCombs
	case MoveGame:
		placeCombs = combin.Binomial(sCount, movable)
		totalCombinations = removeCombs * placeCombs
	}

	return &Run{
		field:             lvl.Field,
		matchCount:        mCount,
		spaceCount:        sCount,
		movable:           movable,
		removeCombsTotal:  removeCombs,
		placeCombsTotal:   placeCombs,
		shapesRequired:    lvl.ShapesRequired,
		gameType:          lvl.GameType,
		totalCombinations: totalCombinations,
		printer:           message.NewPrinter(language.English),
	}
}

// PrintStats prints statistics to the log file.
func (r *Run) PrintStats() {
	fmt.Println("matches", r.matchCount)
	if r.gameType == MoveGame {
		fmt.Println("spaces", r.spaceCount)
	}
	fmt.Println("movable", r.movable)
	r.printer.Printf("remove combs %d\n", r.removeCombsTotal)
	if r.gameType == MoveGame {
		r.printer.Printf("place combs %d\n", r.placeCombsTotal)
	}
	r.printer.Printf("total %d\n\n", r.totalCombinations)
}

// SolveGame runs the Run and returns solutions.
// It returns a slice of fields in the solved state (empty slice if no solutions).
// If oneSolution is set, SolveGame will return only the first solution that it finds.
func (r *Run) SolveGame(oneSolution bool) []FieldI {
	switch r.gameType {
	case RemoveGame:
		return r.RemoveGame(oneSolution)
	case MoveGame:
		return r.MoveGame(oneSolution)
	default:
		panic("Unknown Game Type")
	}
}

// RemoveGame runs the Run as the remove game type and returns solutions.
// It returns a slice of fields in the solved state (empty slice if no solutions).
// If oneSolution is set, SolveGame will return only the first solution that it finds.
func (r *Run) RemoveGame(oneSolution bool) []FieldI {
	solutions := make([]FieldI, 0)
	removable := r.movable
	removeComb := make([]int, removable)

	// init removeComb to [0, 1 , 2 ... r.movable-1]
	for i := 0; i < r.movable; i++ {
		removeComb[i] = i
	}
	removeCombIndex := combin.CombinationIndex(removeComb, r.matchCount, r.movable)

	for removeCombIndex < r.removeCombsTotal {
		// remove the matches that we guess we need to remove
		r.field.ChangeToState(removeComb, field.Match, field.Space)

		// check if solving combination found
		if r.field.CheckSquares(r.shapesRequired) {
			solution := r.field.Copy(true).(FieldI)
			solutions = append(solutions, solution)
			if oneSolution {
				return solutions
			}
		}

		// put the matches we removed back
		r.field.ChangeToState(removeComb, field.Match, field.Match)

		ec.NextCombination(removeComb, r.matchCount, r.movable)
		removeCombIndex++
	}

	return solutions
}

// MoveGame runs the Run as the move game type and returns solutions.
// It returns a slice of fields in the solved state (empty slice if no solutions).
// If oneSolution is set, SolveGame will return only the first solution that it finds.
func (r *Run) MoveGame(oneSolution bool) []FieldI {
	removeComb := make([]int, r.movable)
	// init removeComb to [0, 1 , 2 ... r.movable-1]
	for i := 0; i < r.movable; i++ {
		removeComb[i] = i
	}
	removeCombIndex := combin.CombinationIndex(removeComb, r.matchCount, r.movable)

	found := make(chan *taskReturn)

	// this task runs through the place combinations and sends any solutions it finds
	task := func(tp *taskParams) {
		placeComb := make([]int, r.movable)
		// init placeComb to [0, 1 , 2 ... r.movable-1]
		for i := 0; i < r.movable; i++ {
			placeComb[i] = i
		}
		placeCombIndex := 0
		for placeCombIndex < r.placeCombsTotal {
			// place the matches where we guess they should go
			tp.f.ChangeToState(placeComb, field.Space, field.Match)

			if tp.f.CheckSquares(r.shapesRequired) {
				// solving combinations found, send solution
				found <- &taskReturn{
					f: tp.f.Copy(true).(FieldI),
				}
			}

			// remove the matches we placed
			tp.f.ChangeToState(placeComb, field.Space, field.Space)
			ec.NextCombination(placeComb, r.spaceCount, r.movable)
			placeCombIndex++
		}
	}
	workers := Workers(task, found)

	go func() {
		startTime := time.Now()
		var checks int64 = 0
		for removeCombIndex < r.removeCombsTotal {
			// remove the matches that we guess we need to remove
			r.field.ChangeToState(removeComb, field.Match, field.Space)

			params := taskParams{
				f:               r.field.Copy(false).(FieldI),
				removeCombIndex: removeCombIndex,
			}
			workers <- &params

			if removeCombIndex%50 == 0 {
				checks++
				average := time.Duration(int64(time.Now().Sub(startTime)) / checks)
				r.printer.Printf("%d out of %d average: %v\n",
					removeCombIndex, r.removeCombsTotal, average)
			}
			// put the matches we removed back
			r.field.ChangeToState(removeComb, field.Match, field.Match)
			ec.NextCombination(removeComb, r.matchCount, r.movable)
			removeCombIndex++
		}
		close(workers)
	}()

	solutions := make([]FieldI, 0)
	for result := range found {
		solutions = append(solutions, result.f)
		if oneSolution {
			break
		}
	}

	return solutions
}
