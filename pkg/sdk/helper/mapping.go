package helper

func Mapping[K any, V any](params []K, f func(x K) V) []V {

	result := make([]V, len(params))

	for i := range params {
		x := params[i]
		result[i] = f(x)
		//result = append(result, f(x))
	}

	return result
}

func MappingIndexed[K any, V any](params []K, f func(x K, i int) V) []V {

	result := make([]V, len(params))

	for i := range params {
		x := params[i]
		result[i] = f(x, i)
		//result = append(result, f(x))
	}

	return result
}

func Filter[T any](params []T, f func(param T) bool) []T {

	var result []T

	for i := range params {
		x := params[i]
		if f(x) {
			result = append(result, x)
		}
	}

	return result
}

func FilterAndMapping[T any, V any](params []T, f func(param T) *V) []V {

	var result []V

	for i := range params {
		x := params[i]

		val := f(x)
		if val != nil {
			result = append(result, *val)
		}
	}

	return result
}

//type Chain[T any] struct {
//	source []T
//}
//
//func NewChain[T any](source []T) Chain[T] {
//	return Chain[T]{
//		source: source,
//	}
//}

//func (c Chain[T]) Mapping() []V {
//
//}
