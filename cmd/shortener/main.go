package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kornetvba/YandShortUrl/internal/config/compress"
	"github.com/kornetvba/YandShortUrl/internal/config/config"
	"github.com/kornetvba/YandShortUrl/internal/config/logger"
	"github.com/kornetvba/YandShortUrl/internal/handler"
	"go.uber.org/zap"
	"log"
)

func init() {
	err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	gin.SetMode(gin.ReleaseMode)
	err := logger.Initialization(config.LevelLog)
	if err != nil {
		return err
	}
	logger.Log.Info("server is running", zap.String("address", config.Addr.String()))

	r := gin.New()
	r.Use(logger.HTTPLoggerMiddleWare())
	r.Use(compress.GzipCompressMiddleWare())
	//r.Use(gz.Gzip(gz.BestCompression))

	r.Use(gin.Recovery())
	r.POST("/", handler.TextPlainPage)
	r.POST("/api/shorten", handler.PostURL)
	r.GET("/:id", handler.GetTextPlainPage)

	return r.Run(fmt.Sprintf(":%d", config.Addr.Port))

}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
