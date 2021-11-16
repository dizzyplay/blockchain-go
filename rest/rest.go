package rest

import (
	"encoding/json"
	"fmt"
	"github.com/dizzyplay/blockchain-go/blockchain"
	"github.com/dizzyplay/blockchain-go/p2p"
	"github.com/dizzyplay/blockchain-go/utils"
	"github.com/dizzyplay/blockchain-go/wallet"
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

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayload struct {
	To string
	Amount int
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
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
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "GET txouts for an address",
		},
		{
			URL:         url("/mempool"),
			Method:      "GET",
			Description: "mempool",
		},
		{
			URL:         url("/transactions"),
			Method:      "POST",
			Description: "make tx",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to Web socket",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain())))
		return
	case "POST":
		blockchain.BlockChain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain())))
		return
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

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req.URL)
		next.ServeHTTP(w, req)
	})
}
func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	isTotal := r.URL.Query().Get("total")
	if isTotal == "true" {
		total := blockchain.TotalBalanceByAddress(address, blockchain.BlockChain())
		utils.HandleError(json.NewEncoder(rw).Encode(balanceResponse{address, total}))
	} else {
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.BlockChain())))
	}

}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

func transaction(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleError(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

//type myWalletResponse struct {
//	Address string `json:"address"`
//}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	//json.NewEncoder(rw).Encode(myWalletResponse{
	//	Address: address,
	//})
	json.NewEncoder(rw).Encode(struct {Address string `json:"address"`}{
		Address: address,
	})
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool)
	router.HandleFunc("/wallet", myWallet)
	router.HandleFunc("/transactions", transaction).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade)

	fmt.Printf("http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
