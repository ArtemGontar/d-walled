package main

import (
	"fmt"

	"github.com/ArtemGontar/d-wallet/wallet"
	"github.com/ArtemGontar/d-wallet/wallets"
)

func main() {
	req := &wallet.GetWalletInfoRequest{
		Wallet:     "walletName",
		Passphrase: "ssss",
	}
	path := "Users/anduser"
	s, err := wallets.InitialiseStore(path)
	if err != nil {
		fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	wallet.GetWalletInfo(s, req)

}
