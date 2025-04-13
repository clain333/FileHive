package db

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
)

func (c *Config) ConfigSave() error {
	db := getdb()

	return db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(c)
		if err != nil {
			return err
		}
		return txn.Set([]byte("config"), data)
	})
}

func ConfigGet() (*Config, error) {
	db := getdb()

	var c Config
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("config"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &c)
		})
	})
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func ConfigSet(c Config) error {
	db := getdb()
	cc, _ := json.Marshal(c)
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("config"), cc)
	})
}
