package projpb

import "store/pkg/sdk/helper"

func (t *Commodity) GetRefs() []string {

	if len(t.Medias) > 0 {
		return helper.Mapping(t.Medias, func(x *Media) string {
			return x.Url
		})
	}

	return t.Images
}
