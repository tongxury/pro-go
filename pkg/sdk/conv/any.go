package conv

func AnySlice[T any](values []T) []any {

	rsp := []any{}

	for i := range values {
		v := values[i]
		rsp = append(rsp, v)
	}

	return rsp
}
