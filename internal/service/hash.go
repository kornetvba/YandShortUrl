package service

import (
	"fmt"
	"hash/fnv"
)

func HashPlainText(text []byte) (string, error) {
	hashed := fnv.New32a()
	_, err := hashed.Write(text)
	if err != nil {
		return "", err
	}
	hashedText := hashed.Sum32()
	return fmt.Sprintf("%08x", hashedText), nil
}

//func DeHashText(shortURL string) ([]byte, error) {
//	record, err := handler.Storage.GetRecord(shortURL)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return []byte(record.OriginalURL), nil
//}
