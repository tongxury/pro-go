package helper

func FindInStringMap(src map[string]string, key string) string {
	if len(src) == 0 {
		return ""
	}

	return src[key]
}

func AppendMap[K comparable, V any](dst map[K]V, src ...map[K]V) map[K]V {
	if dst == nil {
		dst = map[K]V{}
	}

	for _, m := range src {
		for k, v := range m {
			dst[k] = v
		}
	}

	return dst
}

func CreateStringMap(params ...string) map[string]string {
	rsp := make(map[string]string)

	for index := 0; index < len(params)-1; index = index + 2 {
		if index%2 == 0 {
			rsp[params[index]] = params[index+1]
		}
	}

	return rsp
}

func HasKey(m map[string]interface{}) bool {
	return len(MapKeys(m)) > 0
}

func ContainsKey(m map[string]interface{}, key string) bool {
	if m == nil {
		return false
	}

	keys := MapKeys(m)

	return Contains(keys, key)
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func MapValues[T any](m map[string]T) []T {
	values := make([]T, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
