package compress

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"github.com/kornetvba/YandShortUrl/internal/config/logger"
	"go.uber.org/zap"
	"io"
	"strings"
)

type gzWriter struct {
	gin.ResponseWriter
	writer io.Writer
}

func (gz *gzWriter) Write(data []byte) (int, error) {
	return gz.writer.Write(data)
}

type gzReader struct {
	r  io.ReadCloser
	gz *gzip.Reader
}

func (g gzReader) Read(d []byte) (int, error) {
	return g.gz.Read(d)
}

func (g gzReader) Close() error {
	return g.gz.Close()
}

func unCompressData(r io.ReadCloser) (*gzReader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &gzReader{gz: gz, r: r}, nil

}

func GzipCompressMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		if strings.Contains(context.GetHeader("Accept-Encoding"), "gzip") {
			context.Writer.Header().Set("Content-Encoding", "gzip")
			gz, err := gzip.NewWriterLevel(context.Writer, gzip.BestCompression)
			if err != nil {
				logger.Log.Error("compress data", zap.Error(err))
				return
			}
			defer gz.Close()
			gzipWriter := gzWriter{ResponseWriter: context.Writer, writer: gz}
			context.Writer = &gzipWriter

		}
		if strings.Contains(context.GetHeader("Content-Encoding"), "gzip") {

			body, err := unCompressData(context.Request.Body)
			if err != nil {
				logger.Log.Error("uncompress data", zap.Error(err))
			}

			context.Request.Body = body

		}
		context.Next()

	}
}
