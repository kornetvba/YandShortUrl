package memory_manager

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/kornetvba/YandShortUrl/internal/config/config"
	"github.com/kornetvba/YandShortUrl/internal/repository"
	"github.com/kornetvba/YandShortUrl/internal/repository/memory"

	"os"
)

type BackupMemory struct {
	MemoryStorage *memory.URLRecords
}

func NewBackupStorage(store *memory.URLRecords) *BackupMemory {
	return &BackupMemory{MemoryStorage: store}
}

func (bm *BackupMemory) SaveFile(filePath *config.FilePathType) (err error) {
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
	file, err := os.OpenFile(filePath.String(), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record repository.URLRecord
		err = json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return err
		}
		existingData[record.ShortURL] = true
	}
	for _, v := range *bm.MemoryStorage {
		if existingData[v.ShortURL] {
			continue
		}
		if err = json.NewEncoder(file).Encode(&v); err != nil {
			return err
		}

	}

	return nil

}

func (bm *BackupMemory) DownloadRecords(filePath *config.FilePathType) error {
	if !filePath.IsEnabled() {
		return errors.New("writing/reading to a file is disabled")
	}
	file, err := os.Open(filePath.String())
	defer file.Close()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record repository.URLRecord
		err = json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return err
		}

		*bm.MemoryStorage = append(*bm.MemoryStorage, record)

	}

	return nil
}
