package storage

type DataStorager interface {
	Set(key, value string)
	Get(key string) (string, bool)
}

type InMemoryStorage struct {
	data map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

func (ds *InMemoryStorage) Set(key, value string) {
	ds.data[key] = value
}

func (ds *InMemoryStorage) Get(key string) (string, bool) {
	value, exists := ds.data[key]
	return value, exists
}
