package wallet

import "github.com/ethereum/go-ethereum/accounts"

type Store interface {
	WalletExists(address string) bool
	ListWallets() ([]string, error)
	GetWallet(privateKey string, passphrase string) (string, error)
	SaveWallet(passphrase string) (accounts.Account, error)
	ImportWallet(jsonKey []byte, passphrase string, newPassphrase string) (accounts.Account, error)
	DeleteWallet(name string) error
	GetWalletPath(name string) string
}
