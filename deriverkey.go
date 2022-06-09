package wallet

import (
	"crypto/ecdsa"
	"errors"
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/tyler-smith/go-bip39"
)

func CreateMnemonic() (string, error) {
	b, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatalf("Failed to NewEntropy  %s \n ", err)
		return "", err
	}

	nm, err := bip39.NewMnemonic(b)
	if err != nil {
		log.Fatalf("Failed to NewMnemonic  %s \n ", err)
		return "", err
	}
	log.Printf("CreateMnemonic :  %s \n", nm)
	return nm, nil

}

func DerivePrivateKeyFromMnemonic(nm string) (*ecdsa.PrivateKey, error) {
	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/1 ")
	if err != nil {
		log.Fatalf("Failed to ParseDerivationPath %s \n", err)
		return nil, err
	}
	seed, err := bip39.NewSeedWithErrorChecking(nm, "")
	if err != nil {
		log.Fatalf("Failed to NewSeedWithErrorChecking %s \n", err)
		return nil, err
	}

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("Failed to NewMaster %s \n", err)
		return nil, err
	}

	priv, err := DerivePrivateKey(path, masterKey)
	if err != nil {
		log.Fatalf("Failed to DerivePrivateKey %s \n", err)
		return nil, err
	}
	return priv, nil

}

func DerivePrivateKey(path accounts.DerivationPath, masterkey *hdkeychain.ExtendedKey) (*ecdsa.PrivateKey, error) {
	key := masterkey
	var err error
	for _, n := range path {
		key, err = key.Child(n)
		if err != nil {
			return nil, err
		}

	}

	privateKey, err := key.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	if err != nil {
		return nil, err
	}
	return privateKeyECDSA, nil
}

func DerivePublicKey(privatekey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	publicKey := privatekey.Public()
	publicKeyEDCSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("FAILED")
	}
	return publicKeyEDCSA, nil
}
