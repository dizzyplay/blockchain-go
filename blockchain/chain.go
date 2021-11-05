package blockchain

import (
	"github.com/dizzyplay/blockchain-go/db"
	"github.com/dizzyplay/blockchain-go/utils"
	"sync"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"current_difficulty"`
}

var b *blockchain
var once sync.Once

func persistBlockChain(b *blockchain) {
	db.SaveBlockChain(utils.ToBytes(b))
}

func Blocks(b *blockchain) []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func calculateDifficulty(b *blockchain) int {
	allBlocks := Blocks(b)
	newestBlock := allBlocks[0]
	lastCalculatedBlock := allBlocks[difficultyInterval-1]
	gap := (newestBlock.Timestamp / 60) - (lastCalculatedBlock.Timestamp / 60)
	expectedTime := blockInterval * difficultyInterval
	if gap <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if gap >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func BlockChain() *blockchain {
	once.Do(func() {
		b = &blockchain{Height: 0}
		checkpoint := db.Checkpoint()
		if checkpoint == nil {
			b.AddBlock()
		} else {
			b.restore(checkpoint)
		}
	})
	return b
}

func GetDifficulty(b *blockchain) int {
	if b.Height == 0 {
		b.CurrentDifficulty = defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		b.CurrentDifficulty = calculateDifficulty(b)
	}
	return b.CurrentDifficulty
}

func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var UTxOuts []*UTxOut
	SpentTxIns := map[string]bool{}

	for _, block := range Blocks(b) {
		for _, txs := range block.Transactions {
			for _, input := range txs.TxIns {
				if address == input.Owner {
					SpentTxIns[input.TxId] = true
				}
			}
			for idx, txout := range txs.TxOuts {
				if txout.Owner == address {
					if _, exist := SpentTxIns[txs.Id]; !exist {
						utxOut := &UTxOut{txs.Id, idx, txout.Amount}
						if !isOnMempool(utxOut) {
							UTxOuts = append(UTxOuts, utxOut)
						}
					}
				}
			}
		}
	}
	return UTxOuts
}

func TotalBalanceByAddress(address string, b *blockchain) int {
	txs := UTxOutsByAddress(address, b)
	total := 0
	for _, t := range txs {
		total += t.Amount
	}
	return total
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1, GetDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	persistBlockChain(b)
}
