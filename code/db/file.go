package db

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
)

// 保存 File
func (f *File) FileSave() error {
	db := getdb()

	return db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(f)
		if err != nil {
			return err
		}
		return txn.Set([]byte(f.HashValue), data)
	})
}

// 删除 File
func (f *File) FileDelete() error {
	db := getdb()

	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(f.HashValue))
	})
}

// 查找有没有相同的
func FileGet(hashValue string) error {
	db := getdb()

	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(hashValue))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func FileGets(hashValue string) File {
	db := getdb()

	var f File
	db.View(func(txn *badger.Txn) error {
		e, err := txn.Get([]byte(hashValue))
		if err != nil {
			return err
		}
		return e.Value(func(v []byte) error {
			return json.Unmarshal(v, &f)
		})
	})
	return f
}
