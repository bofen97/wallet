package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type (
	HDKeyStore struct {
		keyDirPath string
		scryptN    int
		scryptP    int
		Key        keystore.Key
	}
)

func NewHDKeyStore(path string, privateKey *ecdsa.PrivateKey) *HDKeyStore {
	uuid := []byte(NewRandom())
	var uuid16 [16]byte
	copy(uuid16[:], uuid)
	key := keystore.Key{
		Id:         uuid16,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}
	return &HDKeyStore{
		keyDirPath: path,
		scryptN:    keystore.LightScryptN,
		scryptP:    keystore.LightScryptP,
		Key:        key}

}

func (ks HDKeyStore) StoreKey(filename string, key *keystore.Key, auth string) error {
	//对称加密，拿到加密后的私钥
	keyjson, err := keystore.EncryptKey(key, auth, ks.scryptN, ks.scryptP)
	if err != nil {
		return err
	}
	return WriteKeyFile(filename, keyjson)
}

func (ks HDKeyStore) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}

	return filepath.Join(ks.keyDirPath, filename)
}

func (ks *HDKeyStore) GetKey(addr common.Address, filename, auth string) (*keystore.Key, error) {
	keyjson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	//通过加密后的私钥和密码 解出明文私钥
	key, err := keystore.DecryptKey(keyjson, auth)
	if err != nil {
		return nil, err
	}
	if key.Address != addr {
		return nil, fmt.Errorf("KEY CONTENT MISMATCH : GOT ACCOUNT %x , EXPECTED %x ", key.Address, addr)
	}
	ks.Key = *key
	return key, nil
}

func (ks HDKeyStore) SignTx(address common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	signer := types.NewEIP155Signer(chainID)
	//使用私钥进行签名
	signed, err := types.SignTx(tx, signer, ks.Key.PrivateKey)
	if err != nil {
		return nil, err
	}

	msg, err := signed.AsMessage(signer, nil)
	if err != nil {
		return nil, err
	}

	sender := msg.From()

	if sender != address {
		return nil, fmt.Errorf("SIGNER MISMATCH : EXPECTED %s , GOT %s ", address.Hex(), sender.Hex())
	}

	return signed, nil

}

//操作智能合约时进行签名
func (ks HDKeyStore) NewTransactionOps() (*bind.TransactOpts, error) {
	return bind.NewKeyedTransactorWithChainID(ks.Key.PrivateKey, ChainID)
}
