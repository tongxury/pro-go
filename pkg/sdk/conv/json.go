package conv

import "encoding/json"

func J2M[K comparable, V any](src string, dest *map[K]V) error {
	if err := json.Unmarshal([]byte(src), dest); err != nil {
		return err
	}
	return nil
}

func B2M[K comparable, V any](src []byte) (map[K]V, error) {
	var dest map[K]V

	if err := json.Unmarshal(src, &dest); err != nil {
		return nil, err
	}
	return dest, nil
}

func MS2M[K comparable, V any](src any) map[K]V {

	var dest map[K]V

	if err := json.Unmarshal(S2B(src), &dest); err != nil {
		panic(err)
	}
	return dest
}

func S2M[K comparable, V any](src any) (map[K]V, error) {

	var dest map[K]V

	if err := json.Unmarshal(S2B(src), &dest); err != nil {
		return nil, err
	}
	return dest, nil
}

func B2S[T any](src []byte, dest *T) error {
	if err := json.Unmarshal(src, dest); err != nil {
		return err
	}
	return nil
}

func S2J(src any) string {
	jsonBytes, err := json.Marshal(src)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

func S2B(src any) []byte {
	bytes, _ := json.Marshal(src)
	return bytes
}

func M2J[K comparable, V any](src map[K]V) string {
	bytes, _ := json.Marshal(src)
	return string(bytes)
}

func M2B[K comparable, V any](src map[K]V) []byte {
	bytes, _ := json.Marshal(src)
	return bytes
}

func J2S[T any](jsonBytes []byte, targetStruct *T) error {
	if len(jsonBytes) == 0 {
		return nil
	}
	if err := json.Unmarshal(jsonBytes, targetStruct); err != nil {
		return err
	}
	return nil
}
