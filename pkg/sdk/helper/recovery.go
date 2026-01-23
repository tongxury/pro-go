package helper

import (
	"context"
	"fmt"
	"google.golang.org/appengine/log"
	"os"
	"runtime"
)

func TryCatch(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()
	f()
	return nil
}

func Recovery(funcs ...func(interface{})) {
	if r := recover(); r != nil {
		recovered := false
		if len(funcs) > 0 {
			for _, fun := range funcs {
				if fun != nil {
					fun(r)
					recovered = true
				}
			}
		}
		if !recovered {
			buf := make([]byte, 1<<18)
			n := runtime.Stack(buf, false)
			log.Errorf(context.Background(), "%v, STACK: %s", r, buf[0:n])

		}
	}
}

func DeferFunc(f ...func()) {
	if r := recover(); r != nil {

		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}

		_, _ = fmt.Fprintf(os.Stderr, "panic: %v\n%s", err, buf)

		//slf.WithError(err).Errorln("[PANIC]...\n" + string(buf))
	}

	if len(f) > 0 {
		f[0]()
	}
}
