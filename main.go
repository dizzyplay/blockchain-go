package main

import (
	"fmt"
	"github.com/dizzyplay/blockchain-go/blockchain"
	"html/template"
	"log"
	"net/http"
)

const port string = ":4000"

type HomeData struct {
	PageTitle string
	Blocks []*blockchain.Block
}

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/pages/home.gohtml"))
		data := HomeData{"Home", blockchain.GetBlockChain().AllBlocks() }
		tmpl.Execute(writer, data)
	})
	fmt.Printf("http://localhost%s",port)
	log.Fatal(http.ListenAndServe(port, nil))
}
