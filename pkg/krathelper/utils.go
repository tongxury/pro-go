package krathelper

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"io"
)

type FormFile struct {
	Body     []byte
	Filename string
	Size     int64

	URL string // fill after upload
	Md5 string // fill after upload
}

type FormFiles []*FormFile

// todo 放到 form-data codec中实现1
func GetMultipartFormFile(ctx context.Context, name, name2 string) ([]*FormFile, error) {

	request, ok := http.RequestFromServerContext(ctx)
	if !ok {
		return nil, errors.New("RequestFromServerContext err")
	}

	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, err
	}
	formFiles := request.MultipartForm.File[name]

	if len(formFiles) == 0 {
		formFiles = request.MultipartForm.File[name2]

		if len(formFiles) == 0 {
			return nil, errors.New("form files err")
		}
	}

	var files []*FormFile
	for _, x := range formFiles {
		fileBody, _ := x.Open()

		fileBytes, err := io.ReadAll(fileBody)
		if err != nil {
			return nil, err
		}

		files = append(files, &FormFile{
			Body:     fileBytes,
			Filename: x.Filename,
			Size:     int64(len(fileBytes)),
		})
	}

	return files, nil
}
