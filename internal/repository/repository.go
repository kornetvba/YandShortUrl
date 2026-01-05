package repository

import (
	"github.com/google/uuid"
)

// // ЗДЕСЬ ВМЕСТО УУИД НУЖЕН СТРИНГ, ТАК КАК ХРАНИЛИЩЕ ПОСТРЕСА ТЕПЕРЬ НА СТРИНГЕ
type URLRecord struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
}

func NewURLRecord(shortURL string, originalURL string) *URLRecord {

	return &URLRecord{CorrelationID: uuid.New().String(), ShortURL: shortURL, OriginalURL: originalURL}
}

//var URLStorages URLRecords

type URLStorage interface {
	AppendRecord(string, string) error
	AppendRecords(record *[]URLRecord) ([]URLRecord, error)
	GetRecord(string) (*URLRecord, error)
	//DownloadRecords(*config.FilePathType) error
	//SaveFile(*config.FilePathType) (err error)
	//OpenFile(*config.FilePathType) (*os.File, error)
}
