package store

import (
	"encoding/json"
	"time"
)

type Configuration struct {
	ID        string
	Name      string
	Type      string
	ExpiresAt *time.Time
	Value     string
}

func (rv Configuration) StringifiedJSON() (string, error) {
	var val map[string]interface{}

	err := json.Unmarshal([]byte(rv.Value), &val)

	val["id"] = rv.ID
	val["name"] = rv.Name
	val["type"] = rv.Type
	val["expires_at"] = rv.ExpiresAt
	bytes, err := json.Marshal(&val)

	return string(bytes), err
}
