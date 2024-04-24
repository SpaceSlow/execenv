package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"
)

const CompressionAlgorithm = "gzip"

var SupportedContentTypes = []string{"application/json", "text/html"}

type compressResponseWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func (w compressResponseWriter) Write(b []byte) (int, error) {
	return w.zw.Write(b)
}

func newCompressResponseWriter(w http.ResponseWriter) *compressResponseWriter {
	return &compressResponseWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func WithCompressing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isContainsCompression := strings.Contains(r.Header.Get("Content-Encoding"), CompressionAlgorithm)
		if isContainsCompression {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		isSupportedContentType := slices.Contains(SupportedContentTypes, w.Header().Get("Content-Type"))
		isMatchCompressionAlgorithm := strings.Contains(r.Header.Get("Accept-Encoding"), CompressionAlgorithm)
		if isSupportedContentType && isMatchCompressionAlgorithm {
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer gz.Close()

			w.Header().Set("Content-Encoding", CompressionAlgorithm)
			w = newCompressResponseWriter(w)
		}

		next.ServeHTTP(w, r)
	})
}
