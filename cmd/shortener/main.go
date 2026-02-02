package main

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/kornetvba/YandShortUrl/internal/backup"
	"github.com/kornetvba/YandShortUrl/internal/backup/dbmanager"
	"github.com/kornetvba/YandShortUrl/internal/backup/memorymanager"
	"github.com/kornetvba/YandShortUrl/internal/config/compress"
	"github.com/kornetvba/YandShortUrl/internal/config/config"
	"github.com/kornetvba/YandShortUrl/internal/config/db"
	"github.com/kornetvba/YandShortUrl/internal/config/logger"
	"github.com/kornetvba/YandShortUrl/internal/handler"
	"github.com/kornetvba/YandShortUrl/internal/repository/memory"
	"github.com/kornetvba/YandShortUrl/internal/repository/pg"
	_ "github.com/lib/pq"
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

	r.Use(gin.Recovery())
	r.POST("/", URLHandler.TextPlainPage)
	r.POST("/api/shorten", URLHandler.CreateShortURL)
	r.GET("/:id", URLHandler.GetTextPlainPage)
	r.GET("/ping", URLHandler.PingDB)
	r.POST("/api/shorten/batch", URLHandler.CreateShortURLBatch)

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
	err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	var handlerURL *handler.URLHandler
	var backupStorage *backup.BackupManager

	if config.DatabaseDSN != "" {
		con, err := sql.Open("postgres", config.DatabaseDSN)
		db.DB = con
		if err != nil {
			log.Fatal(err)
		}

		defer con.Close()
		if err := con.Ping(); err != nil {
			log.Fatal(err)
		}
		store := pg.NewStore(con)
		//Создаем хендлер с хранилищем
		handlerURL = handler.NewURLHandler(store)
		//Создаем бэкап нашего хранилища
		backupStore := dbmanager.NewBackupStorage(store)
		backupStorage = backup.NewBackupManager(backupStore)

		//Создаем таблицы
		err = store.BootStrap()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		//Оно не надо, но тз сказало так)
		con, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=short_url sslmode=disable")
		db.DB = con
		if err != nil {
			log.Fatal(err)
		}
		defer con.Close()
		store := memory.NewURLRecords()
		handlerURL = handler.NewURLHandler(store)
		backupMemory := memorymanager.NewBackupStorage(store)
		backupStorage = backup.NewBackupManager(backupMemory)
	}

	//Загружаем данные из файла
	err = backupStorage.DownloadRecordsToFile(config.FilePath)
	if err != nil {
		log.Print(err)
	}
	//Запускаем сервер
	srv, err := run(handlerURL)
	if err != nil {
		logger.Log.Error("serv not running", zap.Error(err))
		return
	}
	//Горутина для сохранения файла, graceful shutdown
	defer func() {
		logger.Log.Info("Saving data to file...")
		if err = backupStorage.SaveDataFile(config.FilePath); err != nil {
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
