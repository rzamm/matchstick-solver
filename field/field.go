// Package field is a representation of the game board.
package field

type (
	// MatchSpace represents a match space on the field, can be Match (placed) or Space (empty space).
	MatchSpace struct {
		S State
	}
	// Cell represents cell in a field.
	Cell struct {
		Top, Bot, Lft, Rgt *State
	}
	// Field represents a match field.
	Field struct {
		matches   int
		spaces    int
		width     int
		height    int
		gridSpace [][]*Cell
		lineSpace []*State
		matchList []*State
		spaceList []*State

		visitedMatches  map[*State]interface{}
		squares         [][]*State // list of combinations of matches that may form a square
		requiredVisited int        // the required number of matches visited
	}
)

// NewField returns a new Field with a width, height and an initial placement of matches.
func NewField(width, height, removableMatches int, initialMatches []*MatchPosition) *Field {
	area := 2*width*height + width + height
	gridSpace := make([][]*Cell, width)
	lineSpace := make([]*State, area)
	squares := make([][]*State, 0)
	createLinkedSpaces(width, height, gridSpace, lineSpace, &squares)

	// place the initial matches on the field
	// update the match and space counts
	matches := 0
	spaces := area
	for _, m := range initialMatches {
		var match *State
		switch m.S {
		case Top:
			match = gridSpace[m.X][m.Y].Top
		case Bot:
			match = gridSpace[m.X][m.Y].Bot
		case Lft:
			match = gridSpace[m.X][m.Y].Lft
		case Rgt:
			match = gridSpace[m.X][m.Y].Rgt
		default:
			panic("unknown side")
		}
		lastState := match
		if *lastState == Space {
			*match = Match
			matches++
			spaces--
		}
	}

	// add matches and spaces to lists, used for trying combinations of removals and placements
	matchList := make([]*State, matches)
	spaceList := make([]*State, spaces)
	for i, mi, si := 0, 0, 0; i < area; i++ {
		m := lineSpace[i]
		if *m == Match {
			matchList[mi] = m
			mi++
		} else {
			spaceList[si] = m
			si++
		}
	}

	requiredVisited := matches - removableMatches

	return &Field{
		matches:   matches,
		spaces:    spaces,
		width:     width,
		height:    height,
		gridSpace: gridSpace,
		lineSpace: lineSpace,
		matchList: matchList,
		spaceList: spaceList,

		visitedMatches:  make(map[*State]interface{}, matches),
		squares:         squares,
		requiredVisited: requiredVisited,
	}
}

// GetWidth returns the width.
func (f *Field) GetWidth() int {
	return f.width
}

// GetHeight returns the height.
func (f *Field) GetHeight() int {
	return f.height
}

// CheckMatch returns the State of a match that is on the given Side of a Cell.
// Ex. CheckMatch(2, 3, Top) returns the State of the match on the Top Side of the Cell at (2, 3).
func (f *Field) CheckMatch(x, y int, side Side) State {
	cell := f.gridSpace[x][y]
	switch side {
	case Top:
		return *cell.Top
	case Bot:
		return *cell.Bot
	case Lft:
		return *cell.Lft
	case Rgt:
		return *cell.Rgt
	default:
		panic("unknown side")
	}
}

// GetMatchesCount returns the number of initial matches.
func (f *Field) GetMatchesCount() int {
	return f.matches
}

// GetSpacesCount returns the number of initial spaces.
func (f *Field) GetSpacesCount() int {
	return f.spaces
}

// ChangeToState will change all matches or spaces from a list of indices to the desired state.
// Ex. ChangeToState([]int{1, 2, 3}, Match, Space)
// finds matches 1, 2, 3 in the match list, and changes them to spaces.
func (f *Field) ChangeToState(l []int, fromState State, toState State) {
	var mapping []*State
	if fromState == Match {
		mapping = f.matchList
	} else {
		mapping = f.spaceList
	}

	for _, v := range l {
		*mapping[v] = toState
	}
}

// CheckSquares returns true if the number of squares is equal to the amount required
// and all matches were visited.
func (f *Field) CheckSquares(requiredShapes int) bool {
	squareCount := 0
	addUnique := func(newMatches ...*State) {
		for _, nm := range newMatches {
			_, visited := f.visitedMatches[nm]
			if !visited {
				f.visitedMatches[nm] = nil
			}
		}
	}
	for _, square := range f.squares {
		present := true // all matches in the square are assumed present
		for _, match := range square {
			if *match == Space {
				present = false
				break
			}
		}
		if present {
			squareCount++
			addUnique(square...)
		}
	}
	visited := len(f.visitedMatches)

	// clear visited set (compiler optimized)
	for k := range f.visitedMatches {
		delete(f.visitedMatches, k)
	}

	return squareCount == requiredShapes && visited == f.requiredVisited
}

// Copy returns a copy of this Field.
// If displayOnly is set, then this copy can only be used to display a state, and does not require a spaceList.
// If displayOnly is not set, the Field's spaceList will also be copied and
// this copy can be used for generating more possible field states.
// todo: don't copy squares if it's for display only
func (f *Field) Copy(displayOnly bool) Copyable {
	w := f.width
	h := f.height

	gridSpace := make([][]*Cell, w)
	lineSpace := make([]*State, len(f.lineSpace))
	squares := make([][]*State, 0)
	createLinkedSpaces(w, h, gridSpace, lineSpace, &squares)
	// set to current match layout
	for i, m := range f.lineSpace {
		*lineSpace[i] = *m
	}

	newField := &Field{
		matches:         f.matches,
		spaces:          f.spaces,
		width:           w,
		height:          h,
		gridSpace:       gridSpace,
		lineSpace:       lineSpace,
		matchList:       nil, // not needed
		spaceList:       nil, // may be included
		visitedMatches:  make(map[*State]interface{}, f.matches),
		squares:         squares,
		requiredVisited: f.requiredVisited,
	}

	if !displayOnly {
		spaceList := make([]*State, len(f.spaceList))
		spaceListIndex := 0
		for i, m := range f.lineSpace {
			if m == f.spaceList[spaceListIndex] {
				spaceList[spaceListIndex] = lineSpace[i]
				spaceListIndex++
			}
		}

		newField.spaceList = spaceList
	}

	return newField
}

func createLinkedSpaces(width, height int, gridSpace [][]*Cell, lineSpace []*State, squares *[][]*State) {
	lineSpaceIndex := 0

	addToLine := func(m *State) {
		lineSpace[lineSpaceIndex] = m
		lineSpaceIndex++
	}

	// first add only the tops and lefts
	for i := 0; i < width; i++ {
		gridSpace[i] = make([]*Cell, height)
		for j := 0; j < height; j++ {
			c := &Cell{}
			var mTop State
			var mLeft State
			c.Top = &mTop
			c.Lft = &mLeft

			gridSpace[i][j] = c
			addToLine(&mTop)
			addToLine(&mLeft)

			// link previous cell's right to our left
			if i > 0 {
				gridSpace[i-1][j].Rgt = &mLeft
			}
			// link previous row's bottom to our top
			if j > 0 {
				gridSpace[i][j-1].Bot = &mTop
			}
		}
	}

	// then add the last row of bottoms
	for i := 0; i < width; i++ {
		c := gridSpace[i][height-1]
		var mBottom State
		c.Bot = &mBottom
		addToLine(&mBottom)
	}
	// and the last column of rights
	for j := 0; j < height; j++ {
		c := gridSpace[width-1][j]
		var mRight State
		c.Rgt = &mRight
		addToLine(&mRight)
	}

	// init squares
	// this is a list of a set of states
	// each set represents a set matches that form a square
	// each set length is a multiple of four
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			for size := 1; size <= width-i && size <= height-j; size++ {
				square := make([]*State, 0, 4*size)
				for k := 0; k < size; k++ {
					mTop := gridSpace[i+k][j].Top
					mBot := gridSpace[i+k][j+size-1].Bot
					mLft := gridSpace[i][j+k].Lft
					mRgt := gridSpace[i+size-1][j+k].Rgt
					square = append(square, mTop, mBot, mLft, mRgt)
				}
				*squares = append(*squares, square)
			}
		}
	}
}
