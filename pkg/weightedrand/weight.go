package weightedrand

import (
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
)

type Item struct {
	Key    int64
	Weight int
}

func WeightedRandom(items []Item, count int) ([]int64, error) {
	totalWeight := 0
	for _, item := range items {
		totalWeight += item.Weight
	}
	if totalWeight == 0 {
		return nil, errors.Wrap(fmt.Errorf("total weight is 0"), fmt.Sprintf("items: %v", items))
	}

	result := make([]int64, 0, count)

	for i := 0; i < count; i++ {
		r := rand.Intn(totalWeight)
		for _, item := range items {
			if r < item.Weight {
				result = append(result, item.Key)
				break
			}
			r -= item.Weight
		}
	}

	return result, nil
}
