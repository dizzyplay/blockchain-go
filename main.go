package main

import (
	"fmt"
	"github.com/dizzyplay/blockchain-go/blockchain"
	"html/template"
	"log"
	"net/http"
)

const (
	port        string = ":4000"
	templateDir string = "templates/"
)

type HomeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

var templates *template.Template

func home(rw http.ResponseWriter, r *http.Request) {
	data := HomeData{"Home", blockchain.GetBlockChain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r * http.Request) {
	fmt.Println(r)
	templates.ExecuteTemplate(rw, "add", nil)
}

func main() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	fmt.Printf("http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
