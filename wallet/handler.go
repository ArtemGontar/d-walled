package wallet

import (
	"context"
	"crypto/ecdsa"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type GetWalletInfoRequest struct {
	Address string
}

type GetWalletInfoResponse struct {
	Balance *big.Float `json:"balance"`
	Address string     `json:"address"`
}

func GetWalletInfo(store Store, client *ethclient.Client, req *GetWalletInfoRequest) (*GetWalletInfoResponse, error) {
	address := common.HexToAddress(req.Address)
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}

	balanceEther := convertWeiToEth(balance)

	return &GetWalletInfoResponse{
			Balance: balanceEther,
			Address: req.Address,
		},
		nil
}

func convertWeiToEth(balance *big.Int) *big.Float {
	fBalance := new(big.Float)
	fBalance.SetString(balance.String())
	balanceEther := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18)))
	return balanceEther
}

type CreateWalletRequest struct {
	Passphrase string
}

type CreateWalletResponse struct {
	Address    string            `json:"address"`
	PublicKey  *ecdsa.PublicKey  `json:"publicKey"`
	PrivateKey *ecdsa.PrivateKey `json:"privateKey"`
}

func CreateWallet(store Store, req *CreateWalletRequest) (*CreateWalletResponse, error) {
	resp := &CreateWalletResponse{}

	account, err := NewHDWallet(store, req.Passphrase)
	if err != nil {
		return nil, err
	}
	resp.Address = account.Address.Hex()
	resp.PublicKey = &account.PublicKey
	resp.PrivateKey = &account.PrivateKey
	return resp, nil
}

type ImportWalletRequest struct {
	PrivateKey    string
	File          os.File
	Passphrase    string
	NewPassphrase string
}

type ImportWalletResponse struct {
	Address string `json:"address"`
	// PublicKey  *ecdsa.PublicKey  `json:"publicKey"`
	// PrivateKey *ecdsa.PrivateKey `json:"privateKey"`
}

func ImportWallet(store Store, req *ImportWalletRequest) (*ImportWalletResponse, error) {
	resp := &ImportWalletResponse{}
	jsonBytes, err := ioutil.ReadFile(req.File.Name())
	if err != nil {
		log.Fatal(err)
	}

	// TODO: find file by address
	ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)

	account, err := ks.Import(jsonBytes, req.Passphrase, req.NewPassphrase)

	if err != nil {
		log.Fatal(err)
	}

	if err := os.Remove(req.File.Name()); err != nil {
		log.Fatal(err)
	}

	resp.Address = account.Address.Hex()
	// resp.PublicKey = key.PublicECDSA
	// resp.PrivateKey = key.PrivateECDSA
	return resp, nil
}

type DeleteWalletRequest struct {
	Wallet string
}

func DeleteWallet(store Store, req *DeleteWalletRequest) error {
	return nil
}
