package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kornetvba/YandShortUrl/internal/config/compress"
	"github.com/kornetvba/YandShortUrl/internal/config/config"
	"github.com/kornetvba/YandShortUrl/internal/config/db"
	"github.com/kornetvba/YandShortUrl/internal/config/logger"
	"github.com/kornetvba/YandShortUrl/internal/handler"
	"github.com/kornetvba/YandShortUrl/internal/repository"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func run(URLHandler *handler.URLHandler) (*http.Server, error) {

	gin.SetMode(gin.ReleaseMode)
	err := logger.Initialization(config.LevelLog)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("server is running", zap.String("address", config.Addr.String()))

	r := gin.New()
	r.Use(logger.HTTPLoggerMiddleWare())
	r.Use(compress.GzipCompressMiddleWare())

	//r.Use(gz.Gzip(gz.BestCompression))

	r.Use(gin.Recovery())
	r.POST("/", URLHandler.TextPlainPage)
	r.POST("/api/shorten", URLHandler.PostURL)
	r.GET("/:id", URLHandler.GetTextPlainPage)
	r.GET("/ping", URLHandler.PingHandler)

	srv := &http.Server{
		Addr:    config.Addr.Host + ":" + strconv.Itoa(config.Addr.Port),
		Handler: r,
	}
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("Server failed", zap.Error(err))
		}
	}()
	return srv, nil

}

func main() {
	handlerURL := handler.NewURLHandler(repository.NewURLRecords())

	err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	database := db.Database{nil}
	_, err = database.New(config.DatabaseDSN)
	if err != nil {
		log.Print(err)
	}

	err = handlerURL.Storage.LoadRecords(config.FilePath)
	if err != nil {
		log.Print(err)
	}

	srv, err := run(handlerURL)
	if err != nil {
		logger.Log.Error("serv not running", zap.Error(err))
		return
	}

	defer func() {
		logger.Log.Info("Saving data to file...")
		if err = handlerURL.Storage.SaveFile(config.FilePath); err != nil {
			if err.Error() == "Writing/reading to a file is disabled" {
				logger.Log.Info("Writing/reading to a file is disabled")
				return
			}
			logger.Log.Error("Save file failed", zap.Error(err))
		} else {
			logger.Log.Info("Data saved successfully")
		}
	}()

	// ИСПРАВЛЕНИЕ: буферизованный канал
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Log.Info("Received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server shutdown error", zap.Error(err))
	}

	logger.Log.Info("Application stopped gracefully")
}
