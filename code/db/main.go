package db

import (
	"github.com/dgraph-io/badger/v4"
	"log"
)

var (
	DbInstance *badger.DB
)

func initdb() *badger.DB {
	opts := badger.DefaultOptions("./db/badgerdb")
	opts.Logger = nil
	DbInstance, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return DbInstance
}
func getdb() *badger.DB {
	if DbInstance == nil {
		DbInstance = initdb()
	}
	return DbInstance
}
func Start() {
	DbInstance = nil
	getdb()
}
