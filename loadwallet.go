package wallet

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/howeyc/gopass"
)

func LoadWallet(filename, datadir string) (*HDWallet, error) {
	hdks := NewHDKeyStoreNoKey(datadir)
	fmt.Println("Please input password for : ", filename)
	pass, _ := gopass.GetPasswd()
	fromaddr := common.HexToAddress(filename)
	_, err := hdks.GetKey(fromaddr, hdks.JoinPath(filename), string(pass))
	if err != nil {
		log.Panic("Failed to GetKey ", err)
	}
	return &HDWallet{fromaddr, hdks}, nil
}
