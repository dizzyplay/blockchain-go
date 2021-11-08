package blockchain

import (
	"errors"
	"github.com/dizzyplay/blockchain-go/utils"
	"github.com/dizzyplay/blockchain-go/wallet"
	"time"
)

const (
	mineReward int = 50
)

var ErrorNoMoney = errors.New("not enough funds")
var ErrorInvalidTx = errors.New("invalid transaction")

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
	TxId      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address  string `json:"address"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxId   string
	Index  int
	Amount int
}

func (tx *Tx)sign(){
	for _, txIn := range tx.TxIns {
		txIn.Signature = wallet.Sign(tx.Id, wallet.Wallet())
	}
}

func validate(tx *Tx) bool {
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(BlockChain(),txIn.TxId)
		if prevTx == nil {
			return false
		}
		address := prevTx.TxOuts[txIn.Index].Address
		if !wallet.Verify(txIn.Signature, tx.Id, address) {
			return false
		}
	}
	return true
}

func (tx *Tx) getId() {
	tx.Id = utils.GetHash(tx)
}

func makeCoinBaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{wallet.Wallet().Address, mineReward},
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


func makeTx(from, to string, amount int) (*Tx, error) {
	if TotalBalanceByAddress(from, BlockChain()) < amount {
		return nil, ErrorNoMoney
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, BlockChain())
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxId, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	if rest := total - amount; rest > 0 {
		txOut := &TxOut{from, rest}
		txOuts = append(txOuts, txOut)
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
	tx.sign()
	if !validate(tx) {
		return nil, ErrorInvalidTx
	}
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinBaseTx(wallet.Wallet().Address)
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}

func isOnMempool(uTxOut *UTxOut) bool {
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			return input.TxId == uTxOut.TxId && input.Index == uTxOut.Index
		}
	}
	return false
}
