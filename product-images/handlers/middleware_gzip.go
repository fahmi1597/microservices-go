package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// GzipEncoding is where GzipMiddleware sits
type GzipEncoding struct {
}

// GzipMiddleware allows contents in http response to be compressed
func (ge *GzipEncoding) GzipMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			gresp := NewCustomResponseWriter(resp)
			gresp.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(gresp, req)
			defer gresp.Flush()

			return
		}

		next.ServeHTTP(resp, req)
	})
}

// CustomResponseWriter wraps original http.ResponseWriter
type CustomResponseWriter struct {
	resp  http.ResponseWriter
	gresp *gzip.Writer
}

// NewCustomResponseWriter creates a new custom http.ResponseWriter with Gzip Compression
func NewCustomResponseWriter(resp http.ResponseWriter) *CustomResponseWriter {
	gresp := gzip.NewWriter(resp)
	return &CustomResponseWriter{resp: resp, gresp: gresp}
}

// Header implements http.ResponseWriter
func (crw *CustomResponseWriter) Header() http.Header {
	return crw.resp.Header()
}

// Write implements http.ResponseWriter
func (crw *CustomResponseWriter) Write(b []byte) (int, error) {
	// return new compressed data
	return crw.gresp.Write(b)
}

// WriteHeader implements http.ResponseWriter interface
func (crw *CustomResponseWriter) WriteHeader(code int) {
	crw.resp.WriteHeader(code)
}

// Flush implements Flusher and is implemented by http.ResponseWriter
func (crw *CustomResponseWriter) Flush() {
	crw.gresp.Flush()
	crw.gresp.Close()
}
