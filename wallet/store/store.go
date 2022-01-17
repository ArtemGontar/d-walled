package walletstore

import (
	"fmt"

	dfs "github.com/ArtemGontar/d-wallet/fs"
	"github.com/ArtemGontar/d-wallet/wallet"
)

type Store struct {
	walletsHomePath string
}

func InitialiseStore(walletsHomePath string) (*Store, error) {
	if err := dfs.EnsureDir(walletsHomePath); err != nil {
		return nil, fmt.Errorf("couldn't ensure directories at %s: %w", walletsHomePath, err)
	}

	return &Store{
		walletsHomePath: walletsHomePath,
	}, nil
}

func (s *Store) WalletExists(name string) bool {
	return false
}

func (s *Store) ListWallets() ([]string, error) {
	return nil, nil
}

func (s *Store) GetWallet(name, passphrase string) (wallet.Wallet, error) {
	return nil, nil
}

func (s *Store) SaveWallet(w wallet.Wallet, passphrase string) error {
	return nil
}

func (s *Store) DeleteWallet(name string) error {
	return nil
}

func (s *Store) GetWalletPath(name string) string {
	return string("")
}
