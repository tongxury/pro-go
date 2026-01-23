package eventsource

import (
	"context"
	"fmt"
	"io"
	"store/pkg/confcenter"
	"store/pkg/krathelper"
	"store/pkg/sdk/conv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type Ctx struct {
	context.Context
}

//
//func (emptyCtx) Deadline() (deadline time.Time, ok bool) {
//	return
//}
//
//func (emptyCtx) Done() <-chan struct{} {
//	return nil
//}
//
//func (emptyCtx) Err() error {
//	return nil
//}
//
//func (emptyCtx) Value(key any) any {
//	return nil
//}

func (t Ctx) UserId() string {
	return t.Value("d").(map[string]any)["userId"].(string)
}

func (t Ctx) write(chunk *Chunk) error {

	writer := t.Value("d").(map[string]any)["writer"].(http.ResponseWriter)

	_, err := fmt.Fprint(writer, chunk.String())
	writer.(http.Flusher).Flush()
	return err
}

func (t Ctx) Write(data any) error {
	return t.write(NewChunk("", 0, data, ""))
}

func (t Ctx) WriteV2(id string, data any) error {
	return t.write(NewChunk(id, 0, data, ""))
}

func (t Ctx) WriteDone(data any) error {
	return t.write(NewChunk("", 1000, data, ""))
}

func (t Ctx) Abort(code int, message string) error {
	return t.write(NewChunk("", code, nil, message))
}

type HandleFunc[T any] func(ctx Ctx, params T) (int, error)

//func HandleV2(f EventSourceHandleFunc[T]) func(w http.ResponseWriter, req *http.Request) {
//	return func(w http.ResponseWriter, req *http.Request) {
//
//		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
//		w.Header().Set("Content-Type", "text/event-stream")
//		w.Header().Set("Cache-Control", "no-cache")
//		w.Header().Set("Connection", "keep-alive")
//
//		bodyBytes, err := io.ReadAll(req.Body)
//		if err != nil {
//			_, _ = fmt.Fprintf(w, NewChunk(400, nil, "").String())
//			return
//		}
//
//		var params T
//		_ = conv.J2S(bodyBytes, &params)
//
//		log.Debugw("params", params)
//
//		// auth
//		ctx := req.Context()
//
//		token := krathelper.GetHeader(ctx, "Authorization")
//		if token == "" {
//			token = krathelper.FindCookie(ctx, confcenter.AuthCookieName)
//		}
//
//		if token == "" {
//			w.WriteHeader(401)
//			return
//		}
//
//		userID, err := krathelper.ParseUserID(token)
//		if err != nil {
//			w.WriteHeader(401)
//			return
//		}
//
//		ectx := Ctx{
//			Context: context.WithValue(context.Background(), "d", map[string]any{"writer": w, "userId": userID}),
//		}
//
//		err = f(ectx, params)
//		if err != nil {
//			_, _ = fmt.Fprintf(w, NewChunk(1500, nil, "").String())
//			return
//		}
//	}
//}

func Handle[T any](f HandleFunc[T]) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			_, _ = fmt.Fprintf(w, NewChunk("", 400, nil, "").String())
			return
		}

		var params T
		_ = conv.J2S(bodyBytes, &params)

		log.Debugw("params", params)

		// auth
		ctx := req.Context()

		token := krathelper.GetHeader(ctx, "Authorization")
		if token == "" {
			token = krathelper.FindCookie(ctx, confcenter.AuthCookieName)
		}

		if token == "" {
			w.WriteHeader(401)
			return
		}

		userID, err := krathelper.ParseUserID(token)
		if err != nil {
			w.WriteHeader(401)
			return
		}

		ectx := Ctx{
			Context: context.WithValue(context.Background(), "d", map[string]any{"writer": w, "userId": userID}),
		}

		code, err := f(ectx, params)
		if err != nil {
			_, _ = fmt.Fprintf(w, NewChunk("", code, nil, "").String())
			return
		}
	}
}
