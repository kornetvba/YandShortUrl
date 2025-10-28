package service

import (
	"fmt"
	"github.com/kornetvba/YandShortUrl/internal/config/db"
	"hash/fnv"
)

func HashPlainText(text []byte) (string, error) {
	hashed := fnv.New32a()
	_, err := hashed.Write(text)
	if err != nil {
		return "", err
	}
	hashedText := hashed.Sum32()
	_ = db.URLStorages.AppendRecord(fmt.Sprintf("%08x", hashedText), string(text))
	return fmt.Sprintf("%08x", hashedText), nil
}

func DeHashText(shortURL string) ([]byte, error) {
	record, err := db.URLStorages.GetRecord(shortURL)
	if err != nil {
		return nil, err
	}

	return []byte(record.OriginalURL), nil
}
