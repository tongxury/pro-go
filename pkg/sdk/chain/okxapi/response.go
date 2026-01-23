package okxapi

type Response[T any] struct {
	Code string `json:"code"`
	Data T      `json:"data"`
	Msg  string `json:"msg"`
}
