package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kornetvba/YandShortUrl/internal/handler"
)

func run() error {
	gin.SetMode(gin.ReleaseMode)
	//mux := http.NewServeMux()
	//mux.HandleFunc("/", handler.TextPlainPage)
	//	mux.HandleFunc("/{id}", handler.GetTextPlainPage)

	//return http.ListenAndServe(":8080", mux)

	r := gin.Default()
	r.POST("/", handler.TextPlainPage)
	r.GET("/:id", handler.GetTextPlainPage)

	return r.Run()

}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
