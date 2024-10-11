package weightedrand

import (
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

const (
	random    = 1
	weight    = 2
	week      = 3
	lucky     = 4
	timestamp = 5
)

func randomCall(items []*Item, count int) ([]int64, error) {
	in := func(sli []int64, key int64) bool {
		for i := range sli {
			if sli[i] == key {
				return true
			}
		}
		return false
	}

	results := make([]int64, count)
	for i := 0; i < count; i++ {
		result := items[rand.Intn(len(items))].Key
		if !in(results, result) {
			results[i] = result
		} else {
			i--
		}
	}
	return results, nil
}

func weightedCall(items []*Item, count int) ([]int64, error) {
	if len(items) <= count {
		result := make([]int64, len(items))
		for i := 0; i < len(items); i++ {
			result[i] = items[i].Key
		}
		return result, nil
	}

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

func luckyCall(items []*Item, count int, num int8) ([]int64, error) {
	ready := make([]*Item, 0)
	for _, item := range items {
		if item.Key%10 == int64(num) {
			ready = append(ready, item)
		}
	}
	if len(ready) == 0 {
		ready = items
	}
	return weightedCall(ready, count)
}

func weekCall(items []*Item, count int) ([]int64, error) {
	num := time.Now().Weekday()
	if num == time.Sunday {
		num = 7
	}
	return luckyCall(items, count, int8(num))
}

func timestampCall(items []*Item, count int) ([]int64, error) {
	num := time.Now().Unix() % 10
	return luckyCall(items, count, int8(num))
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
