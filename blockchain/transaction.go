package blockchain

import (
	"errors"
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

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

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

func makeTx(from, to string , amount int) (*Tx, error) {
	 if BlockChain().TotalBalanceByAddress(from) < amount {
		 return nil, errors.New("not enough money")
	 }
	 var txIns []*TxIn
	 var txOuts []*TxOut
	 total := 0
	 oldTxs := BlockChain().TxOutsByAddress(from)
	 for _, otx := range oldTxs {
		 if total > amount {
			 break
		 }
		 txIn := &TxIn{otx.Owner, otx.Amount}
		 txIns = append(txIns, txIn)
		 total += otx.Amount
	 }

	 rest := total-amount

	 if rest != 0 {
		 txOuts = append(txOuts, &TxOut{
			 from,
			 rest,
		 })
	 }

	txOuts = append(txOuts, &TxOut{
		to,
		amount,
	})

	 tx := &Tx{
		 "",
		 int(time.Now().Unix()),
		txIns,
		txOuts,
	 }
	 tx.getId()
	 return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("me", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinBaseTx("me")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}