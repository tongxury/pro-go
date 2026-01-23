package encoder

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"time"

	responsepb "store/api/public/response"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type Response struct {
	Code    string      `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func WriteResponse(res nethttp.ResponseWriter, req *nethttp.Request, code int, data any, message string) {

	codec, _ := http.CodecForRequest(req, "Accept")

	d := map[string]any{
		"data": data,
	}

	if code != 200 {
		d = map[string]any{
			"code":    code,
			"message": message,
		}
	}

	body, err := codec.Marshal(d)

	if err != nil {
		res.WriteHeader(nethttp.StatusInternalServerError)
		return
	}

	res.WriteHeader(200)
	_, err = res.Write(body)
}

func ResponseEncoder(res nethttp.ResponseWriter, req *nethttp.Request, data interface{}) error {

	switch t := data.(type) {
	case *responsepb.RedirectResponse:
		for _, c := range t.Cookies {
			nethttp.SetCookie(res, &nethttp.Cookie{
				Name:     c.Name,
				Value:    c.Value,
				Domain:   c.Domain,
				Expires:  time.Now().Add(3 * 360 * 24 * time.Hour),
				Path:     "/",
				MaxAge:   int(3 * 360 * 24 * time.Hour),
				SameSite: nethttp.SameSiteLaxMode,
				Secure:   true,
			})

		}

		if t.Url != "" {
			nethttp.Redirect(res, req, t.Url, 302)
			return nil
		}

		res.WriteHeader(200)
	case *responsepb.MapResponse:
		res.Header().Set("Content-Type", "application/json")

		marshal, _ := json.Marshal(t.Data)

		_, err := res.Write(marshal)
		if err != nil {
			return errors.InternalServer(err.Error(), "")
		}
		res.WriteHeader(200)

	case *responsepb.IntMapResponse:
		res.Header().Set("Content-Type", "application/json")

		marshal, _ := json.Marshal(t.Data)

		_, err := res.Write(marshal)
		if err != nil {
			return errors.InternalServer(err.Error(), "")
		}
		res.WriteHeader(200)

	case *responsepb.TextResponse:

		contentType := "text/plain"
		if t.ContentType != "" {
			contentType = t.ContentType
		}

		res.Header().Set("Content-Type", contentType)

		_, err := res.Write([]byte(t.Value))
		if err != nil {
			return errors.InternalServer(err.Error(), "")
		}
		res.WriteHeader(200)

	case *responsepb.HtmlResponse:
		res.Header().Set("Content-Type", "text/html")

		_, err := res.Write([]byte(t.Html))
		if err != nil {
			return errors.InternalServer(err.Error(), "")
		}
		res.WriteHeader(200)

	case *responsepb.FileResponse:

		res.Header().Set("Content-Type", t.Type)
		res.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", t.Filename))

		_, err := res.Write([]byte(t.Content))
		if err != nil {
			return errors.InternalServer(err.Error(), "")
		}
		res.WriteHeader(200)

	default:

		resp := Response{
			Data: data,
		}

		codec, _ := http.CodecForRequest(req, "Accept")

		res.Header().Set("Content-Type", "application/"+codec.Name())
		res.WriteHeader(200)

		body, err := codec.Marshal(resp)
		if err != nil {
			return errors.InternalServer(err.Error(), "")
		}

		_, err = res.Write(body)
		if err != nil {
			return errors.InternalServer(err.Error(), "")
		}
	}

	return nil
}
