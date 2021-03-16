package repository

import (
	"errors"
	"sync"
)

var (
	errNotFoundValue = errors.New("Cant find value with such key")
)

type Storage struct {
	Mu *sync.RWMutex
	DB map[string]string
}

type Entry struct {
	Key   string `json:"key,omitempty"`
	Value string ` json:"value,omitempty"`
}

func NewStorage() *Storage {
	db := make(map[string]string)
	mu := &sync.RWMutex{}
	return &Storage{
		Mu: mu,
		DB: db,
	}
}

func (s *Storage) GetValue(key string) (string, error) {
	s.Mu.RLock()

	value, ok := s.DB[key]

	s.Mu.RUnlock()

	if !ok {

		return "", errNotFoundValue

	}

	return value, nil
}

func (s *Storage) List() ([]*Entry, error) {
	var entries []*Entry

	s.Mu.RLock()

	for k, v := range s.DB {
		entry := &Entry{}

		entry.Key = k

		entry.Value = v

		entries = append(entries, entry)
	}
	s.Mu.RUnlock()

	return entries, nil
}

func (s *Storage) Upsert(key, value string) error {
	s.Mu.Lock()

	s.DB[key] = value

	s.Mu.Unlock()

	return nil
}

func (s *Storage) Delete(key string) error {
	s.Mu.Lock()

	_, ok := s.DB[key]

	if !ok {
		s.Mu.Unlock()
		return errNotFoundValue

	}

	delete(s.DB, key)

	s.Mu.Unlock()
	return nil
}

//curl -H 'Content-Type:application/json' -d '{"key":"122223","value":"vasya"}' -X POST http://127.0.0.1:8086/upsert
