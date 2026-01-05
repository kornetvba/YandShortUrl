package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kornetvba/YandShortUrl/internal/repository/memory"
	"github.com/kornetvba/YandShortUrl/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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
			resBody:     "/b8d369a6",
			statusCode:  201,
			contentType: "text/plain",
		},
		{
			name:        "Test3",
			httpMethod:  http.MethodPost,
			body:        "http://youtube.com",
			resBody:     "/9cd9263b",
			statusCode:  201,
			contentType: "text/plain",
		},
		{
			name:        "Test4",
			httpMethod:  http.MethodPost,
			body:        "http://youtube.com",
			resBody:     "/9cd9263b",
			statusCode:  400,
			contentType: "json",
		},
		{
			name:        "Test5",
			httpMethod:  http.MethodPost,
			body:        "http://youtube.com",
			resBody:     "/9cd9263b",
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
			handlerURL := NewURLHandler(memory.NewURLRecords())
			handlerURL.TextPlainPage(c)
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
			pathID:     "b",
			bodyURL:    "http://htex/nubcnadqasd",
		},
	}

	for _, test := range TableTests {

		t.Run(test.name, func(t *testing.T) {
			handlerURL := NewURLHandler(memory.NewURLRecords())
			hashURL, err := service.HashPlainText([]byte(test.bodyURL))
			require.NoError(t, err)

			_ = handlerURL.Storage.AppendRecord(hashURL, test.bodyURL)

			router := gin.New()
			router.GET("/:id", handlerURL.GetTextPlainPage)

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

func TestPostURL(t *testing.T) {
	type TestRequestURL struct {
		URL string `json:"url"`
	}
	type TestResponseURL struct {
		Result string `json:"result"`
	}
	respURL := TestResponseURL{}
	tableTest := []struct {
		name        string
		reqBody     TestRequestURL
		want        TestResponseURL
		status      int
		contentType string
		method      string
	}{
		{
			name:        "Test1",
			reqBody:     TestRequestURL{URL: "youtube.com321"},
			want:        TestResponseURL{Result: "/a80f7ccb"},
			status:      http.StatusCreated,
			contentType: "application/json",
			method:      http.MethodPost,
		},
		{
			name:        "Test2",
			reqBody:     TestRequestURL{URL: "youtube.com3221"},
			want:        TestResponseURL{Result: "/a80f7ccb"},
			status:      http.StatusNotFound,
			contentType: "application/json",
			method:      http.MethodGet,
		},
		{
			name:        "Test3",
			reqBody:     TestRequestURL{URL: "youtube.com321"},
			want:        TestResponseURL{Result: "/a80f7ccb"},
			status:      http.StatusCreated,
			contentType: "application/json",
			method:      http.MethodPost,
		},
	}
	for _, ts := range tableTest {
		t.Run(ts.name, func(t *testing.T) {

			data, err := json.Marshal(ts.reqBody)
			if err != nil {
				log.Fatal(err)
			}

			req, err := http.NewRequest(ts.method, "/", bytes.NewReader(data))
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("Content-Type", ts.contentType)

			router := gin.New()
			handlerURL := NewURLHandler(memory.NewURLRecords())
			router.POST("/", handlerURL.PostURL)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			res := w.Result()
			defer func() {
				err = res.Body.Close()
				if err != nil {
					log.Printf("Body close error: %v", err)
				}
			}()

			//	req.Header.Set("Content-Type", ts.contentType)
			//	w := httptest.NewRecorder()
			//	c, _ := gin.CreateTestContext(w)
			//	c.Request = req
			//	PostUrl(c)
			//	res := w.Result()

			assert.Equal(t, ts.status, res.StatusCode)
			if res.StatusCode == http.StatusCreated {
				if err = json.NewDecoder(res.Body).Decode(&respURL); err != nil {
					log.Fatal(err)
				}
				assert.Equal(t, ts.want, respURL)

			}
		})
	}
}

//func TestGetTextPlainMock(t *testing.T) {
//	tableTests := []struct {
//		name       string
//		id         string
//		statusCode int
//		resp       *repository.URLRecord
//	}{
//		{
//			name:       "test1",
//			id:         "testIDRecord",
//			statusCode: http.StatusTemporaryRedirect,
//			resp: &repository.URLRecord{
//				CorrelationID: uuid.New().String(),
//				ShortURL:      "shortURL",
//				OriginalURL:   "originalURl",
//			},
//		},
//	}
//
//	for _, tt := range tableTests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//
//			m := mock_repository.NewMockURLStorage(ctrl)
//			m.EXPECT().GetRecord(tt.id).Return(tt.resp, nil)
//			app := NewURLHandler(m)
//
//			router := gin.New()
//			router.GET("/:id", app.GetTextPlainPage)
//
//			w := httptest.NewRecorder()
//			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.id), nil)
//			require.NoError(t, err)
//
//			router.ServeHTTP(w, req)
//
//			resp := w.Result()
//			defer resp.Body.Close()
//
//			require.Equal(t, tt.statusCode, resp.StatusCode)
//			if resp.StatusCode == http.StatusTemporaryRedirect {
//				assert.Equal(t, tt.resp.OriginalURL, resp.Header.Get("Location"))
//			}
//
//		})
//	}
//
//}
