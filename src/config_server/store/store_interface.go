package store

import "time"

type Store interface {
	Put(key string, typ string, expiresAt *time.Time, value string) (string, error)
	GetByName(name string) (Configurations, error)
	GetByID(id string) (Configuration, error)
	Delete(key string) (int, error)
	GetExpired() (Configurations, error)
}
