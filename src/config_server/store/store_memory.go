package store

import "strconv"

type MemoryStore struct {
	db map[string]Configuration
}

var dbCounter int

func NewMemoryStore() MemoryStore {
	dbCounter = 0
	return MemoryStore{db: make(map[string]Configuration)}
}

func (store MemoryStore) Put(key string, value string) error {
	val, ok := store.db[key]

	if ok == false {
		store.db[key] = Configuration{
			Key:   key,
			Value: value,
			Id:    strconv.Itoa(dbCounter),
		}
		dbCounter++
	} else {
		val.Value = value
		store.db[key] = val
	}

	return nil
}

func (store MemoryStore) GetByKey(key string) (Configuration, error) {
	return store.db[key], nil
}

func (store MemoryStore) GetById(id string) (Configuration, error) {
	result := Configuration{}

	for _, config := range store.db {
		if config.Id == id {
			result = config
			break
		}
	}

	return result, nil
}

func (store MemoryStore) Delete(key string) (bool, error) {
	deleted := false
	result, _ := store.GetByKey(key)

	// map contains key, delete
	if len(result.Value) > 0 {
		delete(store.db, key)
		deleted = true
	}

	return deleted, nil
}
