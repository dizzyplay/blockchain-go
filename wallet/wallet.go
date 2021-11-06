package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"github.com/dizzyplay/blockchain-go/utils"
	"os"
)

const (
	walletFileName string = "me.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(walletFileName)
	return os.IsExist(err)
}

func createPrivateKey() *ecdsa.PrivateKey{
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	return key
}
func persistWallet(key *ecdsa.PrivateKey){
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleError(err)
	os.WriteFile(walletFileName, bytes, 0644)
}
func Wallet(){
	if w == nil {
		w := &wallet{}
		if hasWalletFile() {
			//read wallet
		}else {
			key := createPrivateKey()
			persistWallet(key)
			w.privateKey = key
		}
	}
}