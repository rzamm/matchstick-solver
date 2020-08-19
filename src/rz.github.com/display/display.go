// Package display contains functions for rendering fields to the log file.
package display

import (
	"SolveMatches/src/rz.github.com/field"
	"SolveMatches/src/rz.github.com/logg"
)

const bot = '_'
const vrt = '|'

// FieldI represents a drawable field.
type FieldI interface {
	GetWidth() int
	GetHeight() int
	CheckMatch(x, y int, side field.Side) field.State
}

// Draw draws a field to the log file, does not flush.
func Draw(f FieldI) {
	w := f.GetWidth()
	h := f.GetHeight()

	getChar := func(x, y int, side field.Side, r rune) rune {
		if f.CheckMatch(x, y, side) == field.Match {
			return r
		}
		return ' '
	}

	for i := 0; i < w; i++ {
		logg.Printf("--")
	}
	logg.Println("-")

	for i := 0; i < w; i++ {
		logg.Printf(" %c", getChar(i, 0, field.Top, bot))
	}
	logg.Println(" ")

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			logg.Printf("%c%c", getChar(i, j, field.Lft, vrt), getChar(i, j, field.Bot, bot))
		}
		logg.Printf("%c\n", getChar(w-1, j, field.Rgt, vrt))
	}

	for i := 0; i < w; i++ {
		logg.Printf("--")
	}
	logg.Println("-")
}
