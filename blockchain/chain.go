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

func (b *blockchain) persist() {
	db.SaveBlockChain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
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

func BlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{Height: 0}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock()
			} else {
				b.restore(checkpoint)
			}
		})
	}
	return b
}

func (b *blockchain) calculateDifficulty() int {
	allBlocks := b.Blocks()
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

func (b *blockchain) Difficulty() int {
	if b.Height == 0 {
		b.CurrentDifficulty = defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		b.CurrentDifficulty = b.calculateDifficulty()
	}
	return b.CurrentDifficulty
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}


func (b *blockchain) UTxOutsByAddress(address string) []*UTxOut {
	var UTxOuts []*UTxOut
	SpentTxIns := map[string]bool{}

	for _, block := range BlockChain().Blocks(){
		for _, txs := range block.Transactions{
			for _, input := range txs.TxIns{
				if address == input.Owner {
					SpentTxIns[input.TxId] = true
				}
			}
			for idx, txout := range txs.TxOuts {
				if txout.Owner == address {
					if _, exist := SpentTxIns[txs.Id]; !exist{
						UTxOuts = append(UTxOuts, &UTxOut{
							txs.Id,
							idx,
							txout.Amount,
						})
					}
				}
			}
		}
	}
	return UTxOuts
}

func (b *blockchain) TotalBalanceByAddress(address string) int {
	txs := b.UTxOutsByAddress(address)
	total := 0
	for _, t := range txs {
		total += t.Amount
	}
	return total
}