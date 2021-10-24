package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type blockchain struct {
	blocks []*block
}

type block struct {
	Data     string
	Hash     string
	PrevHash string
}

var b *blockchain
var once sync.Once

func (b *blockchain) getLastHash() string {
	totalLength := len(b.blocks)
	if totalLength == 0 {
		return ""
	}
	return GetBlockChain().blocks[totalLength - 1].Hash
}

func (b *block) calculateHash() {
	b.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(b.Hash+b.PrevHash)))
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func createBlock(data string) *block {
	 newBlock := block{data, "",b.getLastHash()}
	 newBlock.calculateHash()
	 return &newBlock
}

func GetBlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("genesis block")
		})
	}
	return b
}

func (b *blockchain) AllBlocks() []*block {
	return b.blocks
}