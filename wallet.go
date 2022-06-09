package wallet

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type (
	HDWallet struct {
		Address    common.Address
		HdKeyStore *HDKeyStore
	}
)

func NewWallet(keypath string) (*HDWallet, error) {
	mn, err := CreateMnemonic()
	if err != nil {
		fmt.Println("Failed to CreateMnemonic ", err)
		return nil, err
	}
	privateKey, err := DerivePrivateKeyFromMnemonic(mn)
	if err != nil {
		fmt.Println("Failed to DerivePrivateKeyFromMnemonic ", err)
		return nil, err
	}

	publicKey, err := DerivePublicKey(privateKey)
	if err != nil {
		fmt.Println("Failed to DerivePublicKey ", err)
		return nil, err
	}
	address := crypto.PubkeyToAddress(*publicKey)
	hdks := NewHDKeyStore(keypath, privateKey)
	return &HDWallet{Address: address, HdKeyStore: hdks}, nil

}

func (w HDWallet) StoreKey(pass string) error {
	filename := w.HdKeyStore.JoinPath(w.Address.Hex())
	return w.HdKeyStore.StoreKey(filename, &w.HdKeyStore.Key, pass)
}
