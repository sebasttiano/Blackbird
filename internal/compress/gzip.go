package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
)

var compressedTypes = []string{
	"application/json",
	"text/html",
}

// GZIPWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type GZIPWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewGZIPWriter(w http.ResponseWriter) *GZIPWriter {
	return &GZIPWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *GZIPWriter) Header() http.Header {
	return c.w.Header()
}

func (c *GZIPWriter) Write(p []byte) (int, error) {
	for _, t := range compressedTypes {
		if slices.Contains(c.Header().Values("Content-Type"), t) {
			c.WriteHeader(http.StatusOK)
			return c.zw.Write(p)
		}
	}
	return c.w.Write(p)
}

func (c *GZIPWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *GZIPWriter) Close() error {
	return c.zw.Close()
}

// GZIPReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type GZIPReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewZIPReader(r io.ReadCloser) (*GZIPReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &GZIPReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c GZIPReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *GZIPReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
