package blockchain

import (
	"fmt"
	"github.com/dizzyplay/blockchain-go/db"
	"github.com/dizzyplay/blockchain-go/utils"
	"sync"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}


var b *blockchain
var once sync.Once

func (b *blockchain) persist() {
	db.SaveBlockChain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height + 1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func BlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("genesis block")
			}else {
				b.restore(checkpoint)
				fmt.Printf("n hash %s height %d", b.NewestHash, b.Height)
			}
		})
	}
	return b
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}