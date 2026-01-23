package conv

import (
	"encoding/json"
	yaml2 "github.com/ghodss/yaml"
)

func Y2S[T any](src []byte, dest *T) error {

	// yaml 驼峰字段无法对应 用json中转一下
	jBytes, err := yaml2.YAMLToJSON(src)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jBytes, dest); err != nil {
		return err
	}
	return nil
}
