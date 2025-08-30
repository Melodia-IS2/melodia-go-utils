package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func RequestLogger(flusher Flusher) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			l := &Log{
				ID:        uuid.NewString(),
				Endpoint:  r.URL.Path,
				Method:    r.Method,
				RemoteIP:  r.RemoteAddr,
				UserAgent: r.UserAgent(),
				Timestamp: start,
				Entries:   []Entry{},
			}
			ctx := WithLog(r.Context(), l)

			defer func() {
				l.Status = ww.Status()
				l.Bytes = ww.BytesWritten()
				l.Duration = time.Since(start)

				if flusher != nil {
					if err := flusher.Flush(r.Context(), l); err != nil {
						fmt.Printf("error flushing log: %v\n", err)
					}
				}
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}
