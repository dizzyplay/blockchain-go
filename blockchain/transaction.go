package blockchain

import (
	"github.com/dizzyplay/blockchain-go/utils"
	"time"
)

const (
	mineReward int = 50
)

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	Owner  string
	Amount int
}

type TxOut struct {
	Owner  string
	Amount int
}

func (tx *Tx) getId() {
	tx.Id = utils.GetHash(tx)
}

func makeCoinBaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"COINBASE", mineReward},
	}
	txOuts := []*TxOut{
		{"me", mineReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}
