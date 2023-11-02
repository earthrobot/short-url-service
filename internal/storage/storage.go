package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/earthrobot/short-url-service/internal/models"
)

type InMemoryStorage struct {
	data       map[string]string
	saveToFile bool
	file       *os.File
}

func NewInMemoryStorage(filePath string) (*InMemoryStorage, error) {
	data := make(map[string]string)
	saveToFile := filePath != ""

	storage := &InMemoryStorage{
		data:       data,
		saveToFile: saveToFile,
	}

	storage.LoadFromFile(filePath)

	if saveToFile {
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		storage.file = file
	}

	return storage, nil
}

func (ds *InMemoryStorage) Set(key, value string) error {
	ds.data[key] = value
	if ds.saveToFile {
		ds.saveToFileStorage(key, value)
	}
	return nil
}

func (ds *InMemoryStorage) Get(key string) (string, bool) {
	value, exists := ds.data[key]
	return value, exists
}

func (ds *InMemoryStorage) LoadFromFile(filePath string) error {

	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)

	for {
		data, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		shURL := models.ShortenedURL{}
		err = json.Unmarshal(data, &shURL)
		if err != nil {
			return err
		}

		ds.Set(shURL.ShortenedURL, shURL.URL)
	}

	file.Close()

	return nil
}

func (ds *InMemoryStorage) saveToFileStorage(key, value string) error {

	shURL := models.ShortenedURL{URL: value, ShortenedURL: key}
	data, err := json.Marshal(&shURL)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	_, err = ds.file.Write(data)
	return err

}
