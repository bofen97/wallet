package wallet

import "github.com/ethereum/go-ethereum/accounts/keystore"

func NewHDKeyStoreNoKey(path string) *HDKeyStore {
	return &HDKeyStore{
		keyDirPath: path,
		scryptN:    keystore.LightScryptN,
		scryptP:    keystore.LightScryptP,
		Key:        keystore.Key{},
	}

}
