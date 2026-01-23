package mathz

// 1个基点等于0.01%  0.0001
func BasicPoints[T Float](f T) int64 {
	return int64(f * 10000)
}
