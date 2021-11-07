package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/dizzyplay/blockchain-go/utils"
	"math/big"
	"os"
)

const (
	walletFileName string = "me.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(walletFileName)
	return os.IsExist(err)
}

func createPrivateKey() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	return key
}
func persistWallet(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleError(err)
	os.WriteFile(walletFileName, bytes, 0644)
}
func restoreKey(filename string) (key *ecdsa.PrivateKey) {
	bytes, err := os.ReadFile(filename)
	utils.HandleError(err)
	key, err = x509.ParseECPrivateKey(bytes)
	utils.HandleError(err)
	return
}
func addressFromPrivateKey(key *ecdsa.PrivateKey) string {
	return fmt.Sprintf("%x",append(key.X.Bytes(), key.Y.Bytes()...))
}
func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			//read wallet
			w.privateKey = restoreKey(walletFileName)
		} else {
			privateKey := createPrivateKey()
			persistWallet(privateKey)
			w.privateKey = privateKey
		}
	}
	w.Address = addressFromPrivateKey(w.privateKey)
	return w
}

func Sign(payload string, w *wallet) string {
	hashedPayload, err := hex.DecodeString(payload)
	utils.HandleError(err)
	r, s, err :=ecdsa.Sign(rand.Reader, w.privateKey, hashedPayload)
	signature := append(r.Bytes(), s.Bytes()...)
	return fmt.Sprintf("%x", signature)
}

func restoreSignature(signature string) (*big.Int,*big.Int) {
	return getEachHalfBytes(signature)
}

func restoreAddress(address string) (*big.Int, *big.Int) {
	return getEachHalfBytes(address)
}

func getEachHalfBytes(data string) (*big.Int,*big.Int) {
	bytes,err := hex.DecodeString(data)
	utils.HandleError(err)
	first := bytes[:len(bytes)/2]
	second := bytes[len(bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(first)
	bigB.SetBytes(second)
	return &bigA,&bigB
}

func Verify(signature string, payloadHash string, address string) bool{
	r,s := restoreSignature(signature)
	x, y := restoreAddress(address)
	hPayload,err := hex.DecodeString(payloadHash)
	utils.HandleError(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	return ecdsa.Verify(&publicKey,hPayload,r,s)
}