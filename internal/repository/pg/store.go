package pg

import (
	"database/sql"
	"github.com/kornetvba/YandShortUrl/internal/repository"
	"github.com/kornetvba/YandShortUrl/internal/service"
)

type Store struct {
	Conn *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{Conn: db}
}

func (s *Store) BootStrap() error {
	_, err := s.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS short_url_records (
			id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::VARCHAR,
			shortened_url VARCHAR(150) NOT NULL UNIQUE,
			original_url TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	_, err = s.Conn.Exec(`
		CREATE INDEX IF NOT EXISTS short_url_index ON short_url_records (shortened_url)
		`)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) AppendRecord(shortURL string, originalURL string) error {

	_, err := s.Conn.Exec(`
				INSERT INTO short_url_records ( shortened_url, original_url) 
				VALUES ($1, $2)
				`, shortURL, originalURL)
	if err != nil {
		return err
	}

	return nil
}
func (s *Store) GetRecord(shortURL string) (*repository.URLRecord, error) {
	row := s.Conn.QueryRow(`
	SELECT id, shortened_url, original_url FROM short_url_records
	WHERE shortened_url = $1
	`, shortURL)

	var url repository.URLRecord

	err := row.Scan(&url.CorrelationID, &url.ShortURL, &url.OriginalURL)
	if err != nil {

		return nil, err
	}

	return &url, nil
}

func (s *Store) AppendRecords(records *[]repository.URLRecord) ([]repository.URLRecord, error) {
	tx, err := s.Conn.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`
		INSERT INTO short_url_records
		(id, shortened_url, original_url)
		VALUES ($1, $2, $3)
`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	responseBody := []repository.URLRecord{}
	for _, record := range *records {
		shortest, err := service.HashPlainText([]byte(record.OriginalURL))
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		record.ShortURL = shortest

		_, err = stmt.Exec(record.CorrelationID, record.ShortURL, record.OriginalURL)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		record.OriginalURL = ""
		responseBody = append(responseBody, record)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return responseBody, nil

}

//func (s *Store) DownloadRecords(filePath *config.FilePathType) error {
//	if !filePath.IsEnabled() {
//		return errors.New("writing/reading to a file is disabled")
//	}
//	file, err := s.OpenFile(filePath)
//	if err != nil {
//		return err
//	}
//
//	tx, err := s.Conn.Begin()
//	defer tx.Commit()
//	if err != nil {
//		return err
//	}
//
//	scanner := bufio.NewScanner(file)
//	for scanner.Scan() {
//		var record repository.URLRecord
//		err = json.Unmarshal(scanner.Bytes(), &record)
//		if err != nil {
//			return err
//		}
//		_, err := tx.Exec(`
//			INSERT INTO short_url_records (id, shortened_url, original_url)
//			VALUES ($1, $2, $3)
//			ON CONFLICT DO NOTHING
//			`, record.ID, record.ShortURL, record.OriginalURL,
//		)
//		if err != nil {
//			trErr := tx.Rollback()
//			if trErr != nil {
//				logger.Log.Info("transaction error", zap.Error(trErr))
//			}
//			return err
//		}
//
//	}
//	tx.Commit()
//	return nil
//}

//func (s *Store) SaveFile(filePath *config.FilePathType) (err error) {
//	if !filePath.IsEnabled() {
//		return errors.New("writing/reading to a file is disabled")
//	}
//	var existingData = make(map[string]bool)
//	filepath := filePath.Dir()
//	if len(filepath) > 2 {
//		err := os.MkdirAll(filePath.Dir(), 0777)
//		if err != nil {
//			return err
//		}
//	}
//
//	if file, err := os.Open(filePath.String()); err == nil {
//
//		scanner := bufio.NewScanner(file)
//		for scanner.Scan() {
//			var record repository.URLRecord
//			err = json.Unmarshal(scanner.Bytes(), &record)
//			if err != nil {
//				return err
//			}
//			existingData[record.ShortURL] = true
//		}
//
//		file.Close()
//
//	}
//
//	file, err := os.OpenFile(filePath.String(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
//	defer file.Close()
//	if err != nil {
//		return err
//	}
//
//	rows, err := s.Conn.Query(`
//			SELECT * FROM short_url_records
//`)
//	if err != nil {
//		return err
//	}
//	if err := rows.Scan(nil); err == nil {
//		return err
//	}
//	urls := make([]repository.URLRecord, 0, 25)
//	for rows.Next() {
//		var url repository.URLRecord
//
//		err = rows.Scan(&url.ID, &url.ShortURL, &url.OriginalURL)
//		if err != nil {
//			return err
//		}
//		urls = append(urls, url)
//	}
//	for _, v := range urls {
//		_, ok := existingData[v.ShortURL]
//		if ok {
//			continue
//		}
//		if err = json.NewEncoder(file).Encode(&v); err != nil {
//			return err
//		}
//	}
//	defer file.Close()
//	return nil
//
//}
//
//func (s *Store) OpenFile(filePath *config.FilePathType) (*os.File, error) {
//	if !filePath.IsEnabled() {
//		return nil, errors.New("writing/reading to a file is disabled")
//	}
//	log.Print("Read json file...")
//	//if filename == "" {
//	//	return nil, errors.New("state saver off")
//	//}
//	file, err := os.Open(filePath.String())
//	if err != nil {
//		return nil, err
//	}
//	return file, err
//}
