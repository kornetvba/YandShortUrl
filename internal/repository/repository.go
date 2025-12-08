package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/kornetvba/YandShortUrl/internal/config/config"
	"log"
	"os"
)

type URLRecord struct {
	ID          uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

func NewURLRecord(shortURL string, originalURL string) *URLRecord {
	return &URLRecord{ID: uuid.New(), ShortURL: shortURL, OriginalURL: originalURL}
}

type URLRecords []URLRecord

func NewURLRecords() *URLRecords {
	return &URLRecords{}
}

//var URLStorages URLRecords

type URLStorage interface {
	AppendRecord(string, string) error
	GetRecord(string) (*URLRecord, error)
	LoadRecords(*config.FilePathType) error
	SaveFile(*config.FilePathType) (err error)
	ReadFile(*config.FilePathType) (*os.File, error)
}

func (ur *URLRecords) AppendRecord(shortURL string, originalURL string) error {
	double, _ := ur.GetRecord(shortURL)
	if double != nil {
		return errors.New("double url")
	}

	record := NewURLRecord(shortURL, originalURL)
	*ur = append(*ur, *record)
	return nil
}

func (ur *URLRecords) GetRecord(shortURL string) (*URLRecord, error) {
	for _, record := range *ur {
		if record.ShortURL == shortURL {
			return &record, nil
		}
	}
	return nil, errors.New("id not found")
}

func (ur *URLRecords) LoadRecords(filePath *config.FilePathType) error {
	if !filePath.IsEnabled() {
		return errors.New("writing/reading to a file is disabled")
	}
	file, err := ur.ReadFile(filePath)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record URLRecord
		err = json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return err
		}
		*ur = append(*ur, record)
	}

	defer file.Close()
	log.Print("File read success")
	return scanner.Err()
}

func (ur *URLRecords) SaveFile(filePath *config.FilePathType) (err error) {
	if !filePath.IsEnabled() {
		return errors.New("writing/reading to a file is disabled")
	}
	var existingData = make(map[string]bool)
	filepath := filePath.Dir()
	if len(filepath) > 2 {
		err := os.MkdirAll(filePath.Dir(), 0777)
		if err != nil {
			return err
		}
	}

	if file, err := os.Open(filePath.String()); err == nil {

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var record URLRecord
			err = json.Unmarshal(scanner.Bytes(), &record)
			if err != nil {
				return err
			}
			existingData[record.ShortURL] = true
		}

		file.Close()

	}

	file, err := os.OpenFile(filePath.String(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	for _, v := range *ur {
		_, ok := existingData[v.ShortURL]
		if ok {
			continue
		}
		if err = json.NewEncoder(file).Encode(&v); err != nil {
			return err
		}
	}

	defer file.Close()
	return nil
}

func (ur *URLRecords) ReadFile(filePath *config.FilePathType) (*os.File, error) {
	if !filePath.IsEnabled() {
		return nil, errors.New("writing/reading to a file is disabled")
	}
	log.Print("Read json file...")
	//if filename == "" {
	//	return nil, errors.New("state saver off")
	//}
	file, err := os.Open(filePath.String())
	if err != nil {
		return nil, err
	}
	return file, err
}
