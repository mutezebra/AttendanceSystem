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

func WeightedRandom(items []*Item, count int) ([]int64, error) {
	totalWeight := 0
	for _, item := range items {
		totalWeight += item.Weight
	}
	if totalWeight == 0 {
		return nil, errors.Wrap(fmt.Errorf("total weight is 0"), fmt.Sprintf("items: %v", items))
	}

	originSize := len(items)
	for i := originSize; i > count; i-- {
		r := rand.Intn(totalWeight)
		for index, item := range items {
			if r < item.Weight {
				items = removeKey(items, index)
				totalWeight = subTotalWeight(totalWeight, item.Weight)
				break
			}
			r -= item.Weight
		}
	}

	result := make([]int64, count)
	for i := 0; i < count; i++ {
		result[i] = items[i].Key
	}

	return result, nil
}

func subTotalWeight(origin, num int) int {
	return origin - num
}

func removeKey(items []*Item, index int) []*Item {
	if index < 0 || index >= len(items) {
		// 索引越界，直接返回原切片
		return items
	}
	return append(items[:index], items[index+1:]...)
}
