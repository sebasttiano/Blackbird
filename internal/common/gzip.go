package common

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/klauspost/compress/gzip"

	"github.com/sebasttiano/Blackbird.git/internal/logger"
)

// compressedTypes типы Сontent-Type, которые разрешенно сжимать
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

// NewGZIPWriter конструктор для GZIPWriter
func NewGZIPWriter(w http.ResponseWriter) *GZIPWriter {
	return &GZIPWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header метод возращает карту HTTP заголовков
func (c *GZIPWriter) Header() http.Header {
	return c.w.Header()
}

// Write отправляет сжатые данные
func (c *GZIPWriter) Write(p []byte) (int, error) {
	for _, t := range compressedTypes {
		if slices.Contains(c.Header().Values("Content-Type"), t) {
			c.WriteHeader(http.StatusOK)
			return c.zw.Write(p)
		}
	}
	return c.w.Write(p)
}

// WriteHeader добавляет статус код в заголовки
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

// NewZIPReader конструктор для GZIPReader
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

// Read считывает сжатые данные
func (c GZIPReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close закрывает GZIPReader
func (c *GZIPReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// Compress сжимает слайс байт.
func Compress(data []byte) (*bytes.Buffer, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return &b, nil
}

// Decompress распаковывает слайс байт.
func Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		logger.Log.Error("failed to init gzip reader")
	}
	defer r.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}
