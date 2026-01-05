package dbmanager

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/kornetvba/YandShortUrl/internal/config/config"
	"github.com/kornetvba/YandShortUrl/internal/config/logger"
	"github.com/kornetvba/YandShortUrl/internal/repository"
	"github.com/kornetvba/YandShortUrl/internal/repository/pg"
	"go.uber.org/zap"
	"os"
)

type BackupPG struct {
	Store *pg.Store
}

func NewBackupStorage(store *pg.Store) *BackupPG {
	return &BackupPG{Store: store}
}

func (bs *BackupPG) SaveFile(filePath *config.FilePathType) (err error) {
	if !filePath.IsEnabled() {
		return errors.New("writing/reading to a file is disabled")
	}
	var existingData = make(map[string]bool)

	file, err := os.OpenFile(filePath.String(), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record repository.URLRecord
		err = json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return err
		}
		existingData[record.ShortURL] = true
	}

	rows, err := bs.Store.Conn.Query(`
			SELECT * FROM short_url_records
	`)
	if err != nil {
		return err
	}
	if err := rows.Scan(nil); err == nil {
		return err
	}

	for rows.Next() {
		var url repository.URLRecord

		err = rows.Scan(&url.CorrelationID, &url.ShortURL, &url.OriginalURL)
		if err != nil {
			return err
		}
		if existingData[url.ShortURL] {
			continue
		}

		if err = json.NewEncoder(file).Encode(&url); err != nil {
			return err
		}

	}
	if rows.Err() != nil {
		return err
	}

	return nil

}

func (bs *BackupPG) DownloadRecords(filePath *config.FilePathType) error {

	if !filePath.IsEnabled() {
		return errors.New("writing/reading to a file is disabled")
	}
	file, err := os.Open(filePath.String())
	if err != nil {
		return err
	}
	defer file.Close()

	tx, err := bs.Store.Conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
			INSERT INTO short_url_records (id, shortened_url, original_url)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
			`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record repository.URLRecord
		err = json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err := stmt.Exec(record.CorrelationID, record.ShortURL, record.OriginalURL)
		if err != nil {
			trErr := tx.Rollback()
			if trErr != nil {
				logger.Log.Info("transaction error", zap.Error(trErr))
			}
			return err
		}

	}
	if err := scanner.Err(); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
