package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
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
