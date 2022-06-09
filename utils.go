package wallet

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type (
	UUID []byte
)

func NewRandom() UUID {
	uuid := make([]byte, 16)
	io.ReadFull(rand.Reader, uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return uuid

}

func WriteKeyFile(file string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(file), 0700); err != nil {
		return err
	}
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return err
	}

	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}

	f.Close()
	return os.Rename(f.Name(), file)
}
