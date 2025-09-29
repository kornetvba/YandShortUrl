package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kornetvba/YandShortUrl/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestTextPlainPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	TableTests := []struct {
		name        string
		httpMethod  string
		body        string
		resBody     string
		statusCode  int
		contentType string
	}{
		{
			name:        "Test1",
			httpMethod:  http.MethodGet,
			body:        "",
			resBody:     "",
			statusCode:  400,
			contentType: "",
		},
		{
			name:        "Test2",
			httpMethod:  http.MethodPost,
			body:        "http://htgfnn.yandex/nubcnadqasd",
			resBody:     "http://example.com/b8d369a6",
			statusCode:  201,
			contentType: "text/plain",
		},
		{
			name:        "Test3",
			httpMethod:  http.MethodPost,
			body:        "http://youtube.com",
			resBody:     "http://example.com/9cd9263b",
			statusCode:  201,
			contentType: "text/plain",
		},
		{
			name:        "Test4",
			httpMethod:  http.MethodPost,
			body:        "http://youtube.com",
			resBody:     "http://example.com/9cd9263b",
			statusCode:  400,
			contentType: "json",
		},
		{
			name:        "Test5",
			httpMethod:  http.MethodPost,
			body:        "http://youtube.com",
			resBody:     "http://example.com/9cd9263b",
			statusCode:  400,
			contentType: "json",
		},
	}

	for _, test := range TableTests {
		t.Run(test.name, func(t *testing.T) {

			req := httptest.NewRequest(test.httpMethod, "/", strings.NewReader(test.body))
			req.Header.Set("Content-Type", test.contentType)

			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req
			TextPlainPage(c)
			res := w.Result()
			//status
			require.Equal(t, test.statusCode, res.StatusCode)
			if res.StatusCode == http.StatusCreated {
				//content-type
				assert.Equal(t, test.contentType, res.Header.Get("Content-Type"))
				//body
				defer func() {
					err := res.Body.Close()
					if err != nil {
						log.Printf("body close error: %v", err)
					}

				}()
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, test.resBody, string(resBody))
				//content-len
				assert.Equal(t, strconv.Itoa(len(test.resBody)), res.Header.Get("Content-Length"))

			}
		})
	}

}

func TestGetTextPlainPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	TableTests := []struct {
		name       string
		httpMethod string
		statusCode int
		pathID     string
		bodyURL    string
	}{
		{
			name:       "test1",
			httpMethod: http.MethodGet,
			statusCode: http.StatusTemporaryRedirect,
			pathID:     "b8d369a6",
			bodyURL:    "http://htgfnn.yandex/nubcnadqasd",
		},
		{
			name:       "test2",
			httpMethod: http.MethodGet,
			statusCode: http.StatusBadRequest,
			pathID:     "b8d36932a6",
			bodyURL:    "http://htgfnn.yandex/nubcnadqasd",
		},
	}

	for _, test := range TableTests {

		t.Run(test.name, func(t *testing.T) {
			_, _ = service.HashPlainText([]byte(test.bodyURL))
			router := gin.New()
			router.GET("/:id", GetTextPlainPage)

			req := httptest.NewRequest(test.httpMethod, fmt.Sprintf("/%s", test.pathID), nil)
			req.Header.Set("Content-Type", "text/plain")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			res := w.Result()
			defer func() {
				err := res.Body.Close()
				if err != nil {
					log.Printf("Body close error: %v", err)
				}
			}()
			require.Equal(t, test.statusCode, res.StatusCode)
			if res.StatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, res.Header.Get("Location"), test.bodyURL)
			}
		})
	}

}
