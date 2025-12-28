package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

var Log = zap.NewNop()

func Initialization(level string) error {

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	lg, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = lg
	return nil
}

func HTTPLoggerMiddleWare() gin.HandlerFunc {
	fnHan := func(hand *gin.Context) {

		timeStart := time.Now()

		hand.Next()

		duration := time.Since(timeStart)

		Log.Info("HTTP completed",
			zap.String(" URI", hand.Request.RequestURI),
			zap.String("method", hand.Request.Method),
			zap.Duration("time", duration),
			zap.Int("status", hand.Writer.Status()),
			zap.Int("size", hand.Writer.Size()),
		)

	}
	return fnHan

}
