package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"github.com/kornetvba/YandShortUrl/internal/handler"
	"github.com/kornetvba/YandShortUrl/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type request struct {
	contentEnc string
	acceptEnc  string
	body       []byte
}

func (req request) checkCompress() []byte {
	if req.contentEnc == "gzip" {
		req.body, _ = compressBody(req.body)
		return req.body
	}
	return req.body
}

type response struct {
	contentEnc string
	acceptEnc  string
	resBody    []byte
}

func compressBody(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write(data)
	err := gz.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func deCompressBody(data []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return io.ReadAll(gz)

}

func hashMockTest(data []byte) []byte {
	body, _ := service.HashPlainText(data)

	return []byte("/" + body)
}

func TestCompressTextPlain(t *testing.T) {
	tableTests := []struct {
		name   string
		status int
		req    request
		res    response
	}{
		{
			name:   "test1",
			status: 201,
			req: request{
				acceptEnc:  "",
				contentEnc: "",
				body:       []byte("hello world!"),
			},
			res: response{
				acceptEnc:  "",
				contentEnc: "",
				resBody:    hashMockTest([]byte("hello world!")),
			},
		},
		{
			name:   "test2",
			status: 201,
			req: request{
				acceptEnc:  "",
				contentEnc: "gzip",
				body:       []byte("hello world!"),
			},
			res: response{
				acceptEnc:  "",
				contentEnc: "",
				resBody:    hashMockTest([]byte("hello world!")),
			},
		},
		{
			name:   "test3",
			status: 201,
			req: request{
				acceptEnc:  "gzip",
				contentEnc: "gzip",
				body:       []byte("hello world!"),
			},
			res: response{
				acceptEnc:  "",
				contentEnc: "gzip",
				resBody:    hashMockTest([]byte("hello world!")),
			},
		},
	}
	for _, tt := range tableTests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			router.Use(GzipCompressMiddleWare())
			router.POST("/", handler.TextPlainPage)

			w := httptest.NewRecorder()
			body := tt.req.checkCompress()
			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			req.Header.Set("Content-Encoding", tt.req.contentEnc)
			req.Header.Set("Accept-Encoding", tt.req.acceptEnc)
			req.Header.Set("Content-Type", "text/plain")
			if err != nil {
				log.Fatal(err)
			}
			router.ServeHTTP(w, req)
			res := w.Result()
			assert.Equal(t, tt.status, res.StatusCode)
			respBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.res.contentEnc, res.Header.Get("Content-Encoding"))
			if tt.req.acceptEnc == "gzip" {
				noCompBody, err := deCompressBody(respBody)
				require.NoError(t, err)
				assert.Equal(t, tt.res.resBody, noCompBody)
			} else {
				assert.Equal(t, tt.res.resBody, respBody)
			}

		})
	}

}
