package main

import (
	"github.com/dizzyplay/blockchain-go/cli"
	"github.com/dizzyplay/blockchain-go/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
