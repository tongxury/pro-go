package conv

func SliceToMap[T comparable](slice []T) map[T]bool {

	tmp := make(map[T]bool, len(slice))

	for _, x := range slice {
		tmp[x] = true
	}

	return tmp
}
