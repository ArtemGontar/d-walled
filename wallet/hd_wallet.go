package wallet

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type HDWallet struct {
	Address    string
	PrivateKey ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
}

// NewHDWallet creates a wallet with auto-generated recovery phrase. This is
// useful to create a brand-new wallet, without having to take care of the
// recovery phrase generation.
// The generated recovery phrase is returned alongside the created wallet.
func NewHDWallet(passphrase string) (*HDWallet, error) {
	ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(passphrase)
	if err != nil {
		return nil, err
	}
	return &HDWallet{
		Address: string(account.Address.Hex()),
	}, nil
}

// ImportHDWallet creates a wallet based on the recovery phrase in input. This
// is useful import or retrieve a wallet.
func ImportHDWallet(name, recoveryPhrase string) (*HDWallet, error) {
	return nil, nil
}

func (w *HDWallet) publicKey() ecdsa.PublicKey {
	return w.PublicKey
}

func (w *HDWallet) address() string {
	return w.Address
}
