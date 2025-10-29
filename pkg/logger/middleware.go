package logger

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func RequestLogger(flusher Flusher, appName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			l := &Log{
				ID:        uuid.NewString(),
				AppName:   appName,
				Endpoint:  r.URL.Path,
				Method:    r.Method,
				RemoteIP:  r.RemoteAddr,
				UserAgent: r.UserAgent(),
				Timestamp: start,
				Entries:   []Entry{},
			}
			ctx := WithLog(r.Context(), l)

			defer func() {
				shouldPanic := false
				rec := recover()
				if rec != nil {
					l.Status = http.StatusInternalServerError
					l.Entries = append(l.Entries, Entry{
						Level:   Error,
						Message: fmt.Sprintf("panic: %v", rec),
						Data:    string(debug.Stack()),
					})
					shouldPanic = true
				}

				l.Status = ww.Status()
				l.Bytes = ww.BytesWritten()
				l.Duration = time.Since(start)

				if flusher != nil {
					if err := flusher.Flush(r.Context(), l); err != nil {
						fmt.Printf("error flushing log: %v\n", err)
					}
				}

				if shouldPanic {
					panic(rec)
				}
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}
