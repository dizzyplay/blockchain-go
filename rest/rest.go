package rest

import (
	"encoding/json"
	"fmt"
	"github.com/dizzyplay/blockchain-go/blockchain"
	"github.com/dizzyplay/blockchain-go/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("http://localhost%s%s", port, u)), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "",
		},
		{
			URL: url("/status"),
			Method: "GET",
			Description: "block chain status",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See a block",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

type addBlockBody struct {
	Message string
}

func blocks(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.BlockChain().Blocks()))
	case "POST":
		var b addBlockBody
		utils.HandleError(json.NewDecoder(req.Body).Decode(&b))
		blockchain.BlockChain().AddBlock(b.Message)
		rw.WriteHeader(http.StatusCreated)
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.BlockChain().Blocks()))
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.BlockChain())
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}


func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	fmt.Printf("http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
