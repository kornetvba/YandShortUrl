package handler

import (
	"fmt"
	"github.com/kornetvba/YandShortUrl/internal/service"
	"io"
	"log"
	"net/http"
)

func TextPlainPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Header.Get("Content-Type") != "text/plain" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	defer func() {
		err := req.Body.Close()
		if err != nil {
			log.Printf("Body close error: %v", err)
		}

	}()
	PostText, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	hashText, errors := service.HashPlainText(PostText)

	if errors != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	result := fmt.Sprintf("%s%s%s%s", "http://", req.Host, req.URL.String(), hashText)
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Content-Length", fmt.Sprint(len(result)))
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write([]byte(result))
	if err != nil {
		log.Printf("Body err: %v", err)
		return
	}

}

func GetTextPlainPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	//if req.Header.Get("Content-Type") != "text/plain" {
	//	res.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	id := req.PathValue("id")
	resURL, err := service.DeHashText(id)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", string(resURL))
	res.WriteHeader(http.StatusTemporaryRedirect)

}
