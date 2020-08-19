package field

import (
	"fmt"
)

type (
	// BitField represents a match field.
	BitField struct {
		matches       int
		spaces        int
		width         int
		height        int
		linearMapping map[int]uint64
		matchSpace    *uint64
		matchList     []uint64
		spaceList     []uint64
		squares       []uint64 // list of combinations of matches that can form a square
	}
)

// NewBitField returns a new BitField with a width, height and an initial placement of matches.
func NewBitField(width, height int, initialMatches []*MatchPosition) *BitField {
	area := 2*width*height + width + height
	if area > 64 {
		panic(fmt.Sprintf("cannot fit field with %d spaces into int64", area))
	}

	linearMapping := make(map[int]uint64)

	to1D := func(x, y int, s Side) int {
		return int(s) + width*(y+height*x)
	}

	matchSpace := uint64(0)
	matchBit := uint64(1)
	// first add only the tops and lefts
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			mTop := matchBit
			matchBit = matchBit << 1
			mLeft := matchBit
			matchBit = matchBit << 1

			linearMapping[to1D(i, j, Top)] = mTop
			linearMapping[to1D(i, j, Lft)] = mLeft

			// previous cell's right has the same value as the current left
			if i > 0 {
				linearMapping[to1D(i-1, j, Rgt)] = mLeft
			}
			// previous cell's bottom has the same value as the current top
			if j > 0 {
				linearMapping[to1D(i, j-1, Bot)] = mTop
			}
		}
	}

	// then add the last row of bottoms
	for i := 0; i < width; i++ {
		mBottom := matchBit
		matchBit = matchBit << 1
		linearMapping[to1D(i, height-1, Bot)] = mBottom
	}
	// and the last column of rights
	for j := 0; j < height; j++ {
		mRight := matchBit
		matchBit = matchBit << 1
		linearMapping[to1D(width-1, j, Rgt)] = mRight
	}

	// init squares
	// this is a list of a set of states
	// each set represents a set matches that form a square
	// each set length is a multiple of four
	squares := make([]uint64, 0)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			for size := 1; size <= width-i && size <= height-j; size++ {
				square := uint64(0)
				for k := 0; k < size; k++ {
					mTop := linearMapping[to1D(i+k, j, Top)]
					mBot := linearMapping[to1D(i+k, j+size-1, Bot)]
					mLft := linearMapping[to1D(i, j+k, Lft)]
					mRgt := linearMapping[to1D(i+size-1, j+k, Rgt)]
					square |= mTop | mBot | mLft | mRgt
				}
				squares = append(squares, square)
			}
		}
	}

	// place the initial matches on the field
	// update the match and space counts
	matches := 0
	spaces := area
	for _, m := range initialMatches {
		match := linearMapping[to1D(m.X, m.Y, m.S)]
		lastState := matchSpace & match
		if lastState == 0 {
			matchSpace |= match
			matches++
			spaces--
		}
	}

	// add matches and spaces to lists, used for trying combinations of removals and placements
	matchList := make([]uint64, matches)
	spaceList := make([]uint64, spaces)
	for i, mi, si := 0, 0, 0; i < area; i++ {
		m := uint64(1 << i)
		if matchSpace&m > 0 {
			matchList[mi] = m
			mi++
		} else {
			spaceList[si] = m
			si++
		}
	}

	return &BitField{
		matches:       matches,
		spaces:        spaces,
		width:         width,
		height:        height,
		linearMapping: linearMapping,
		matchList:     matchList,
		spaceList:     spaceList,
		matchSpace:    &matchSpace,

		squares: squares,
	}
}

func (f *BitField) to1D(x, y int, s Side) int {
	return int(s) + f.width*(y+f.height*x)
}

func (f *BitField) getMatchBit(x, y int, s Side) uint64 {
	bit, ok := f.linearMapping[f.to1D(x, y, s)]
	if !ok {
		panic("out of bounds")
	}
	return bit
}

// GetWidth returns the width.
func (f *BitField) GetWidth() int {
	return f.width
}

// GetHeight returns the height.
func (f *BitField) GetHeight() int {
	return f.height
}

// CheckMatch returns the State of a match that is on the given Side of a Cell.
// Ex. CheckMatch(2, 3, Top) returns the State of the match on the Top Side of the Cell at (2, 3).
func (f *BitField) CheckMatch(x, y int, s Side) State {
	matchBit := f.getMatchBit(x, y, s)
	if matchBit&*f.matchSpace > 0 {
		return Match
	}
	return Space
}

// GetMatchesCount returns the number of initial matches.
func (f *BitField) GetMatchesCount() int {
	return f.matches
}

// GetSpacesCount returns the number of initial spaces.
func (f *BitField) GetSpacesCount() int {
	return f.spaces
}

// ChangeToState will change all matches or spaces from a list of indices to the desired state.
// Ex. ChangeToState([]int{1, 2, 3}, Match, Space)
// finds matches 1, 2, 3 in the match list, and changes them to spaces.
func (f *BitField) ChangeToState(l []int, fromState State, toState State) {
	var list []uint64
	if fromState == Match {
		list = f.matchList
	} else {
		list = f.spaceList
	}

	var changeBit func(b uint64)
	if toState == Match {
		// set bit
		changeBit = func(b uint64) {
			*f.matchSpace |= b
		}
	} else {
		// clear bit
		changeBit = func(b uint64) {
			*f.matchSpace &^= b
		}
	}

	for _, v := range l {
		changeBit(list[v])
	}
}

// CheckSquares returns true if the number of squares is equal to the amount required
// and all matches were visited.
func (f *BitField) CheckSquares(requiredShapes int) bool {
	count := 0
	visitedMatches := uint64(0)
	for _, s := range f.squares {
		if *f.matchSpace&s == s {
			count++
			visitedMatches |= s
		}
	}

	return count == requiredShapes && *f.matchSpace^visitedMatches == 0
}

// Copy returns a copy of this BitField.
func (f *BitField) Copy(bool) Copyable {
	newMatchSpace := *f.matchSpace
	return &BitField{
		matches:       f.matches,
		spaces:        f.spaces,
		width:         f.width,
		height:        f.height,
		linearMapping: f.linearMapping,
		matchSpace:    &newMatchSpace,
		matchList:     f.matchList,
		spaceList:     f.spaceList,
		squares:       f.squares,
	}
}
