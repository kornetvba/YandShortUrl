package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kornetvba/YandShortUrl/internal/config/config"
	"github.com/kornetvba/YandShortUrl/internal/handler"
	"log"
)

func init() {
	config.ParseFlags()
}

func run() error {

	gin.SetMode(gin.ReleaseMode)
	//mux := http.NewServeMux()
	//mux.HandleFunc("/", handler.TextPlainPage)
	//	mux.HandleFunc("/{id}", handler.GetTextPlainPage)

	//return http.ListenAndServe(":8080", mux)

	r := gin.Default()
	r.POST("/", handler.TextPlainPage)
	r.GET("/:id", handler.GetTextPlainPage)
	fmt.Println("port in", config.Addr.Port)

	return r.Run(fmt.Sprintf(":%d", config.Addr.Port))

}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
