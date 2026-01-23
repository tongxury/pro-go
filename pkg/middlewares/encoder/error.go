package encoder

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	nethttp "net/http"
)

func ErrorEncoder(res nethttp.ResponseWriter, req *nethttp.Request, err error) {
	codec, _ := http.CodecForRequest(req, "Accept")

	res.Header().Set("Content-Type", "application/"+codec.Name())

	var code = "errorServer"
	var message string
	var statusCode int = 500

	switch t := err.(type) {
	case *errors.Error:

		if t.Reason != "UNKNOWN" {
			code = (t.Reason)
			message = t.Message
			statusCode = 200
		}

	default:
		//err = errors.InternalServer(err.Error(), "")
		//message = err.Error()
	}

	log.Errorw("server err", err)

	resp := Response{
		Code:    code,
		Message: message,
	}

	res.WriteHeader(statusCode)

	body, _ := codec.Marshal(resp)
	_, _ = res.Write(body)

}
