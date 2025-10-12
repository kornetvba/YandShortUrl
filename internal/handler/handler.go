package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kornetvba/YandShortUrl/internal/config/config"
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

	resultResText := fmt.Sprintf("%s%s%s", config.ResultURL, "/", hashText)

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

func PostURL(c *gin.Context) {
	type requestURL struct {
		URL string `json:"url"`
	}
	type responseURL struct {
		Result string `json:"result"`
	}
	var reqU requestURL
	var resU responseURL
	if c.GetHeader("Content-type") != "application/json" {
		c.String(http.StatusInternalServerError, "")
		return
	}

	//	if err := c.ShouldBindJSON(&reqU); err != nil {
	//	c.String(http.StatusInternalServerError, "")
	//	return
	//}
	data, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	err = json.Unmarshal(data, &reqU)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	res, err := service.HashPlainText([]byte(reqU.URL))
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	resU.Result = fmt.Sprintf("%s%s%s", config.ResultURL, "/", res)
	resp, err := json.Marshal(resU)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(resp)

}
