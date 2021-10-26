package cli

import (
	"flag"
	"fmt"
	"github.com/dizzyplay/blockchain-go/explorer"
	"github.com/dizzyplay/blockchain-go/rest"
	"os"
)

func usage(){
	fmt.Println("usage:")
	fmt.Println("	-mode=rest")
	fmt.Println("	-mode=explorer")
	os.Exit(0)
}

func Start() {
	mode := flag.String("mode", "rest", "Sets the mode of the server")
	port := flag.Int("port", 4000, "Sets the mode of the server")
	flag.Parse()

	switch *mode {
	case "rest":
		fmt.Printf("\nStart %s server...on port %d\n\n",*mode, *port)
		rest.Start(*port)
	case "explorer":
		fmt.Printf("\nStart %s server...on port %d\n\n",*mode, *port)
		explorer.Start(*port)
	default:
		usage()
	}
}

