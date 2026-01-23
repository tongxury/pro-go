package codecs

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"reflect"
)

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/textproto"
)

func Register() {
	log.Debugw("Register codec")
	encoding.RegisterCodec(formCodec{})
}

type formCodec struct{}

func (formCodec) Marshal(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	val := reflect.ValueOf(v)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		partHeader := make(textproto.MIMEHeader)
		partHeader.Add("Content-Disposition", `form-data; name="`+field.Name+`"`)
		part, err := writer.CreatePart(partHeader)
		if err != nil {
			return nil, err
		}

		_, err = part.Write([]byte(value.String()))
		if err != nil {
			return nil, err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (formCodec) Unmarshal(data []byte, v interface{}) error {
	reader := multipart.NewReader(bytes.NewReader(data), "boundary")

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(part)
		if err != nil {
			return err
		}

		field := reflect.ValueOf(v).Elem().FieldByName(part.FormName())
		if field.IsValid() {
			field.SetString(buf.String())
		}
	}

	return nil
}

func (formCodec) Name() string {
	return "multipart/form-data"
}
