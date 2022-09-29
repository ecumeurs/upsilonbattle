package pattern

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position"
)

func TestPathTo(t *testing.T) {
	startPos := position.New(30, 20, 0)
	endPos := position.New(40, 10, 0)
	path := PathTo2D(endPos.Substract(startPos))

	// check that all path steps are adjascent to the previous one.
	for i := 1; i < len(path); i++ {
		if !path[i].IsAdjacent(path[i-1], 2) {
			t.Errorf("Path step %d is not adjacent to the previous one", i)
		}
	}

	t.Log("Start position:", startPos)
	expandedPath := path.Apply2D(startPos)
	for i := 0; i < len(expandedPath); i++ {
		t.Logf("%d: %v", i, expandedPath[i])
	}
	t.Log("End position:", endPos)

	// check that the last path step is the end position.
	if !expandedPath[len(expandedPath)-1].Equals(endPos) {
		t.Errorf("Last path step is not the end position")
	}
}
