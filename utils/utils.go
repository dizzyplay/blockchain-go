package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func GetHash(data interface{}) string {
	hash := sha256.Sum256([]byte(fmt.Sprint(data)))
	return fmt.Sprintf("%x",hash)
}

func ToBytes(i interface{}) []byte{
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	HandleError(enc.Encode(i))
	return buff.Bytes()
}

func FromBytes(i interface{}, d []byte) {
	dec := gob.NewDecoder(bytes.NewReader(d))
	err := dec.Decode(i)
	HandleError(err)
}
