package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kornetvba/YandShortUrl/internal/config/config"
	"github.com/kornetvba/YandShortUrl/internal/config/db"
	"github.com/kornetvba/YandShortUrl/internal/repository"
	"github.com/kornetvba/YandShortUrl/internal/service"
	"net/http"
	"time"
)

type URLHandler struct {
	Storage repository.URLStorage
}

func NewURLHandler(storage repository.URLStorage) *URLHandler {
	return &URLHandler{Storage: storage}
}

func (mh *URLHandler) TextPlainPage(c *gin.Context) {

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

	err = mh.Storage.AppendRecord(hashText, string(postTest))
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	resultResText := fmt.Sprintf("%s%s%s", config.ResultURL, "/", hashText)

	c.Header("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Write([]byte(resultResText))

}

func (mh *URLHandler) GetTextPlainPage(c *gin.Context) {

	id := c.Param("id")

	resURL, err := mh.Storage.GetRecord(id)

	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "")
		return
	}
	c.Writer.Header().Set("Location", resURL.OriginalURL)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)

}

func (mh *URLHandler) CreateShortURL(c *gin.Context) {
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

	shortURL, err := service.HashPlainText([]byte(reqU.URL))
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	err = mh.Storage.AppendRecord(shortURL, reqU.URL)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	resU.Result = fmt.Sprintf("%s%s%s", config.ResultURL, "/", shortURL)
	resp, err := json.Marshal(resU)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(resp)

}

func (mh *URLHandler) PingDB(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := db.DB.PingContext(ctx)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func (mh *URLHandler) CreateShortURLBatch(c *gin.Context) {

	if c.GetHeader("Content-type") != "application/json" {
		c.String(http.StatusUnsupportedMediaType, "")
		return
	}
	//helo
	records := make([]repository.URLRecord, 0, 20)

	data, err := c.GetRawData()

	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&records); err != nil {
		c.String(http.StatusUnprocessableEntity, "")
		return
	}

	if len(records) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	resData, err := mh.Storage.AppendRecords(&records)

	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	type responseData struct {
		ID       string `json:"correlation_id"`
		ShortURL string `json:"short_url"`
	}
	var respItems []responseData

	for _, r := range resData {
		respItems = append(respItems, responseData{
			ID:       r.CorrelationID,
			ShortURL: r.ShortURL,
		})
	}

	c.JSON(http.StatusCreated, respItems)

}
