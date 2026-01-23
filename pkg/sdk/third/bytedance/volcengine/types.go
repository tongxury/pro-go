package volcengine

type Response[T any] struct {
	Code             int    `json:"Code"`
	Message          string `json:"Message"`
	ResponseMetadata struct {
		RequestId string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Service   string `json:"Service"`
		Region    string `json:"Region"`
		Code      int    `json:"Code"`
		Domain    string `json:"Domain"`
		Message   string `json:"Message"`
	} `json:"ResponseMetadata"`
	Result T `json:"Result"`
}

type Owner struct {
	Type string `json:"Type"`
	Id   int64  `json:"Id"`
}

type StoreItem struct {
	Md5              string `json:"Md5"`
	Size             int64  `json:"Size"`
	SkipDataComplete bool   `json:"SkipDataComplete"`
	Filename         string `json:"Filename"`
	FileExtension    string `json:"FileExtension"`
}

type CreateMaterialInfo struct {
	Title              string   `json:"Title"`
	MediaType          int      `json:"MediaType"`
	MediaFirstCategory string   `json:"MediaFirstCategory"`
	Tags               []string `json:"Tags"`
	MediaExtension     string   `json:"MediaExtension"`
}

type GetUploadStateRequest struct {
	Owner *Owner `json:"Owner"`
	Md5   string `json:"Md5"`
	Crc   uint32 `json:"Crc"`
	Size  int64  `json:"Size"`
	Start int64  `json:"Start"`
	End   int64  `json:"End"`
}

type GetUploadStateResult struct {
	SkipDataComplete bool `json:"SkipDataComplete"`
}

type CreateMaterialRequest struct {
	Owner              *Owner              `json:"Owner"`
	StoreItem          *StoreItem          `json:"StoreItem"`
	CreateMaterialInfo *CreateMaterialInfo `json:"CreateMaterialInfo"`
}

type CreateMaterialResult struct {
	MediaId string `json:"MediaId"`
}
