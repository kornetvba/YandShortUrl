package repository

import (
	"github.com/google/uuid"
)

type URLRecord struct {
	ID          uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

func NewURLRecord(shortURL string, originalURL string) *URLRecord {
	return &URLRecord{ID: uuid.New(), ShortURL: shortURL, OriginalURL: originalURL}
}

//var URLStorages URLRecords

type URLStorage interface {
	AppendRecord(string, string) error
	GetRecord(string) (*URLRecord, error)
	//DownloadRecords(*config.FilePathType) error
	//SaveFile(*config.FilePathType) (err error)
	//OpenFile(*config.FilePathType) (*os.File, error)
}
