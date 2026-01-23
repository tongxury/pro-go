package entz

type Extra map[string]interface{}

func (t Extra) GetString(name string) string {

	v, ok := t[name]
	if !ok {
		return ""
	}

	return v.(string)
}
