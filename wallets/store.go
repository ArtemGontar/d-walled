package wallets

import (
	"fmt"

	"code.vegaprotocol.io/shared/paths"
	walletstore "github.com/ArtemGontar/d-wallet/wallet/store"
)

// InitialiseStore builds a wallet Store specifically for users wallets.
func InitialiseStore(vegaHome string) (*walletstore.Store, error) {
	p := paths.New(vegaHome)
	walletsHome, err := p.CreateDataPathFor(paths.WalletsDataHome)
	if err != nil {
		return nil, fmt.Errorf("couldn't get wallets data home path: %w", err)
	}
	return walletstore.InitialiseStore(walletsHome)
}
