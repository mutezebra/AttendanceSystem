package weightedrand

type Item struct {
	Key    int64
	Weight int
}

// WeightedRandom 权重越高,越不容易被抽中,return uids
func WeightedRandom(items []*Item, count int, action, number int8) ([]int64, error) {
	switch action {
	case random:
		return randomCall(items, count)
	case weight:
		return weightedCall(items, count)
	case week:
		return weekCall(items, count)
	case lucky:
		return luckyCall(items, count, number)
	case timestamp:
		return timestampCall(items, count)
	default:
		return weightedCall(items, count)
	}
}
