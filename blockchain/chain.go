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

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
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
				b.AddBlock("Genesis")
			} else {
				b.restore(checkpoint)
			}
		})
	}
	return b
}

func (b *blockchain) calculateDifficulty() int {
	block, _ := FindBlock(b.NewestHash)
	targetHash := block.PrevHash
	lastTime := block.Timestamp
	var firstTime int
	for i := 0; i < difficultyInterval-1; i++ {
		p, _ := FindBlock(targetHash)
		targetHash = p.PrevHash
		if i == difficultyInterval-2 {
			firstTime = p.Timestamp
		}
	}
	t := lastTime - firstTime
	if t < blockInterval * difficultyInterval * 60  {
		return b.CurrentDifficulty + 1
	} else {
		return b.CurrentDifficulty - 1
	}
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
