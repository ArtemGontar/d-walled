package wallet

import (
	"context"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/store_mock.go -package mocks code.vegaprotocol.io/vegawallet/wallet Store
type Store interface {
	WalletExists(name string) bool
	SaveWallet(w Wallet, passphrase string) error
	GetWallet(name, passphrase string) (Wallet, error)
	DeleteWallet(name string) error
	GetWalletPath(name string) string
	ListWallets() ([]string, error)
}

type GetWalletInfoRequest struct {
	Address string
}

type GetWalletInfoResponse struct {
	Balance big.Float `json:"balance"`
	Address string    `json:"address"`
}

func GetWalletInfo(store Store, req *GetWalletInfoRequest) (*GetWalletInfoResponse, error) {
	client, err := ethclient.DialContext(context.Background(), "https://ropsten.infura.io/v3/da52061ed5a94d03949fb39417aa8b7e")
	if err != nil {
		return nil, err
	}
	defer client.Close()

	address := common.HexToAddress(req.Address)
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}

	fBalance := new(big.Float)
	fBalance.SetString(balance.String())
	balanceEther := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18)))
	return &GetWalletInfoResponse{
			Balance: *balanceEther,
			Address: req.Address,
		},
		nil
}

type CreateWalletRequest struct {
	Passphrase string
}

type CreateWalletResponse struct {
	Address string `json:"address"`
	// PublicKey  *ecdsa.PublicKey  `json:"publicKey"`
	// PrivateKey *ecdsa.PrivateKey `json:"privateKey"`
}

func CreateWallet(store Store, req *CreateWalletRequest) (*CreateWalletResponse, error) {
	resp := &CreateWalletResponse{}

	ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(req.Passphrase)
	if err != nil {
		log.Fatal(err)
	}

	resp.Address = account.Address.Hex()
	// resp.PublicKey = key.PublicECDSA
	// resp.PrivateKey = key.PrivateECDSA
	return resp, nil
}

type ImportWalletRequest struct {
	Address       string
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

	// TODO: find file by address
	file := "./wallets/UTC--2022-01-17T10-39-00.719119000Z--8f808684001e0e8db07e6b001b4f5ea630f64067"
	ks := keystore.NewKeyStore("./wallets/tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	account, err := ks.Import(jsonBytes, req.Passphrase, req.NewPassphrase)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.Remove(file); err != nil {
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
