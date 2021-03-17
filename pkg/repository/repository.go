package repository

import (
	"errors"
	"sync"
)

var (
	errNotFoundValue = errors.New("Cant find value with such key")
	errListIsEmpty   = errors.New("Storage list is empty")
)

//структура для работы с  key-value хранилищем
type Storage struct {
	Mu *sync.RWMutex
	DB map[string]string
}

//интерфейс для работы с key-value хранилищем
type Storager interface {
	GetValue(key string) (string, error)
	List() ([]*Entry, error)
	Upsert(key, value string) error
	Delete(key string) error
}

//структура для добавления/изменения данных в хранилище
type Entry struct {
	Key   string `json:"key"`
	Value string ` json:"value"`
}

//создание нового экземпляра хранилища
func NewStorage() *Storage {
	db := make(map[string]string)
	mu := &sync.RWMutex{}
	return &Storage{
		Mu: mu,
		DB: db,
	}
}

//получение данных из хранилища
func (s *Storage) GetValue(key string) (string, error) {
	s.Mu.RLock()

	value, ok := s.DB[key]

	s.Mu.RUnlock()

	if !ok {

		return "", errNotFoundValue

	}

	return value, nil
}

//вывод списка данных из хранилища
func (s *Storage) List() ([]*Entry, error) {
	s.Mu.RLock()

	//если в хранилище нет записей, информируем об этом пользователя
	if len(s.DB) == 0 {
		s.Mu.RUnlock()
		return nil, errListIsEmpty
	}
	s.Mu.RUnlock()

	var entries []*Entry

	s.Mu.RLock()

	//проходимся по всем значениям хранилища и записываем их в слайс записей для выдачи пользователю
	for k, v := range s.DB {
		entry := &Entry{}

		entry.Key = k

		entry.Value = v

		entries = append(entries, entry)
	}
	s.Mu.RUnlock()

	return entries, nil
}

//добавление или внесение изменений в хранилище
func (s *Storage) Upsert(key, value string) error {
	s.Mu.Lock()

	s.DB[key] = value

	s.Mu.Unlock()

	return nil
}

//удаление данных из хранилища
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
