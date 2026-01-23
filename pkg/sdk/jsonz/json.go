package jsonz

import jsoniter "github.com/json-iterator/go"

type API struct {
	jsoniter.API
}

func New() API {
	return API{
		API: jsoniter.Config{
			//IndentionStep:                 0,
			//MarshalFloatWith6Digits:       false,
			EscapeHTML: false,
			//SortMapKeys:                   false,
			//UseNumber:                     false,
			//DisallowUnknownFields:         false,
			//TagKey:                        "",
			//OnlyTaggedField:               false,
			//ValidateJsonRawMessage:        false,
			//ObjectFieldMustBeSimpleString: true,
			//CaseSensitive:                 false,
		}.Froze(),
	}
}

func (t *API) MapToStruct(src any, dest any) error {

	b, err := t.Marshal(src)
	if err != nil {
		return err
	}

	err = t.Unmarshal(b, &dest)
	if err != nil {
		return err
	}
	return nil
}
