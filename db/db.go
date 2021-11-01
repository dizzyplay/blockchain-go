package db

import (
	"github.com/boltdb/bolt"
	"github.com/dizzyplay/blockchain-go/utils"
)

const (
	dbName = "blockchain.db"
	dataBucket = "data"
	blockBucket = "blocks"
	checkpoint = "checkpoint"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		db = dbPointer
		utils.HandleError(err)
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleError(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blockBucket))
			return err
		})
		utils.HandleError(err)
		
	}
	return db
}

func SaveBlock(hash string, data []byte) {
	save([]byte(hash), data, blockBucket)
}

func SaveBlockChain(data []byte) {
	save([]byte(checkpoint), data, dataBucket)
}

func save(key []byte, data []byte, bucketName string) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		err := bucket.Put(key, data)
		return err
	})
	utils.HandleError(err)
}

func Checkpoint() []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

func Close(){
	DB().Close()
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}