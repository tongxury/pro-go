package helper

import "fmt"

func Combine[T any](iterable []T, r int) [][]T {

	var rsp [][]T

	pool := iterable
	n := len(pool)

	if r > n {
		return [][]T{iterable}
	}

	indices := make([]int, r)
	for i := range indices {
		indices[i] = i
	}

	result := make([]T, r)
	for i, el := range indices {
		result[i] = pool[el]
	}

	rsp = app(rsp, result)

	for {
		i := r - 1
		for ; i >= 0 && indices[i] == i+n-r; i -= 1 {
		}

		if i < 0 {
			return rsp
		}

		indices[i] += 1
		for j := i + 1; j < r; j += 1 {
			indices[j] = indices[j-1] + 1
		}

		for ; i < len(indices); i += 1 {
			result[i] = pool[indices[i]]
		}
		rsp = app(rsp, result)
	}
}

func app[T any](dst [][]T, src []T) [][]T {

	var tmp []T
	for _, i := range src {
		tmp = append(tmp, i)
	}

	return append(dst, tmp)

}

func combinations(iterable []int, r int) {
	pool := iterable
	n := len(pool)

	if r > n {
		return
	}

	indices := make([]int, r)
	for i := range indices {
		indices[i] = i
	}

	result := make([]int, r)
	for i, el := range indices {
		result[i] = pool[el]
	}

	fmt.Println(result)

	for {
		i := r - 1
		for ; i >= 0 && indices[i] == i+n-r; i -= 1 {
		}

		if i < 0 {
			return
		}

		indices[i] += 1
		for j := i + 1; j < r; j += 1 {
			indices[j] = indices[j-1] + 1
		}

		for ; i < len(indices); i += 1 {
			result[i] = pool[indices[i]]
		}
		fmt.Println(result)

	}

}

func permutations(iterable []int, r int) {
	pool := iterable
	n := len(pool)

	if r > n {
		return
	}

	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}

	cycles := make([]int, r)
	for i := range cycles {
		cycles[i] = n - i
	}

	result := make([]int, r)
	for i, el := range indices[:r] {
		result[i] = pool[el]
	}

	fmt.Println(result)

	for n > 0 {
		i := r - 1
		for ; i >= 0; i -= 1 {
			cycles[i] -= 1
			if cycles[i] == 0 {
				index := indices[i]
				for j := i; j < n-1; j += 1 {
					indices[j] = indices[j+1]
				}
				indices[n-1] = index
				cycles[i] = n - i
			} else {
				j := cycles[i]
				indices[i], indices[n-j] = indices[n-j], indices[i]

				for k := i; k < r; k += 1 {
					result[k] = pool[indices[k]]
				}

				fmt.Println(result)

				break
			}
		}

		if i < 0 {
			return
		}

	}

}

func product(argsA, argsB []int) {

	pools := [][]int{argsA, argsB}
	npools := len(pools)
	indices := make([]int, npools)

	result := make([]int, npools)
	for i := range result {
		result[i] = pools[i][0]
	}

	fmt.Println(result)

	for {
		i := npools - 1
		for ; i >= 0; i -= 1 {
			pool := pools[i]
			indices[i] += 1

			if indices[i] == len(pool) {
				indices[i] = 0
				result[i] = pool[0]
			} else {
				result[i] = pool[indices[i]]
				break
			}

		}

		if i < 0 {
			return
		}

		fmt.Println(result)
	}
}
