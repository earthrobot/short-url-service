package storage

type DataStorage interface {
	Set(key, value string)
	Get(key string) (string, bool)
}

type inMemoryStorage struct {
	data map[string]string
}

func NewInMemoryStorage() DataStorage {
	return &inMemoryStorage{
		data: make(map[string]string),
	}
}

func (ds *inMemoryStorage) Set(key, value string) {
	ds.data[key] = value
}

func (ds *inMemoryStorage) Get(key string) (string, bool) {
	value, exists := ds.data[key]
	return value, exists
}
