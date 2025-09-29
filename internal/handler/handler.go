package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kornetvba/YandShortUrl/internal/service"
	"log"
	"net/http"
)

func TextPlainPage(c *gin.Context) {

	if c.Request.Method != http.MethodPost {
		c.String(http.StatusBadRequest, "")
		return
	}

	if c.ContentType() != "text/plain" {
		c.String(http.StatusBadRequest, "")
		return
	}

	postTest, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	hashText, err := service.HashPlainText(postTest)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	resultResText := fmt.Sprintf("%s%s%s%s", "http://", c.Request.Host, c.Request.URL.String(), hashText)

	c.Header("Content-Type", "text/plain")
	c.Header("Content-Length", fmt.Sprint(len(resultResText)))
	c.Writer.WriteHeader(http.StatusCreated)
	_, err = c.Writer.Write([]byte(resultResText))
	if err != nil {
		// Логируем ошибку, но НЕ отправляем новый ответ клиенту
		log.Printf("Failed to write response: %v", err)
		return // Прерываем выполнение
	}

}

func GetTextPlainPage(c *gin.Context) {

	if c.Request.Method != http.MethodGet {
		c.String(http.StatusBadRequest, "")
		return
	}

	//if c.ContentType() != "text/plain" {
	//	c.String(http.StatusBadRequest, "")
	//	return
	//}

	id := c.Param("id")
	resURL, err := service.DeHashText(id)

	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	c.Writer.Header().Set("Location", string(resURL))
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)

}
