package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type blockchain struct {
	Blocks []*Block
}

type Block struct {
	Data     string
	Hash     string
	PrevHash string
}

var b *blockchain
var once sync.Once

func (b *blockchain) getLastHash() string {
	totalLength := len(b.Blocks)
	if totalLength == 0 {
		return ""
	}
	return GetBlockChain().Blocks[totalLength - 1].Hash
}

func (b *Block) calculateHash() {
	b.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(b.Hash+b.PrevHash)))
}

func (b *blockchain) AddBlock(data string) {
	b.Blocks = append(b.Blocks, createBlock(data))
}

func createBlock(data string) *Block {
	 newBlock := Block{data, "",b.getLastHash()}
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

func (b *blockchain) AllBlocks() []*Block {
	return b.Blocks
}