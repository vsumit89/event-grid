package middlewares

import (
	"bytes"
	"net/http"
	"server/pkg/logger"
	"time"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
	Body       bytes.Buffer
}

func NewResponseWriterWrapper(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

func (w *ResponseWriterWrapper) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriterWrapper) Write(data []byte) (int, error) {
	bytesWritten, err := w.Body.Write(data)
	if err != nil {
		return bytesWritten, err
	}
	return w.ResponseWriter.Write(data)
}

func (w *ResponseWriterWrapper) CopyToOriginalWriter() {
	w.ResponseWriter.WriteHeader(w.StatusCode)
	_, _ = w.ResponseWriter.Write(w.Body.Bytes())
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Info("request received", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)

		wrappedWriter := NewResponseWriterWrapper(w)

		next.ServeHTTP(wrappedWriter, r)

		logger.Info("request completed", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr, "status_code", wrappedWriter.StatusCode, "message", http.StatusText(wrappedWriter.StatusCode), "duration", time.Since(start))
	})
}
