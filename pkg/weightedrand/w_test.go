package weightedrand

import (
	"fmt"
	"testing"
)

func TestWeightedRandom(t *testing.T) {
	items := []Item{
		{1, 0},
		{2, 1},
		{3, 0},
		{4, 0},
	}
	times := make(map[int64]int)
	keys, err := WeightedRandom(items, 1)
	if err != nil {
		t.Fatal(err)
	}
	for i := range keys {
		times[keys[i]]++
	}
	fmt.Println(times)
}
