package wallet

import (
	"crypto/ecdsa"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
)

type HDWallet struct {
	Address    common.Address
	PrivateKey ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
}

// NewHDWallet creates a wallet with auto-generated recovery phrase. This is
// useful to create a brand-new wallet, without having to take care of the
// recovery phrase generation.
// The generated recovery phrase is returned alongside the created wallet.
func NewHDWallet(store Store, passphrase string) (*HDWallet, error) {
	account, err := store.SaveWallet(passphrase)
	if err != nil {
		return nil, err
	}
	return &HDWallet{
		Address: account.Address,
	}, nil
}

// ImportHDWallet creates a wallet based on the recovery phrase in input. This
// is useful import or retrieve a wallet.
func ImportHDWallet(store Store, fileName string, passphrase string, newPassphrase string) (*HDWallet, error) {
	jsonBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	account, err := store.ImportWallet(jsonBytes, passphrase, newPassphrase)
	if err != nil {
		return nil, err
	}

	return &HDWallet{
		Address: account.Address,
	}, nil
}

func (w *HDWallet) address() common.Address {
	return w.Address
}

func (w *HDWallet) publicKey() ecdsa.PublicKey {
	return w.PublicKey
}

func (w *HDWallet) privateKey() ecdsa.PrivateKey {
	return w.PrivateKey
}
