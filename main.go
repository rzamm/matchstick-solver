package main

import (
	"SolveMatches/src/rz.github.com/display"
	"SolveMatches/src/rz.github.com/logg"
	"SolveMatches/src/rz.github.com/run"
)

func main() {
	lvl := run.Lvl19(true)

	runner := run.NewRun(lvl)
	runner.PrintStats()

	logg.Println("Starting Layout: ")
	display.Draw(lvl.Field)
	logg.Println("\n\n")

	fs := runner.SolveGame(true)
	if len(fs) == 0 {
		logg.Println("No Solutions Found")
	} else {
		logg.Println("Solution: ")
		for _, f := range fs {
			display.Draw(f)
		}
	}

	logg.Flush()
}
