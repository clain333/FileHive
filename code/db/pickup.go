package db

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"log"
)

func (p *PickupInfo) PickupSave() error {
	db := getdb()

	return db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return txn.Set([]byte("id_"+p.PickupCode), data)
	})
}

func PickupDelete(id string) error {
	db := getdb()

	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte("id_" + id))
	})
}

func PickupSet(id string) error {
	db := getdb()

	pp, err := PickupGet(id)
	if err != nil {
		return fmt.Errorf("PickupGet 失败: %w", err)
	}

	pp.IsDownloaded = !pp.IsDownloaded

	data, err := json.Marshal(pp)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}

	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("id_"+pp.PickupCode), data)
	})
}
func PickupSets(id string, p *PickupInfo) error {
	db := getdb()

	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}

	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("id_"+id), data)
	})
}

func PickupGet(id string) (*PickupInfo, error) {
	db := getdb()

	var p PickupInfo
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("id_" + id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &p)
		})
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}
func PickupGetAll() ([]*PickupInfo, error) {
	db := getdb()

	var result []*PickupInfo
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("id_")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {

				var p PickupInfo
				if err := json.Unmarshal(v, &p); err != nil {
					return err
				}
				result = append(result, &p)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("遍历出错:", err)
	}
	return result, nil
}

func PickupDelteAll() {
	db := getdb()
	err := db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // 不需要获取值，提高效率
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("id_")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.KeyCopy(nil)
			err := txn.Delete(key)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func PickSetAll(bool2 bool) {
	db := getdb()
	db.Update(func(txn *badger.Txn) error {
		// 使用 Iterator 遍历数据库中的所有键值对
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		// 遍历数据库中的所有键值对
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()

			// 如果键以 "id_" 开头
			if string(key[:3]) == "id_" {
				// 获取键对应的值
				err := item.Value(func(val []byte) error {
					var data PickupInfo
					// 将值反序列化成 PickupInfo 结构体
					err := json.Unmarshal(val, &data)
					if err != nil {
						return err
					}

					data.IsDownloaded = bool2

					// 将修改后的数据重新编码成 JSON 格式
					newValue, err := json.Marshal(data)
					if err != nil {
						return err
					}

					// 更新数据库中的值
					err = txn.Set(key, newValue)
					return err
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
