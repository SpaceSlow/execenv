package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"
)

var (
	_ http.ResponseWriter = (*compressResponseWriter)(nil)
	_ io.ReadCloser       = (*compressReader)(nil)
)

const CompressionAlgorithm = "gzip"

var SupportedContentTypes = []string{"application/json", "text/html"}

type compressResponseWriter struct {
	http.ResponseWriter
	compressWriter *gzip.Writer
}

func newCompressResponseWriter(w http.ResponseWriter) *compressResponseWriter {
	return &compressResponseWriter{
		ResponseWriter: w,
		compressWriter: gzip.NewWriter(w),
	}
}

func (w compressResponseWriter) Write(b []byte) (int, error) {
	isSupportedContentType := slices.Contains(SupportedContentTypes, w.Header().Get("Content-Type"))
	if isSupportedContentType {
		return w.compressWriter.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w compressResponseWriter) WriteHeader(statusCode int) {
	isSupportedContentType := slices.Contains(SupportedContentTypes, w.Header().Get("Content-Type"))
	if statusCode < 300 && isSupportedContentType {
		w.Header().Set("Content-Encoding", CompressionAlgorithm)
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

type compressReader struct {
	io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		ReadCloser: r,
		zr:         zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.ReadCloser.Close(); err != nil {
		return err
	}
	return c.zr.Close()
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

		isMatchCompressionAlgorithm := strings.Contains(r.Header.Get("Accept-Encoding"), CompressionAlgorithm)
		if isMatchCompressionAlgorithm {
			cw := newCompressResponseWriter(w)
			w = cw
			defer cw.compressWriter.Close()
		}

		next.ServeHTTP(w, r)
	})
}
