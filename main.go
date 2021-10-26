package main

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

func main() {
	mode := flag.NewFlagSet("mode", flag.ExitOnError)
	modeFlag := flag.String("mode", "rest", "Sets the mode of the server")
	portFlag := flag.Int("port", 4000, "Sets the mode of the server")

	if len(os.Args) < 2 {
		usage()
	}
	mode.Parse(os.Args[1:])
	fmt.Printf("\nStart %s server...on port %d\n\n",*modeFlag, *portFlag)
	switch *modeFlag {
	case "rest":
		rest.Start(*portFlag)
	case "explorer":
		explorer.Start(*portFlag)
	default:
		usage()
	}
}
