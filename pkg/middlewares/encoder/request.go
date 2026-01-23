package encoder

import (
	"bytes"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"io"
)

type formCodec struct{}

func (formCodec) Marshal(v interface{}) ([]byte, error) {
	return nil, nil
}

func (formCodec) Unmarshal(data []byte, v interface{}) error {
	return nil
}

func (formCodec) Name() string {
	return "form-data"
}

func RequestDecoder(r *http.Request, v interface{}) error {

	encoding.RegisterCodec(formCodec{})

	codec, ok := http.CodecForRequest(r, "Content-Type")
	log.Debugw("codec", codec)

	if !ok {
		return errors.BadRequest("CODEC", fmt.Sprintf("unregister Content-Type: %s", r.Header.Get("Content-Type")))
	}
	data, err := io.ReadAll(r.Body)

	// reset body.
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	if err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	if len(data) == 0 {
		return nil
	}
	if err = codec.Unmarshal(data, v); err != nil {
		return errors.BadRequest("CODEC", fmt.Sprintf("body unmarshal %s", err.Error()))
	}
	return nil
}
