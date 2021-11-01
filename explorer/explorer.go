package explorer

import (
	"fmt"
	"github.com/dizzyplay/blockchain-go/blockchain"
	"log"
	"net/http"
	"text/template"
)

const (
	templateDir string = "explorer/templates/"
)

type HomeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

var templates *template.Template

func home(rw http.ResponseWriter, r *http.Request) {
	data := HomeData{"Home", nil}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.BlockChain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusMovedPermanently)
	}
}

func Start(aPort int) {
	handler := http.NewServeMux()
	port := fmt.Sprintf(":%d", aPort)
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
