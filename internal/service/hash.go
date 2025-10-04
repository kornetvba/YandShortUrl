package service

import (
	"fmt"
	"hash/fnv"
)

var HashMap = make(map[string][]byte)

func HashPlainText(text []byte) (string, error) {
	hashed := fnv.New32a()
	_, err := hashed.Write(text)
	if err != nil {
		return "", err
	}
	hashedText := hashed.Sum32()
	HashMap[fmt.Sprintf("%08x", hashedText)] = text

	return fmt.Sprintf("%08x", hashedText), nil
}

func DeHashText(id string) ([]byte, error) {
	value, ok := HashMap[id]
	if !ok {
		return nil, fmt.Errorf("id not found")
	}
	return value, nil
}
