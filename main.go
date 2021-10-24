package main

import (
	"fmt"
	"github.com/dizzyplay/blockchain-go/blockchain"
)


func main() {
	chain := blockchain.GetBlockChain()
	chain.AddBlock("second")
	chain.AddBlock("third")
	for _, block := range chain.AllBlocks() {
		fmt.Printf("Data: %s \n",block.Data)
		fmt.Printf("Hash: %s \n",block.Hash)
		fmt.Printf("Prev Hash: %s \n",block.PrevHash)
	}
}
