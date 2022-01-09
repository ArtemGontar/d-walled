package walletstore

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/ArtemGontar/d-wallet/dcrypto"
	dfs "github.com/ArtemGontar/d-wallet/fs"
	"github.com/ArtemGontar/d-wallet/wallet"
)

type Store struct {
	walletsHome string
}

func InitialiseStore(walletsHome string) (*Store, error) {
	if err := dfs.EnsureDir(walletsHome); err != nil {
		return nil, fmt.Errorf("couldn't ensure directories at %s: %w", walletsHome, err)
	}

	return &Store{
		walletsHome: walletsHome,
	}, nil
}

func (s *Store) DeleteWallet(name string) error {
	walletPath := s.walletPath(name)

	if exists, _ := dfs.PathExists(walletPath); !exists {
		return wallet.ErrWalletDoesNotExists
	}

	if err := os.Remove(walletPath); err != nil {
		return fmt.Errorf("couldn't remove wallet file at %s: %w", walletPath, err)
	}
	return nil
}

func (s *Store) WalletExists(name string) bool {
	walletPath := s.walletPath(name)

	exists, _ := dfs.PathExists(walletPath)
	return exists
}

func (s *Store) ListWallets() ([]string, error) {
	walletsParentDir, walletsDir := filepath.Split(s.walletsHome)
	entries, err := fs.ReadDir(os.DirFS(walletsParentDir), walletsDir)
	if err != nil {
		return nil, fmt.Errorf("couldn't read directory at %s: %w", s.walletsHome, err)
	}
	wallets := make([]string, len(entries))
	for i, entry := range entries {
		wallets[i] = entry.Name()
	}
	sort.Strings(wallets)
	return wallets, nil
}

func (s *Store) GetWallet(name, passphrase string) (wallet.Wallet, error) {
	buf, err := fs.ReadFile(os.DirFS(s.walletsHome), name)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file at %s: %w", s.walletsHome, err)
	}

	decBuf, err := dcrypto.Decrypt(buf, passphrase)
	if err != nil {
		if err.Error() == "cipher: message authentication failed" {
			return nil, wallet.ErrWrongPassphrase
		}
		return nil, err
	}

	versionedWallet := &struct {
		Version uint32 `json:"version"`
	}{}

	err = json.Unmarshal(decBuf, versionedWallet)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal wallet verion: %w", err)
	}

	w := &wallet.HDWallet{}
	err = json.Unmarshal(decBuf, w)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal wallet: %w", err)
	}

	// The wallet name is not saved in the file to avoid de-synchronisation
	// between file name and file content
	w.SetName(name)

	return w, nil
}

func (s *Store) SaveWallet(w wallet.Wallet, passphrase string) error {
	buf, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("couldn't marshal wallet: %w", err)
	}

	encBuf, err := dcrypto.Encrypt(buf, passphrase)
	if err != nil {
		return fmt.Errorf("couldn't encrypt wallet: %w", err)
	}

	walletPath := s.walletPath(w.Name())
	err = dfs.WriteFile(walletPath, encBuf)
	if err != nil {
		return fmt.Errorf("couldn't write wallet file at %s: %w", walletPath, err)
	}
	return nil
}

func (s *Store) GetWalletPath(name string) string {
	return s.walletPath(name)
}

func (s *Store) walletPath(name string) string {
	return filepath.Join(s.walletsHome, name)
}
