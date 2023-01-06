package bar_graph

import (
	"testing"
)

func TestDrawBar(t *testing.T) {
	DrawBarGraph([]int{1, 4, 2, 6, 8}, false, 2, 4)
}
