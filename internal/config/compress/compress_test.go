package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockBody(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write(data)
	err := gz.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type reqHeader struct {
	contentEncoding string
	acceptEncoding  string
}
type respHeader struct {
	contentEncoding string
	acceptEncoding  string
}

func TestGzipCompressMiddleWare(t *testing.T) {

	hand := func(c *gin.Context) {
		_, err := c.GetRawData()
		if err != nil {
			log.Fatal(err)
			return
		}
		c.Writer.WriteHeader(200)
		c.Writer.Write([]byte("mock"))
	}

	testBody, _ := mockBody([]byte("helo fjkdsa fkl sdajlkf jasldfj;la skdjf ;lkasjd fl;asdjf"))

	tableTest := []struct {
		status  int
		name    string
		reqEnc  reqHeader
		respEnc respHeader
	}{
		{
			status:  200,
			name:    "test1",
			reqEnc:  reqHeader{acceptEncoding: "", contentEncoding: ""},
			respEnc: respHeader{acceptEncoding: "", contentEncoding: ""},
		},
		{
			status:  200,
			name:    "test2",
			reqEnc:  reqHeader{acceptEncoding: "gzip", contentEncoding: ""},
			respEnc: respHeader{acceptEncoding: "", contentEncoding: "gzip"},
		},

		{
			status:  200,
			name:    "test3",
			reqEnc:  reqHeader{acceptEncoding: "gzip", contentEncoding: "gzip"},
			respEnc: respHeader{acceptEncoding: "", contentEncoding: "gzip"},
		},
	}
	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			router := gin.New()

			router.Use(gin.Recovery())
			router.Use(GzipCompressMiddleWare())
			router.POST("/", hand)

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(testBody))
			require.NoError(t, err)
			req.Header.Set("Accept-Encoding", tt.reqEnc.acceptEncoding)
			req.Header.Set("Content-Encoding", tt.reqEnc.contentEncoding)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			res := w.Result()
			assert.Equal(t, res.Header.Get("Content-Encoding"), tt.respEnc.contentEncoding)
			assert.Equal(t, res.Header.Get("Accept-Encoding"), tt.respEnc.acceptEncoding)
			assert.Equal(t, res.StatusCode, tt.status)

		})
	}

}
